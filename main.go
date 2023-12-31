package main

import (
	"password-lock/controller"
	"password-lock/db"
	"password-lock/repository"
	"password-lock/server"
	"password-lock/service"

	"github.com/gin-gonic/gin"
)

func main() {

	redis := db.ConnectToRedis()

	pgInstance := db.ConnectToDatabase()

	userRepo := repository.NewUserRepository(pgInstance)
	entityRepo := repository.NewEntityRepository(pgInstance)

	svc := service.NewService(
		redis,
		userRepo,
		entityRepo,
	)

	ctrl := controller.NewController(svc)
	r := gin.Default()

	srv := server.NewServer(r, ctrl, redis)

	// Listen and Server in 0.0.0.0:8080
	srv.Run(":8080")
}
