package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	models "github.com/omcg33/go-metrics/internal/model"
)

var validate = newValidate()

func newValidate() *validator.Validate {
	v := validator.New()
	v.RegisterAlias("metric_type", fmt.Sprintf("oneof=%s %s", models.Gauge, models.Counter))
	return v
}