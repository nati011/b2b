package core

import (
	"database/sql"

	"b2b.nati011.github.com/config"
	payment_db_port "b2b.nati011.github.com/internal/adapter/secondary/application/payment/db"
	category_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	configurable_product_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/configurable_product/db"
	distributor_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/distributor/db"
	invoice_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/invoice/db"
	order_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/order/db"
	product_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/product/db"
	retailer_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/retailer/db"
	application_core "b2b.nati011.github.com/internal/core/application"
	checkout "b2b.nati011.github.com/internal/core/application/checkout"
	payment_verification "b2b.nati011.github.com/internal/core/application/payment_verification"
	"b2b.nati011.github.com/internal/core/domain/catalogue"
	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

/*Dependency Tree*/

// Container
// ├── CategoryService
// │   └── ProductService
// │     └── ConfigurableProductService
// │
// ├──── ProductService
// │ └── ConfigurableProductService
// │   		└── CatalogueService
// │
// ├── InvoiceService
// │   └── OrderService
// │
// ├── RetailerService
// │   └── OrderService
// │
// ├── ProductService
// │   └── OrderService

type Container struct {
	db                         *sql.DB
	CategoryService            category.Provider
	ProductService             product.Provider
	ConfigurableProductService configurable_product.Provider
	InvoiceService             invoice.Provider
	OrderService               order.Provider
	DistributorService         distributor.Provider
	RetailerService            retailer.Provider
	ApplicationServices        application_core.Container
	PaymentVerificationService payment_verification.Provider
	CheckoutService            checkout.Provider
	Pagination                 config.Pagination
	CatalogueService           catalogue.Provider

	//REMOVE ME FROM HERE
	FrontendURL string
	BaseURL     string
}

func NewContainer(application_core application_core.Container, baseUrl string, frontendUrl string, db *sql.DB) *Container {
	container := Container{}
	container.db = db
	container.ApplicationServices = application_core
	container.FrontendURL = frontendUrl
	container.BaseURL = baseUrl

	// ORDER ORDER!!

	//  messing up the order creates chaos
	// //utils
	container.InitPagination()
	container.InitCategoryService()
	container.InitProductService()
	container.InitConfigrableProductService()
	container.InitInvoiceService()
	container.InitRetailerService()
	container.InitOrderService()
	container.InitDistributorService()
	container.InitPaymentVerificationService()
	container.InitCheckoutService()
	container.InitOrderService()
	container.InitCatalogueService()
	return &container
}

func (m *Container) InitPagination() {
	m.Pagination = *config.NewPaginationBuilder().Build()
}

func (m *Container) InitCategoryService() {
	m.CategoryService = category.NewCategory(category_db_port.NewPostgres(m.db))
}

func (m *Container) InitCatalogueService() {
	m.CatalogueService = catalogue.NewCatalogueService(m.ProductService, m.ConfigurableProductService)
}

func (m *Container) InitProductService() {
	m.ProductService = product.NewProduct(product_db_port.NewPostgres(m.db, &m.Pagination), m.CategoryService)
}

func (m *Container) InitConfigrableProductService() {
	m.ConfigurableProductService = configurable_product.NewConfigurableProductService(configurable_product_db_port.NewPostgres(m.db),
		m.ProductService,
	)
}

func (m *Container) InitInvoiceService() {
	m.InvoiceService = invoice.NewInvoice(invoice_db_port.NewPostgres(m.db, &m.ApplicationServices.Pagination))
}

func (m *Container) InitOrderService() {
	m.OrderService = order.NewOrderService(order_db_port.NewPostgres(m.db, &m.Pagination), m.InvoiceService, m.ProductService, m.RetailerService, m.CheckoutService)
}

func (m *Container) InitDistributorService() {
	m.DistributorService = distributor.NewDistributorService(m.ApplicationServices.UserService, distributor_db_port.NewPostgres(m.db, &m.ApplicationServices.Pagination))
}

func (m *Container) InitRetailerService() {
	m.RetailerService = retailer.NewRetailerService(m.ApplicationServices.UserService, retailer_db_port.NewPostgres(m.db, &m.ApplicationServices.Pagination))
}

func (m *Container) InitPaymentVerificationService() {
	m.PaymentVerificationService = payment_verification.NewPaymentVerificationService(payment_db_port.NewPostgres(m.db), m.ApplicationServices.PaymentPartnerService, m.ApplicationServices.TransactionService, m.OrderService)
}

func (m *Container) InitCheckoutService() {
	m.CheckoutService = checkout.NewCheckoutService(payment_db_port.NewPostgres(m.db), m.ApplicationServices.PaymentPartnerService, m.ApplicationServices.TransactionService, m.FrontendURL, m.BaseURL)
}
