package event

import (
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"

	domain_handler "b2b.nati011.github.com/internal/adapter/primary/event/handler/domain"
)

func BuildEventHandler(applicationServices *application_core.Container, domainServices *domain_core.Container) error {
	err := domain_handler.InitProduct(domainServices.ProductService, applicationServices.Event)
	if err != nil {
		return err
	}

	err = domain_handler.InitAuth(domainServices.RetailerService, applicationServices.UserService, applicationServices.Event)
	if err != nil {
		return err
	}
	return nil
}
