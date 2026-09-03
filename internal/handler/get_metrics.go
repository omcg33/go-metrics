package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type PageData struct{ Title, Data string }

func (controller *Controller) GetMetrics(res http.ResponseWriter, req *http.Request) {
	metrics := controller.service.Metrics()
	jsonData, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		http.Error(res, "failed to encode metrics", http.StatusInternalServerError)
		return
	}

	tpl := template.Must(template.ParseFiles("./internal/views/main.html"))
	_ = tpl.Execute(res, PageData{Title: "Hello Яндекс практика", Data: string(jsonData)})
}