package http

import (
	"net/http"
	"strconv"
	"undina/domain"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

func NewUserHandler(e *gin.RouterGroup, authMiddleware gin.HandlerFunc, userUsecase domain.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: userUsecase,
	}

	userEndpoints := e.Group("user", authMiddleware)
	{
		userEndpoints.GET("", handler.GetUser)
	}
}

func (u *UserHandler) GetUser(c *gin.Context) {
	user := c.Value(domain.CtxUserKey).(*domain.User)
	userId := strconv.Itoa(int(user.Id))

	res, err := u.UserUsecase.GetByID(c, userId)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, res)
}
