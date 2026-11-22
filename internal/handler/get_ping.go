package handler

import (
	"context"
	"net/http"
	"time"
)

func (controller *Controller) GetPing(repository any) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), time.Second)
		defer cancel()

		pinger, ok := repository.(interface{ Ping(context.Context) error })

		if !ok {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		err := pinger.Ping(ctx)

		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}
