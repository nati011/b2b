package handler

import (
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"

	domain_handler "b2b.nati011.github.com/internal/adapter/primary/event/handler/domain"
)

func BuildEventHandler(applicationServices *application_core.Container, domainServices *domain_core.Container) {
	domain_handler.InitProduct(domainServices.ProductService, applicationServices.Event)
}
