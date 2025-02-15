package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	Version   string
	Builder   string
	BuildTime string
)

func init() {
	if Version == "" {
		Version = "development"
	}
}

func BuildInfo(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{
		"version":   Version,
		"builder":   Builder,
		"buildTime": BuildTime,
	})
}
