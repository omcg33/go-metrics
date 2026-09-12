package main

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/omcg33/go-metrics/internal/config"
	"github.com/omcg33/go-metrics/internal/handler"
	logger "github.com/omcg33/go-metrics/internal/logger"
	"github.com/omcg33/go-metrics/internal/middleware"
	"github.com/omcg33/go-metrics/internal/repository"
	"github.com/omcg33/go-metrics/internal/service"
)

func main() {
	if err := logger.Initialize("Info"); err != nil {
		panic(err)
	}

	config := config.NewConfig()
	storage := repository.NewMemStorage()
	svc := service.NewService(storage)
	controller := handler.NewController(svc)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(chiMiddleware.Recoverer)

	router.Post("/update", controller.CreateOrUpdateMetric)
	router.Post("/update/{type}/{name}/{value}", controller.CreateOrUpdateMetric)
	router.Get("/value/{type}/{name}", controller.GetMetric)
	router.Get("/", controller.GetMetrics)

	logger.Log.Info("Server listening", zap.String("address", config.ServerAddress))

	err := http.ListenAndServe(config.ServerAddress, router)
	if err != nil {
		logger.Log.Info("Server failed")
		panic(err)
	}
}
