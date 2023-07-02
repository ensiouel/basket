package rest

import (
	"github.com/ensiouel/basket/internal/config"
	"github.com/ensiouel/basket/internal/transport/rest/handler"
	"github.com/ensiouel/basket/internal/transport/rest/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type Server struct {
	conf   config.Server
	logger *slog.Logger
	engine *gin.Engine
}

func New(conf config.Server, logger *slog.Logger) *Server {
	engine := gin.Default()

	return &Server{
		conf:   conf,
		logger: logger,
		engine: engine,
	}
}

func (server *Server) Handle(fileHandler *handler.FileHandler) *Server {
	api := server.engine.Group("api", middleware.Error(server.logger))
	{
		v1 := api.Group("v1")
		{
			fileHandler.Register(v1.Group("file"))
		}
	}

	return server
}

func (server *Server) Run() error {
	return server.engine.Run(server.conf.Addr)
}
