package logger

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type contextLogger string

var contextLoggerKey contextLogger = "logger"

func Get(ctx context.Context) *slog.Logger {
	loggerAny := ctx.Value(contextLoggerKey)
	if loggerAny == nil {
		return slog.Default()
	}
	return loggerAny.(*slog.Logger)
}

func SetContextLogger(ctx context.Context, newLogger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextLoggerKey, newLogger)
}

func SetContextLoggerGin(c *gin.Context, newLogger *slog.Logger) {
	c.Request = c.Request.WithContext(SetContextLogger(c.Request.Context(), newLogger))
}

func Gin(c *gin.Context) *slog.Logger {
	return Get(c.Request.Context())
}
