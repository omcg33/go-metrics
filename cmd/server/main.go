package main

import (
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/omcg33/go-metrics/internal/config"
	"github.com/omcg33/go-metrics/internal/handler"
	logs "github.com/omcg33/go-metrics/internal/logger"
	"github.com/omcg33/go-metrics/internal/middleware"
	"github.com/omcg33/go-metrics/internal/repository"
	"github.com/omcg33/go-metrics/internal/service"
)

func main() {
	logger, logErr := logs.NewLogger("Info")

	if logErr != nil {
		log.Fatal(logErr)
	}
	defer logger.Sync()

	config := config.NewConfig(logger)
	storage := repository.NewMemStorage()
	fileStorage, fileStorageErr := repository.NewFileStorage(
		repository.FileStorageConfig{
			FilePath:      config.FileStoragePath,
			Restore:       config.Restore,
			StoreInterval: config.StoreInterval,
		},
		storage,
		logger,
	)
	svc := service.NewService(fileStorage)
	controller := handler.NewController(svc)

	if fileStorageErr != nil {
		logger.Fatal("Failed to init FileStorage", zap.Error(fileStorageErr))
	}

	go fileStorage.AutoSave()

	router := chi.NewRouter()

	router.Use(middleware.Logger(logger))
	router.Use(middleware.DecompressRequest)
	router.Use(middleware.CompressResponse)
	router.Use(chiMiddleware.Recoverer)

	router.Post("/value/", controller.GetMetricJSON)
	router.Post("/update/", controller.CreateOrUpdateMetricJSON)
	router.Post("/update/{type}/{name}/{value}", controller.CreateOrUpdateMetric)
	router.Get("/value/{type}/{name}", controller.GetMetric)
	router.Get("/", controller.GetMetrics)

	logger.Info("Server listening", zap.String("address", config.ServerAddress))

	err := http.ListenAndServe(config.ServerAddress, router)

	if err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
