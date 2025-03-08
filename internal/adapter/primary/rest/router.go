package rest

import (
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
)

func BuildRouter(mux *http.ServeMux) error {
	for _, h := range handler.GetHandlers() {
		if err := h.Init(); err != nil {
			return err
		}
		h.Routes(mux)
	}

	return nil
}
