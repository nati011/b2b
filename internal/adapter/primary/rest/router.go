package rest

import (
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	application_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/application"
	application_resource_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/application/resource"
	application_role_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/application/role"
	application_transaction_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/application/transaction"
	application_user_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/application/user"
	application_core "b2b.nati011.github.com/internal/core/application"

	// application_email_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/application/email"
	domain_category_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/domain/category"
	domain_configurable_product_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/domain/configurable_product"
	domain_order_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/domain/order"
	domain_product_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/domain/product"
	domain_domain_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/domain/retailer"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

func BuildRouter(mux *http.ServeMux, applicationServices *application_core.Container, domainServices *domain_core.Container) error {
	application_handler.InitAuth()
	application_handler.InitHealth()
	application_handler.InitPaymentPartner()

	application_resource_handler.InitResource()
	application_user_handler.InitUser()
	application_transaction_handler.InitTransaction()
	application_role_handler.InitRole()

	domain_domain_handler.InitRetailer()
	domain_product_handler.InitProduct()
	domain_category_handler.InitCategory()
	domain_order_handler.InitOrder()
	domain_configurable_product_handler.InitConfigurableProduct()

	for _, h := range handler.GetHandlers() {
		if err := h.Init(applicationServices, domainServices); err != nil {
			return err
		}
		h.Routes(mux)
	}

	return nil
}
