package http

import (
	"net/http"
	"undina/application"
	"undina/domain"
	"undina/participation"
	"undina/room"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ApplicationHandler struct {
	ApplicationUsecase   domain.ApplicationUsecase
	ParticipationUsecase domain.ParticipationUsecase
	RoomUsecase          domain.RoomUsecase
}

func NewApplicationHandler(e *gin.RouterGroup, authMiddleware gin.HandlerFunc, applicationUsecase domain.ApplicationUsecase, participationUsecase domain.ParticipationUsecase, roomUsecase domain.RoomUsecase) {
	handler := &ApplicationHandler{
		ApplicationUsecase:   applicationUsecase,
		ParticipationUsecase: participationUsecase,
		RoomUsecase:          roomUsecase,
	}

	applicationEndpoints := e.Group("application", authMiddleware)
	{
		applicationEndpoints.POST("/accept", handler.AcceptApplication)
		applicationEndpoints.DELETE("/delete", handler.DeleteApplication)
	}
}

func (a *ApplicationHandler) AcceptApplication(c *gin.Context) {
	var body domain.ApplicationRequest
	if err := c.BindJSON(&body); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isAdmin, err := a.ParticipationUsecase.IsAdmin(c, int32(body.RoomId))
	if !isAdmin || err != nil {
		logrus.Error(err)
		if err == participation.ErrNotHost {
			c.AbortWithStatusJSON(http.StatusForbidden, err.Error())
			return
		}
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	roomInfo, _ := a.RoomUsecase.GetRoomInfo(c, int32(body.RoomId))
	members, _ := a.RoomUsecase.GetRoomMembers(c, int32(body.RoomId))
	if len(members) >= int(roomInfo.MaxCount) {
		c.AbortWithStatusJSON(http.StatusForbidden, room.ErrRoomFull.Error())
		return
	}

	err = a.ApplicationUsecase.AcceptApplication(c, body.RoomId, body.UserId)
	if err != nil {
		logrus.Error(err)
		if err == application.ErrApplicationNotFound {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = a.ParticipationUsecase.Create(c, &domain.Participation{
		UserId:        body.UserId,
		RoomId:        body.RoomId,
		PaymentStatus: domain.UNPAID,
		IsHost:        false,
	})
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

func (a *ApplicationHandler) DeleteApplication(c *gin.Context) {
	var body domain.ApplicationRequest
	if err := c.BindJSON(&body); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isAdmin, err := a.ParticipationUsecase.IsAdmin(c, int32(body.RoomId))
	if !isAdmin || err != nil {
		logrus.Error(err)
		if err == participation.ErrNotHost {
			c.AbortWithStatusJSON(http.StatusForbidden, err.Error())
			return
		}
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = a.ApplicationUsecase.DeleteApplication(c, body.RoomId, body.UserId)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}
