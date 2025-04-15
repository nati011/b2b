package rest

import (
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	application_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/application"
	domain_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/domain"

	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

func BuildRouter(mux *http.ServeMux, applicationServices *application_core.Container, domainServices *domain_core.Container) error {
	application_handler.InitAuth()
	application_handler.InitHealth()
	application_handler.InitPaymentPartner()

	application_handler.InitResource()
	application_handler.InitUser()
	application_handler.InitTransaction()
	application_handler.InitRole()

	domain_handler.InitRetailer()
	domain_handler.InitDistributor()
	domain_handler.InitProduct()
	domain_handler.InitCatalogue()
	domain_handler.InitCategory()
	domain_handler.InitOrder()
	domain_handler.InitConfigurableProduct()

	for _, h := range handler.GetHandlers() {
		if err := h.Init(applicationServices, domainServices); err != nil {
			return err
		}
		h.Routes(mux)
	}

	return nil
}
