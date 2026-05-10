package observability

import (
	"log/slog"
	"os"
)

var L *slog.Logger

func InitLogger() *slog.Logger {
	if L != nil {
		return L
	}
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	L = slog.New(h)
	slog.SetDefault(L)
	return L
}
