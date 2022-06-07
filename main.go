package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_authHandlerHttpDelivery "undina/auth/delivery/http"
	_authUsecase "undina/auth/usecase"

	_userHandlerHttpDelivery "undina/user/delivery/http"
	_userRepo "undina/user/repository/mysql"
	_userUsecase "undina/user/usecase"

	_serviceHandlerHttpDelivery "undina/service/delivery/http"
	_serviceRepo "undina/service/repository/mysql"
	_serviceUsecase "undina/service/usecase"

	_roomHandlerHttpDelivery "undina/room/delivery/http"
	_roomRepo "undina/room/repository/mysql"
	_roomUsecase "undina/room/usecase"
)

func init() {
	viper.SetConfigFile("./config/config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool("debug") {
		logrus.Info("Service RUN on DEBUG mode")
	}
}

func sayHello(c *gin.Context) {
	version := os.Getenv("BUILD_VERSION")
	s := fmt.Sprintf("Hello, backend version: %s", version)
	c.String(http.StatusOK, s)
}

func sayPongJSON(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong!",
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

	userRepo := _userRepo.NewmysqlUserRepository(db)
	userUsecase := _userUsecase.NewUserUsecase(userRepo, timeoutContext)

	serviceRepo := _serviceRepo.NewmysqlServiceRepository(db)
	serviceUsecase := _serviceUsecase.NewServiceUsecase(serviceRepo, timeoutContext)

	roomRepo := _roomRepo.NewmysqlRoomRepository(db)
	roomUsecase := _roomUsecase.NewRoomUsecase(roomRepo, timeoutContext)

	authUsecase := _authUsecase.NewAuthUsecase(
		userRepo,
		viper.GetString("auth.hash_salt"),
		[]byte(viper.GetString("auth.signing_key")),
		viper.GetDuration("auth.token_ttl"),
	)
	authMiddleware := _authHandlerHttpDelivery.NewAuthMiddleware(authUsecase)

	v1Router := r.Group("/api/v1/")
	{
		_authHandlerHttpDelivery.NewAuthHandler(v1Router, authUsecase)
		_userHandlerHttpDelivery.NewUserHandler(v1Router, authMiddleware, userUsecase)
		_serviceHandlerHttpDelivery.NewServiceHandler(v1Router, serviceUsecase)
		_roomHandlerHttpDelivery.NewRoomHandler(v1Router, authMiddleware, roomUsecase)
	}

	logrus.Fatal(r.Run(":" + viper.GetString("server.address")))
}
