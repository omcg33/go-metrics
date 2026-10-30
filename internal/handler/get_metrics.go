package handler

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/omcg33/go-metrics/internal/views"
)

type PageData struct{ Title, Data string }

func (controller *Controller) GetMetrics(res http.ResponseWriter, req *http.Request) {
	metrics := controller.service.Metrics()
	jsonData, err := json.MarshalIndent(metrics, "", "  ")

	if err != nil {
		http.Error(res, "failed to encode metrics", http.StatusInternalServerError)
		return
	}

	tpl := template.Must(template.New("main").Parse(views.Main))

	res.Header().Set("Content-Type", "text/html; charset=utf-8")

	_ = tpl.Execute(res, PageData{Title: "Hello Яндекс практика", Data: string(jsonData)})
}
