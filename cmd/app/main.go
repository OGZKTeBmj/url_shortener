package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/OGZKTeBmj/url_shortener/config"
	"github.com/OGZKTeBmj/url_shortener/internal/adapter/postgres"
	"github.com/OGZKTeBmj/url_shortener/internal/adapter/redis"
	"github.com/OGZKTeBmj/url_shortener/internal/controller/http"
	service "github.com/OGZKTeBmj/url_shortener/internal/service/shortener"
	"github.com/OGZKTeBmj/url_shortener/pkg/httpserver"
	"github.com/OGZKTeBmj/url_shortener/pkg/logger"
	"github.com/OGZKTeBmj/url_shortener/pkg/utils"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.InitConfig()
	if err != nil {
		panic(utils.ErrWrap("can't init config", err))
	}

	log := logger.NewSlogLogger(os.Stdout, logger.Config{
		AppName:    cfg.App.Name,
		AppVersion: cfg.App.Version,
		Level:      cfg.Log.Level,
		Env:        cfg.App.Env,
	})

	pgs := postgres.New(cfg.Postgres)
	if err := pgs.Connect(context.Background()); err != nil {
		log.Error("postgres connect", "error", err)
		panic(err)
	}
	defer func() {
		pgs.Close()
		log.Info("postgres stopped")
	}()
	urlStorage := postgres.NewUrlStorage(pgs.DB())

	redisClient := redis.New(cfg.Redis)
	if err := redisClient.Connect(context.Background()); err != nil {
		log.Error("redis connect", "error", err)
		panic(err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Error("redis stopping", "error", err)
		} else {
			log.Info("redis stopped")
		}
	}()

	shortenerService := service.NewShortener(
		log, urlStorage, redisClient,
	)

	router := http.New(log, shortenerService)
	server := httpserver.New(router, httpserver.Config{
		ShutdownTimeout: cfg.HTTP.ShutdownTimeout,
		Port:            cfg.HTTP.Port,
	})

	serverErr := make(chan error, 1)

	go func() {
		err := server.Run()
		if err != nil {
			log.Error("starting http server", "error", err)
			serverErr <- err
		}
	}()
	log.Info("http server started", "port", cfg.HTTP.Port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sign := <-stop:
		log.Info("gracefull shutdown", "signal", sign)
	case err := <-serverErr:
		log.Error("http server", "error", err)
	}

	if err = server.Shutdown(); err != nil {
		log.Error("shutdown http server", "error", err)
	}
}
