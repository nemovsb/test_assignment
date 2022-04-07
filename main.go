package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"test_assignment/internal/configuration/cfg"
	server "test_assignment/internal/http_server"
	"test_assignment/internal/http_server/ginhandlers"
	"test_assignment/internal/storage"
	"test_assignment/pkg/zaplogger"

	group "github.com/oklog/run"
	errors2 "github.com/pkg/errors"
	"go.uber.org/zap"
)

var ErrOsSignal = errors.New("got os signal")

func main() {

	config, err := cfg.ViperConfigurationProvider(os.Getenv("GOLANG_ENVIRONMENT"), false)
	if err != nil {
		log.Fatal("Read config error: ", err)
	}

	logger, zapLoggerCleanup, err := zaplogger.Provider(config.ZapLoggerMode)
	if err != nil {
		log.Fatal(errors2.WithMessage(err, "zap logger provider"))
	}

	logger.Info("application", zap.String("event", "initializing"))
	logger.Info("application", zap.Any("resolved_configuration", config))

	db, err := storage.NewPGDB(cfg.GetDBConfig(config), time.Duration(config.TTL))
	if err != nil {
		logger.Sugar().Fatalf("No DB conn: %s", err)
	}

	cache, err := storage.NewRedisStorage(cfg.GetCacheConfig(config), time.Duration(config.TTL), logger)
	if err != nil {
		logger.Error("No cache connection: ", zap.Error(err))
	}

	siteHandler := ginhandlers.NewSiteHandler(config.TTL, config.Timeout, cache, db, logger)
	handlerSet := server.NewHandlerSet(siteHandler)
	router := server.NewRouter(handlerSet)

	server := server.NewServer(cfg.GetHTTPServerConfig(config), router)

	var (
		serviceGroup        group.Group
		interruptionChannel = make(chan os.Signal, 1)
	)

	serviceGroup.Add(func() error {
		signal.Notify(interruptionChannel, syscall.SIGINT, syscall.SIGTERM)
		osSignal := <-interruptionChannel

		return fmt.Errorf("%w: %s", ErrOsSignal, osSignal)
	}, func(error) {
		interruptionChannel <- syscall.SIGINT
	})

	serviceGroup.Add(func() error {
		logger.Info("server", zap.String("event", "HTTP API started"))

		return server.Run()
	}, func(error) {
		err = server.Shutdown()
		logger.Info("shutdown Http Server error", zap.Error(err))
	})

	err = serviceGroup.Run()
	cache.Client.Close()
	logger.Info("services stopped", zap.Error(err))

	zapLoggerCleanup()

}
