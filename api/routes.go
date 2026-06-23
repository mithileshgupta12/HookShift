package api

import (
	"net/http"

	"github.com/mithileshgupta12/hook-shift/queue"
)

func InitializeRoutes(mux *http.ServeMux, q queue.Queue) {
	handler := NewHandler(q)

	mux.HandleFunc("GET /healthz", handler.handleHealthz)
	mux.HandleFunc("POST /v1/dispatches", handler.handleDispatches)
}
