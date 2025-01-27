package response

import (
	"github.com/gin-gonic/gin"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}

func (err Error) Respond(c *gin.Context) {
	c.AbortWithStatusJSON(err.Code, err)
}
