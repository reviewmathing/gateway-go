package logger

import (
	"log/slog"
	"net/http"
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

type Config interface {
	appHandler() (slog.Handler, func(), error)
	httpHandler() (slog.Handler, func(), error)
}

func SetUp(config Config) (func(), error) {
	handler, appCloser, err := config.appHandler()
	if err != nil {
		return nil, err
	}
	initApp(handler)

	httpHandler, httpCloser, err := config.httpHandler()
	if err != nil {
		return nil, err
	}
	if httpHandler != nil {
		initHttp(httpHandler)
	}
	return func() {
		appCloser()
		if httpCloser != nil {
			httpCloser()
		}
	}, nil
}

func initHttp(handler slog.Handler) {
	HTTP = HttpLogger{
		slog.New(handler),
	}
}

func initApp(handler slog.Handler) {
	App = AppLogger{
		slog.New(handler),
	}
}
