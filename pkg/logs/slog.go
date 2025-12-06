package logs

import (
	"log/slog"
	"os"

	"github.com/bmstu-itstech/scriptum-back/internal/config"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs/handlers/slogpretty"
)

const (
	logDebug = "debug"
	logInfo  = "info"
	logWarn  = "warn"
	logError = "error"
)

func NewLogger(cfg config.Logging) *slog.Logger {
	level := slog.LevelDebug
	switch cfg.Level {
	case logDebug:
		level = slog.LevelDebug
	case logInfo:
		level = slog.LevelInfo
	case logWarn:
		level = slog.LevelWarn
	case logError:
		level = slog.LevelError
	}

	return slog.New(
		slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: level,
			},
		}.NewPrettyHandler(os.Stdout),
	)
}
