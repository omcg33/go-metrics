package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/omcg33/go-metrics/internal/config"
	"github.com/omcg33/go-metrics/internal/handler"
	"github.com/omcg33/go-metrics/internal/repository"
	"github.com/omcg33/go-metrics/internal/service"
)


func main() {
	config := config.NewConfig();
	storage := repository.NewMemStorage()
	svc := service.NewService(storage)
	controller := handler.NewController(svc)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/update/{type}/{name}/{value}", controller.CreateOrUpdateMetric)
	router.Get("/value/{type}/{name}", controller.GetMetric)
	router.Get("/", controller.GetMetrics)

	fmt.Printf("Server listening on http://%s\n", *config.ServerAddress)

	err := http.ListenAndServe(*config.ServerAddress, router)
	if err != nil {
		fmt.Printf("Server failed\n")
		panic(err)
	}
}
