package http

import (
	"net/http"

	"undina/auth"
	"undina/domain"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
}

func NewAuthHandler(e *gin.RouterGroup, authUsecase domain.AuthUsecase) {
	handler := &AuthHandler{
		AuthUsecase: authUsecase,
	}

	authEndpoints := e.Group("auth")
	{
		authEndpoints.POST("/signup", handler.SignUp)
		authEndpoints.POST("/signin", handler.SignIn)
	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var body domain.SignupRequest
	if err := c.BindJSON(&body); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := h.AuthUsecase.SignUp(c.Request.Context(), body.Name, body.Email, body.Password)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var body domain.LoginRequest
	if err := c.BindJSON(&body); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := h.AuthUsecase.SignIn(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		if err == auth.ErrUserNotFound {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}
