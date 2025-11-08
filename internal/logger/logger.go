package logger

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

type HttpLogger struct {
	*slog.Logger
}

func (hl *HttpLogger) LogTransection(r http.Request, status int) {
	hl.Info("HTTP Request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Int("status", status),
		slog.String("user_agent", r.UserAgent()),
	)
}

type AppLogger struct {
	*slog.Logger
}

var (
	App  AppLogger
	HTTP HttpLogger
)

func SetUp(root string) (func(), error) {
	logFolder := filepath.Join(root, "log")
	appFileLogger := &lumberjack.Logger{
		Filename:  filepath.Join(logFolder, "app.log"),
		MaxSize:   100,
		MaxAge:    1,
		LocalTime: true,
		Compress:  true,
	}
	multiWriter := io.MultiWriter(os.Stdout, appFileLogger)
	appHandler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		AddSource: true,
		Level:     nil,
	})
	initApp(appHandler)

	httpFileLogger := &lumberjack.Logger{
		Filename:  filepath.Join(logFolder, "http.log"),
		MaxSize:   100,
		MaxAge:    1,
		LocalTime: true,
		Compress:  true,
	}
	textHandler := slog.NewTextHandler(httpFileLogger, &slog.HandlerOptions{
		AddSource: false,
		Level:     nil,
	})
	initHttp(textHandler)
	return func() {
		_ = httpFileLogger.Close()
		_ = appFileLogger.Close()
	}, nil
}

func initHttp(textHandler *slog.TextHandler) {
	HTTP = HttpLogger{
		slog.New(textHandler),
	}
}

func initApp(textHandler slog.Handler) {
	App = AppLogger{
		slog.New(textHandler),
	}
}
