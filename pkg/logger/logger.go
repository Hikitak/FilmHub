package logger

import (
	"os"
	"time"

	zapsentry "github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

// Init инициализирует глобальный zap логгер.
// В режиме APP_ENV=prod используется Production конфигурация,
// иначе Development (более человекочитаемая).
func Init() {
	env := os.Getenv("APP_ENV")
	dsn := os.Getenv("SENTRY_DSN")

	var l *zap.Logger
	var err error
	if env == "prod" {
		l, err = zap.NewProduction()
	} else {
		l, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}

	// Инициализация Sentry (если задан DSN)
	if dsn != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:         dsn,
			Environment: env,
		}); err != nil {
			l.Warn("Sentry initialization failed", zap.Error(err))
		} else {
			cfg := zapsentry.Configuration{
				Level:             zapcore.ErrorLevel,
				EnableBreadcrumbs: true,
				BreadcrumbLevel:   zapcore.InfoLevel,
			}
			core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(sentry.CurrentHub().Client()))
			if err == nil {
				l = zapsentry.AttachCoreToLogger(core, l)
			} else {
				l.Warn("Failed to attach Sentry zap core", zap.Error(err))
			}
		}
	}
	Log = l.Sugar()
}

// Sync очищает буферы zap и отправляет накопленные события в Sentry.
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
	sentry.Flush(2 * time.Second)
}
