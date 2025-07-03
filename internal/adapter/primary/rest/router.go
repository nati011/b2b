package rest

import (
	"context"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	application_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/application"
	domain_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/domain"

	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/resource"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

func BuildRouter(mux *http.ServeMux, applicationServices *application_core.Container, domainServices *domain_core.Container) error {
	application_handler.InitAuth()
	application_handler.InitHealth()
	application_handler.InitPaymentPartner()
	application_handler.InitMobileClient()
	application_handler.InitEmailTemplate()
	application_handler.InitEmail()
	application_handler.InitResource()
	application_handler.InitUser()
	application_handler.InitTransaction()
	application_handler.InitRole()
	application_handler.InitIdentity()

	domain_handler.InitDistributor()
	domain_handler.InitRetailer()
	domain_handler.InitProduct()
	domain_handler.InitCatalogue()
	domain_handler.InitCategory()
	domain_handler.InitOrder()
	domain_handler.InitPayment()
	domain_handler.InitConfigurableProduct()
	domain_handler.InitInvoice()

	for _, h := range handler.GetHandlers() {
		if err := h.Init(applicationServices.AuthMiddleware, applicationServices, domainServices); err != nil {
			return err
		}
		h.Routes(mux)
	}

	ctx := context.Background()
	for _, r := range handler.GetRoutes() {
		_, err := applicationServices.ResourceService.GetByName(ctx, r.Name)
		if err != nil {
			switch err {
			case resource.ErrNameNotFound:
				continue
			default:
				panic("Failed to fetch resource")
			}
		}
		_, err = applicationServices.ResourceService.Create(ctx, &resource.CreateRequest{
			Name:     r.Name,
			Resource: r.Resource,
			Action:   resource.ANY,
		})
		if err != nil {
			if _, err = applicationServices.ResourceService.Create(ctx,
				&resource.CreateRequest{
					Name: r.Name, Resource: r.Resource, Action: resource.ANY}); err != nil {
				panic("Failed to create resource err")
			}
		}
	}
	return nil
}
