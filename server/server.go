package server

import (
	"github.com/gin-gonic/gin"
	"password-lock/controller"
	mw "password-lock/middleware"
)

type Server struct {
	controller *controller.Controller
	router     *gin.Engine
}

func NewServer(ctrl *controller.Controller, middleware *mw.Middleware) Server {

	router := gin.Default()
	router.Use(middleware.CORS())

	initializeRoutes(router, ctrl, middleware)

	return Server{
		controller: ctrl,
		router:     router,
	}
}

func (s Server) Run(port string) {
	err := s.router.Run(port)
	if err != nil {
		return
	}
}
