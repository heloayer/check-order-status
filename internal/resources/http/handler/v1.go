package handler

import (
	sentryGin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/heloayer/check-order-status/config"
	"github.com/heloayer/check-order-status/internal/service"
)

type Server struct {
	conf    *config.Config
	Router  *gin.Engine
	service service.Service
}

func NewServer(
	conf *config.Config,
	router *gin.Engine,
	service service.Service,
) *Server {
	server := &Server{
		conf:    conf,
		Router:  router,
		service: service,
	}
	server.Register(server.Router)

	server.Router.RedirectTrailingSlash = true
	server.Router.RedirectFixedPath = true
	server.Router.HandleMethodNotAllowed = true
	return server
}

func (srv *Server) Register(engine *gin.Engine) {
	engine.Use(sentryGin.New(sentryGin.Options{
		Repanic: true,
	}))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"PUT", "PATCH", "GET", "DELETE", "POST", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	engine.Use(cors.New(corsConfig))

	v1 := engine.Group("/v1")

	order := v1.Group("/order")
	order.GET("/get-order/:id", srv.GetOrder)
	order.POST("/create-order", srv.CreateOrder)

}
