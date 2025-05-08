package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"proxy-data-filter/internal/config"
	"proxy-data-filter/internal/handler"
)

type Init struct {
	cfg    *config.Config
	router *httprouter.Router
}

func New(cfg *config.Config, router *httprouter.Router) *Init {
	return &Init{
		cfg:    cfg,
		router: router,
	}
}

var conf = []map[string]interface{}{
	{
		"path":   "/online/companies/jQu",
		"method": "POST",
	},
}

func (controller *Init) SetRoutes(ctx context.Context) error {
	handler.InitHandler(controller.cfg)

	for _, item := range conf {
		method, ok := item["method"].(string)
		if !ok {
			return fmt.Errorf("method is not a string")
		}
		path, ok := item["path"].(string)
		if !ok {
			return fmt.Errorf("path is not a string")
		}
		controller.router.HandlerFunc(method, path, Handle)
	}

	return nil
}

func Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(r.Body); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}
