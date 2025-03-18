package rest

import (
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	application_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/application"
	application_core "b2b.nati011.github.com/internal/core/application"

	domain_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/domain"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

func BuildRouter(mux *http.ServeMux, applicationServices *application_core.Container, domainServices *domain_core.Container) error {

	application_handler.InitResource()
	domain_handler.InitProduct()

	for _, h := range handler.GetHandlers() {
		if err := h.Init(applicationServices, domainServices); err != nil {
			return err
		}
		h.Routes(mux)
	}

	return nil
}
