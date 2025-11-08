package main

import (
	"context"
	config2 "gateway-go/internal/config"
	handler "gateway-go/internal/health"
	"gateway-go/internal/logger"
	"gateway-go/internal/router"
	"gateway-go/proxy"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logConfigData, err := config2.GetData(logger.LogConfigFileName)
	if err != nil {
		log.Fatal("log config read fail", err)
	}
	config, err := logger.ReadConfig(logConfigData)
	if err != nil {
		log.Fatal("Failed to get Log Config", err)
		return
	}
	Close, err := logger.SetUp(config)
	defer Close()
	if err != nil {
		log.Fatal("Log init Fail ", err)
		return
	}

	routerConfigData, err := config2.GetData(router.RouterConfigName)
	if err != nil {
		log.Fatal("Failed to get Router Config", err)
		return
	}
	newRouter, err := router.NewRouter(routerConfigData)
	if err != nil {
		logger.App.Error("Failed to initialize router", "error", err)
		return
	}

	newProxy := proxy.NewProxy(newRouter)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthHandler)
	mux.Handle("/", &newProxy)

	// HTTP 서버 설정
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// 서버를 고루틴에서 실행
	go func() {
		logger.App.Info("Gateway server starting", "address", ":8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.App.Error("Server error", "error", err)
			os.Exit(0)
		}
	}()

	// 시그널 대기 (Ctrl+C, kill 등)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.App.Info("Shutting down server gracefully...")

	// 30초 타임아웃으로 graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.App.Error("Forced shutdown", "error", err)
	}

	logger.App.Info("Server stopped")
}
