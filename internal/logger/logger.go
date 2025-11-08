package logger

import (
	"log/slog"
	"net/http"
	"os"
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

func SetUp(root string) error {
	logFolder := root + "/log"
	file, err := os.OpenFile(logFolder+"/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil
	}
	handler := slog.NewTextHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     nil,
	})
	initApp(handler)

	file2, err := os.OpenFile(logFolder+"/middle.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	textHandler := slog.NewTextHandler(file2, &slog.HandlerOptions{
		AddSource: false,
		Level:     nil,
	})
	initHttp(textHandler)
	return nil
}

func initHttp(textHandler *slog.TextHandler) {
	HTTP = HttpLogger{
		slog.New(textHandler),
	}
}

func initApp(handler *slog.TextHandler) {
	App = AppLogger{
		slog.New(handler),
	}
}
