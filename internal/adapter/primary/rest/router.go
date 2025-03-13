package rest

import (
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	application_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/application"
	"b2b.nati011.github.com/internal/core"
)

func BuildRouter(mux *http.ServeMux, services *core.MasterContainer) error {
	application_handler.InitAuth()
	application_handler.InitDistributor()
	application_handler.InitResource()

	for _, h := range handler.GetHandlers() {
		if err := h.Init(services); err != nil {
			return err
		}
		h.Routes(mux)
	}

	return nil
}
