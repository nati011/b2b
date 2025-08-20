package core

import (
	"database/sql"

	category_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	config_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/domain/config"
	configurable_product_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/configurable_product/db"
	distributor_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/distributor/db"
	distributor_approval_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/distributor_approval/db"
	invoice_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/invoice/db"
	order_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/order/db"
	product_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/product/db"
	retailer_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/retailer/db"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/domain/catalogue"
	"b2b.nati011.github.com/internal/core/domain/category"
	config_module "b2b.nati011.github.com/internal/core/domain/config"
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	distributorApproval "b2b.nati011.github.com/internal/core/domain/distributor_approval"
	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/payment_verification"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

/*Dependency Tree*/

// Container
// ├── CategoryService
// │   └── ProductService
// │
// ├── ConfigurableProductService
// │	└── ProductService
// │
// ├──── CatalogueService
// │   	 └── ProductService
// │
// ├──── ConfigurableProductService
// │   	 └── CatalogueService
// │
// ├── InvoiceService
// │   └── OrderService
// │
// ├── RetailerService
// │   └── OrderService
// │
// ├── ProductService
// │   └── OrderService
// │
// ├── DistributorService
// │   └── DistributorApprovalService

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
	CatalogueService           catalogue.Provider
	PaymentVerificationService payment_verification.Provider
	DistributorApprovalService distributorApproval.Provider
	ConfigService              config_module.Provider
}

func NewContainer(application_core application_core.Container, baseUrl string, frontendUrl string, db *sql.DB) *Container {
	container := Container{}
	container.db = db
	container.ApplicationServices = application_core

	// ORDER ORDER!!
	//  messing up the order creates chaos
	container.InitCategoryService()
	container.InitConfigrableProductService()
	container.InitInvoiceService()
	container.InitRetailerService()
	container.InitDistributorApprovalService()
	container.InitDistributorService()
	container.InitProductService()
	container.InitConfigService()
	container.InitOrderService()
	container.InitCatalogueService()
	container.InitPaymentVerificationService()

	return &container
}

func (m *Container) InitCategoryService() {
	m.CategoryService = category.NewCategory(
		category_db_port.NewPostgres(m.db))
}

func (m *Container) InitCatalogueService() {
	m.CatalogueService = catalogue.NewCatalogueService(
		m.ProductService,
		m.ConfigurableProductService)
}

func (m *Container) InitProductService() {
	m.ProductService = product.NewProduct(
		product_db_port.NewPostgres(m.db, m.ApplicationServices.Pagination),
		m.CategoryService,
		m.DistributorService)
}

func (m *Container) InitConfigrableProductService() {
	m.ConfigurableProductService = configurable_product.NewConfigurableProductService(
		configurable_product_db_port.NewPostgres(m.db),
		m.ProductService,
	)
}

func (m *Container) InitInvoiceService() {
	m.InvoiceService = invoice.NewInvoice(
		invoice_db_port.NewPostgres(
			m.db,
			m.ApplicationServices.Pagination),
		m.ApplicationServices.RenderService,
	)
}

func (m *Container) InitDistributorService() {
	m.DistributorService = distributor.NewDistributorService(
		m.ApplicationServices.UserService,
		distributor_db_port.NewPostgres(m.db, m.ApplicationServices.Pagination),
		m.DistributorApprovalService,
		&m.ApplicationServices.Event)
}

func (m *Container) InitDistributorApprovalService() {
	m.DistributorApprovalService = distributorApproval.NewDistributorApprovalService(
		distributor_approval_db_port.NewPostgres(m.db))
}

func (m *Container) InitRetailerService() {
	m.RetailerService = retailer.NewRetailerService(
		m.ApplicationServices.UserService,
		retailer_db_port.NewPostgres(m.db, m.ApplicationServices.Pagination))
}

func (m *Container) InitOrderService() {
	m.OrderService = order.NewOrderService(
		order_db_port.NewPostgres(
			m.db, m.ApplicationServices.Pagination),
		m.DistributorService,
		m.InvoiceService,
		m.ProductService,
		m.RetailerService,
		m.ApplicationServices.CheckoutService,
		m.ConfigService)
}

func (m *Container) InitPaymentVerificationService() {
	m.PaymentVerificationService = payment_verification.NewPaymentVerificationService(
		m.ApplicationServices.PaymentService,
		m.ApplicationServices.PaymentPartnerService,
		m.ApplicationServices.TransactionService,
		m.OrderService)
}

func (m *Container) InitConfigService() {
	m.ConfigService = config_module.NewConfig(config_db_adapter.NewPostgres(m.db))
}
