package service

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/smw1218/sour/logger"
)

func Engine(log *slog.Logger) *gin.Engine {
	e := gin.New()
	e.Use(
		logger.NewLoggerMiddleware(log),
		gin.Recovery())
	// TODO add monitoring etc
	return e
}
