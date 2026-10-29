package agent

import (
	"strings"

	"go.uber.org/zap"
	"resty.dev/v3"

	logs "github.com/omcg33/go-metrics/internal/logger"
	models "github.com/omcg33/go-metrics/internal/model"
)

var _ Service = (*MetricsService)(nil)

type MetricsService struct {
	client *resty.Client
	logger *logs.Logger
}

func NewService(serverAddress string, logger *logs.Logger) *MetricsService {
	if !strings.Contains(serverAddress, "://") {
		serverAddress = "http://" + serverAddress
	}

	client := resty.New().
		SetBaseURL(serverAddress).
		SetRequestMiddlewares(
			resty.MiddlewareRequestCreate,
			gzipMiddleware,
		).
		SetCloseConnection(true)

	return &MetricsService{
		client: client,
		logger: logger,
	}
}

func (service *MetricsService) Close() error {
	return service.client.Close()
}

func (service *MetricsService) Report(report Report) {
	for name, value := range report.gauges {
		service.logger.Info("Send Gauge metric POST to /update", zap.String("name", name), zap.Float64("value", value))

		gaugeValue := value
		_, err := service.client.R().
			SetBody(models.Metrics{
				ID:    name,
				MType: models.Gauge,
				Value: &gaugeValue,
			}).
			SetHeader("Content-Type", "application/json").
			Post("/update/")

		if err != nil {
			service.logger.Error("Failed Gauge metric POST to /update", zap.String("name", name), zap.Float64("value", value), zap.Error(err))
			continue
		}
	}

	for name, value := range report.counters {

		service.logger.Info("Send Counter metric POST to /update", zap.String("name", name), zap.Int64("value", value))

		counterValue := value
		_, err := service.client.R().
			SetBody(models.Metrics{
				ID:    name,
				MType: models.Counter,
				Delta: &counterValue,
			}).
			SetHeader("Content-Type", "application/json").
			Post("/update/")

		if err != nil {
			service.logger.Error("Failed Counter metric %s === %d POST to /update with %v", zap.String("name", name), zap.Int64("value", value), zap.Error(err))
			continue
		}
	}
}
