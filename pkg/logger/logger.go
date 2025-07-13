package logger

import (
	"time"

	zapsentry "github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New создает zap логгер, настраивая уровень и интеграцию с Sentry в зависимости
// от окружения. Возвращает *zap.SugaredLogger.
func New(appEnv, sentryDSN string) *zap.SugaredLogger {
	var l *zap.Logger
	var err error
	if appEnv == "prod" {
		l, err = zap.NewProduction()
	} else {
		l, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}

	// Интеграция Sentry (опционально)
	if sentryDSN != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:         sentryDSN,
			Environment: appEnv,
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

	return l.Sugar()
}

// Sync очищает буферы zap и отправляет накопленные события в Sentry.
func Sync(l *zap.SugaredLogger) {
	if l != nil {
		_ = l.Sync()
	}
	sentry.Flush(2 * time.Second)
}
