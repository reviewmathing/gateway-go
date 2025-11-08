package logger

import (
	"log/slog"
	"os"
)

func TestSetUp() {
	appHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     nil,
	})
	initApp(appHandler)

	httpHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     nil,
	})
	initHttp(httpHandler)
}
