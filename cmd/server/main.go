package main

import (
	"fmt"
	"net/http"

	"github.com/omcg33/go-metrics/internal/handler"
	"github.com/omcg33/go-metrics/internal/repository"
	"github.com/omcg33/go-metrics/internal/service"
)

func main() {
	storage := repository.NewMemStorage()
	svc := service.NewService(storage)
	controller := handler.NewController(svc)

	mux := http.NewServeMux()

	mux.HandleFunc(`POST /update/{type}/{name}/{value}`, controller.CreateOrUpdateMetric)
	mux.HandleFunc(`GET /metrics`, controller.GetMetrics)

	fmt.Printf("Server listening on http://localhost:8080\n")

	err := http.ListenAndServe(`:8080`, handler.Logging(mux))
	if err != nil {
		fmt.Printf("Server failed\n")
		panic(err)
	}
}
