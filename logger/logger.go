package logger

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmittmann/tint"
)

func SetupDefaultSlog(serviceName string) {
	slogger := slog.New(
		tint.NewHandler(log.Default().Writer(),
			&tint.Options{
				Level:      slog.LevelDebug,
				TimeFormat: "2006-01-02 15:04:05.000",
			}))
	slogger = slogger.With("svc", serviceName)
	slog.SetDefault(slogger)
}

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

func NewLoggerMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		handlerLogger := log.With("hand", HandlerName(c.HandlerName()))
		SetContextLoggerGin(c, handlerLogger)
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		endpoint := fmt.Sprintf("%v %v", c.Request.Method, c.Request.URL)
		handlerLogger.Info(endpoint, "tm", latency)
	}
}

// HandlerName makes the handler names shorter and readable ala:
// github.com/smw1218/sour/cmd/test-service/te.(*TEHandlers).Test-fm -> te.Test
// github.com/smw1218/sour/cmd/test-service/te.TopLevel -> te.TopLevel
func HandlerName(fullType string) string {
	slashed := strings.Split(fullType, "/")
	last := slashed[len(slashed)-1]
	last = strings.TrimSuffix(last, "-fm")
	dotted := strings.Split(last, ".")
	if len(dotted) <= 2 {
		return last
	}
	return fmt.Sprintf("%s.%s", dotted[0], strings.Join(dotted[2:], "."))
}
