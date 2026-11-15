package handler

import (
	"context"
	"net/http"
	"time"
)

func (controller *Controller) GetPing(dbStorage DBStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), time.Second)
		defer cancel()

		err := dbStorage.Ping(ctx)

		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}
