package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_serviceHandlerHttpDelivery "undina/service/delivery/http"
	_serviceRepo "undina/service/repository/mysql"
	_serviceUsecase "undina/service/usecase"
)

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool("debug") {
		logrus.Info("Service RUN on DEBUG mode")
	}
}

func sayHello(c *gin.Context) {
	c.String(http.StatusOK, "Hello")
}

func sayPongJSON(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {
	logrus.Info("HTTP server started")

	// Database connection
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("loc", "Asia/Taipei")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		logrus.Fatal(err)
	} else {
		logrus.Info("Database connection ESTABLISHED")
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "DELETE", "OPTIONS", "PUT", "PATCH"},
		AllowHeaders: []string{"Authorization", "Content-Type", "Upgrade", "Origin",
			"Connection", "Accept-Encoding", "Accept-Language", "Host",
			"Access-Control-Request-Method", "Access-Control-Request-Headers"},
		AllowCredentials: true,
	}))
	r.GET("/", sayHello)
	r.GET("/ping", sayPongJSON)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	serviceRepo := _serviceRepo.NewmysqlServiceRepository(db)
	serviceUsecase := _serviceUsecase.NewServiceUsecase(serviceRepo, timeoutContext)
	v1Router := r.Group("/api/v1/")
	{
		_serviceHandlerHttpDelivery.NewServiceHandler(v1Router, serviceUsecase)
	}

	logrus.Fatal(r.Run(":" + viper.GetString("server.address")))
}
