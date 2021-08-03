package http

import (
	"net/http"
	"undina/domain"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ServiceHandler struct {
	serviceUsecase domain.ServiceUsecase
}

func NewServiceHandler(e *gin.RouterGroup, serviceUsecase domain.ServiceUsecase) {
	handler := &ServiceHandler{
		serviceUsecase: serviceUsecase,
	}

	serviceEndpoints := e.Group("services")
	{
		serviceEndpoints.GET("", handler.GetAllServices)
		serviceEndpoints.GET("/:serviceID", handler.GetServiceDetails)
	}
}

func (s *ServiceHandler) GetAllServices(c *gin.Context) {
	services, err := s.serviceUsecase.FetchAll(c)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": services})
}

func (s *ServiceHandler) GetServiceDetails(c *gin.Context) {
	serviceID := c.Param("serviceID")
	logrus.Debug("serviceID:", serviceID)

	serviceDetails, err := s.serviceUsecase.GetDetailByID(c, serviceID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": serviceDetails})
}
