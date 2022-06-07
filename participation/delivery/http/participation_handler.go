package http

import (
	"net/http"
	"undina/domain"
	"undina/room"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ParticipationHandler struct {
	RoomUsecase domain.RoomUsecase
}

func NewParticipationHandler(e *gin.RouterGroup, authMiddleware gin.HandlerFunc, roomUsecase domain.RoomUsecase) {
	handler := &ParticipationHandler{
		RoomUsecase: roomUsecase,
	}

	participantEndpoints := e.Group("participant", authMiddleware)
	{
		participantEndpoints.DELETE("", handler.LeaveRoom)
		participantEndpoints.PATCH("/status", handler.UpdatePaymentStatus)
	}
}

func (p *ParticipationHandler) LeaveRoom(c *gin.Context) {
	var body domain.ParticipationRequest
	if err := c.BindJSON(&body); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := p.RoomUsecase.LeaveRoom(c, body.RoomId, body.UserId)
	if err != nil {
		logrus.Error(err)
		if err == room.ErrNotHost {
			c.AbortWithStatusJSON(http.StatusForbidden, err.Error())
			return
		}
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

func (p *ParticipationHandler) UpdatePaymentStatus(c *gin.Context) {
	var body domain.ParticipationStatusRequest
	if err := c.BindJSON(&body); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err := p.RoomUsecase.UpdatePaymentStatus(c, body.RoomId, body.UserId, body.PaymentStatus)
	if err != nil {
		logrus.Error(err)
		if err == room.ErrNotAuthorized {
			c.AbortWithStatusJSON(http.StatusForbidden, err.Error())
			return
		}
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}
