package main

import (
	"gateway-go/internal/logger"
	"gateway-go/internal/router"
	"gateway-go/proxy"
	"log"
	"net/http"
	"os"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("get root dir fail :", err)
		return
	}
	logger.SetUp(dir)

	file, err := os.ReadFile(dir + "/config.yml")
	if err != nil {
		log.Println("open fail config.yml : ", err)
		return
	}
	newRouter, err := router.NewRouter(file)
	if err != nil {
		logger.App.Error("init router fail : ", err)
		return
	}

	newProxy := proxy.NewProxy(newRouter)

	logger.App.Info("Gateway server starting on :8080")

	if err := http.ListenAndServe(":8080", &newProxy); err != nil {
		logger.App.Error("system 종료 : ", err)
	}
}
