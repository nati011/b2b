package rest

import (
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	"b2b.nati011.github.com/internal/core"
)

func BuildRouter(mux *http.ServeMux, services *core.MasterContainer) error {
	print("Building Routes...")
	for _, h := range handler.GetHandlers() {
		print(h)
		if err := h.Init(services); err != nil {
			return err
		}
		h.Routes(mux)
	}

	return nil
}
