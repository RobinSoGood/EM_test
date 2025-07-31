package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/RobinSoGood/EM_test/internal/config"
	"github.com/RobinSoGood/EM_test/internal/logger"
	"github.com/RobinSoGood/EM_test/internal/server"
	"github.com/RobinSoGood/EM_test/internal/service"
	"github.com/RobinSoGood/EM_test/internal/storage"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"golang.org/x/sync/errgroup"
)

// @title           Subscription API With Swagger
// @version         1.0

// @host      localhost:8081
// @BasePath  /api/v1

func main() {
	cfg := config.ReadConfig()
	log := logger.Get(cfg.Debug)
	log.Debug().Any("cfg", cfg).Msg("config")

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-c

		log.Info().Msg("gracefully stopping...")
		cancel()
	}()

	var subService service.SubService

	err := storage.Migrations(cfg.DbDSN, cfg.MigratePath)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	stor, err := storage.NewDB(context.Background(), cfg.DbDSN)
	if err != nil {
		log.Error().Err(err).Send()
	} else {
		subService = service.NewSubService(stor)
	}
	serve := server.New(cfg, subService)

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		if err := serve.Run(gCtx); err != nil {
			log.Error().Err(err).Send()
			return err
		}
		return nil
	})
	group.Go(func() error {
		log.Debug().Msg("start listening error channel")
		defer log.Debug().Msg("stop listening error channel")
		return <-serve.ErrChan
	})
	group.Go(func() error {
		<-gCtx.Done()
		return serve.Shutdown(gCtx)
	})
	group.Go(func() error {
		<-gCtx.Done()
		return stor.Close()
	})

	if err := group.Wait(); err != nil {
		log.Error().Err(err).Send()
	}
	log.Info().Msg("server stoped")
}
