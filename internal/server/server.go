package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RobinSoGood/EM_test/internal/config"
	"github.com/RobinSoGood/EM_test/internal/logger"
	"github.com/RobinSoGood/EM_test/internal/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	serve    *http.Server
	sService service.SubService
	ErrChan  chan error
}

func New(cfg config.Config, ss service.SubService) *Server {
	addrStr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server := http.Server{
		Addr: addrStr,
	}
	srv := Server{
		serve:    &server,
		sService: ss,
	}
	return &srv
}

func (s *Server) Run(ctx context.Context) error {
	log := logger.Get()
	router := s.configRouting()
	s.serve.Handler = router
	log.Info().Str("addr", s.serve.Addr).Msg("server start")
	if err := s.serve.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("runing server failed")
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.serve.Shutdown(ctx)
}

func (s *Server) configRouting() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) { ctx.String(http.StatusOK, "Hello!") })
	subs := router.Group("/subs")
	{
		subs.GET("/:id", s.getSubHandler)
		subs.GET("/", s.getSubsHandler)
		subs.POST("/add", s.addSubHandler)
		subs.DELETE("/:id", s.deleteSubHandler)
		subs.POST("/total", s.getTotalPriceHandler)
	}
	return router
}
