package logger

import (
	"log/slog"
	"net/http"
)

type HttpLogger struct {
	*slog.Logger
}

func (hl *HttpLogger) LogTransaction(r http.Request, status int) {
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

// SetUp initializes all loggers with the given config.
// Returns a cleanup function that safely closes all open log files.
// The cleanup function is safe to call even if initialization failed.
func SetUp(config Config) (func(), error) {
	var appCloser func()
	var httpCloser func()

	handler, closer, err := config.appHandler()
	if err != nil {
		// Return safe cleanup even on error
		return func() {}, err
	}
	appCloser = closer
	initApp(handler)

	httpHandler, closer, err := config.httpHandler()
	if err != nil {
		// Return safe cleanup even on error
		return func() {
			if appCloser != nil {
				appCloser()
			}
		}, err
	}
	httpCloser = closer

	if httpHandler != nil {
		initHttp(httpHandler)
	}

	// Always return valid cleanup function
	return func() {
		if appCloser != nil {
			appCloser()
		}
		if httpCloser != nil {
			httpCloser()
		}
	}, nil
}

// initHttp initializes HTTP logger with the given handler.
// handler must not be nil (checked by caller).
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
