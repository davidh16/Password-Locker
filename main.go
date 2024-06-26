package main

import (
	"password-lock/config"
	"password-lock/controller"
	"password-lock/db"
	mw "password-lock/middleware"
	"password-lock/repository"
	"password-lock/server"
	"password-lock/service"
)

func main() {

	Cfg := config.GetConfig()

	redis := db.ConnectToRedis()

	pgInstance := db.ConnectToDatabase()

	userRepo := repository.NewUserRepository(pgInstance)
	entityRepo := repository.NewEntityRepository(pgInstance)
	tokenRepo := repository.NewTokenRepository(pgInstance)

	svc := service.NewService(
		redis,
		Cfg,
		userRepo,
		entityRepo,
		tokenRepo,
	)

	ctrl := controller.NewController(svc)

	middleware := mw.InitializeMiddleware(pgInstance, redis)

	srv := server.NewServer(ctrl, middleware)

	// Listen and Server in 0.0.0.0:8080
	srv.Run(":8080")
}
