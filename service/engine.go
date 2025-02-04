package service

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/smw1218/sour/logger"
	"github.com/smw1218/sour/response"
)

func Engine(log *slog.Logger) *gin.Engine {
	// gin debug mode isn't really useful
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.NoRoute(NotFoundHandler)
	e.Use(
		logger.NewLoggerMiddleware(log),
		gin.Recovery())
	// TODO add monitoring etc

	return e
}

func NotFoundHandler(c *gin.Context) {
	response.NotFound(c, "")
}
