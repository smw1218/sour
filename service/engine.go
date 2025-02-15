package service

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smw1218/sour/logger"
	"github.com/smw1218/sour/response"
)

func Engine(log *slog.Logger, svc ServiceInterface) *gin.Engine {
	// gin debug mode isn't really useful
	//gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.NoRoute(NotFoundHandler)
	e.GET(fmt.Sprintf("/%v/version", svc.LongName()), BuildInfo)
	e.GET(fmt.Sprintf("/%v/health", svc.LongName()), HealthHandler)
	e.Use(
		logger.NewLoggerMiddleware(log),
		gin.Recovery())
	// TODO add monitoring etc

	return e
}

func NotFoundHandler(c *gin.Context) {
	response.NotFound(c, "")
}

func HealthHandler(c *gin.Context) {
	c.Data(http.StatusOK, "applicaiton/json", []byte(`{"status":"up"}`))
}
