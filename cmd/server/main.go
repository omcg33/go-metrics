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
	repository "github.com/omcg33/go-metrics/internal/repository/metrics"
	"github.com/omcg33/go-metrics/internal/service"
)

func main() {
	logger, logErr := logs.NewLogger("Info")

	if logErr != nil {
		log.Fatal(logErr)
	}
	defer logger.Sync()

	config := config.NewConfig(logger)

	metricsRepository, err := repository.NewMetricsRepository(config, logger)

	if err != nil {
		logger.Error("Failed to init MetricsRepository", zap.Error(err))
		return
	}
	defer func() {
		if err := metricsRepository.Close(); err != nil {
			logger.Error("Failed to close metrics repository", zap.Error(err))
		}
	}()

	svc := service.NewService(metricsRepository)
	controller := handler.NewController(svc)

	if runner, ok := metricsRepository.(repository.Runner); ok {
		go runner.Run()
	}

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
	router.Get("/ping", controller.GetPing(metricsRepository))

	logger.Info("Server listening", zap.String("address", config.ServerAddress))

	errServerInit := http.ListenAndServe(config.ServerAddress, router)

	if errServerInit != nil {
		logger.Error("Server failed", zap.Error(errServerInit))
		return
	}
}
