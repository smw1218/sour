package authz

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/smw1218/sour/response"
)

type authContextType string

var tokenKey authContextType = "token"

func NewAuthMiddleware(ac *AuthChecker) gin.HandlerFunc {
	am := &authMiddleware{ac}
	return am.Check
}

type authMiddleware struct {
	ac *AuthChecker
}

func (am *authMiddleware) Check(c *gin.Context) {
	tokString, err := request.BearerExtractor{}.ExtractToken(c.Request)
	if err != nil {
		response.NewError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)).Respond(c)
		return
	}
	token, err := am.ac.CheckAuth(tokString)
	if err != nil {
		response.NewError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)).Respond(c)
		return
	}
	ctx := context.WithValue(c.Request.Context(), tokenKey, token)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
