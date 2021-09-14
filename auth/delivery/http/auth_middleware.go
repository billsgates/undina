package http

import (
	"net/http"
	"strings"

	"undina/auth"
	"undina/domain"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	AuthUsecase domain.AuthUsecase
}

func NewAuthMiddleware(authUsecase domain.AuthUsecase) gin.HandlerFunc {
	return (&AuthMiddleware{
		AuthUsecase: authUsecase,
	}).Handle
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if headerParts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := m.AuthUsecase.ParseToken(c.Request.Context(), headerParts[1])
	if err != nil {
		status := http.StatusInternalServerError
		if err == auth.ErrInvalidAccessToken {
			status = http.StatusUnauthorized
		}

		c.AbortWithStatus(status)
		return
	}

	c.Set(domain.CtxUserKey, user)
}
