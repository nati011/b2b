package core

import (
	"database/sql"

	category_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	configurable_product_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/configurable_product/db"
	distributor_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/distributor/db"
	invoice_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/invoice/db"
	order_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/order/db"
	product_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/product/db"
	retailer_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/retailer/db"
	application_core "b2b.nati011.github.com/internal/core/application"
	payment "b2b.nati011.github.com/internal/core/application/payment"
	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/order"
	payment_processor "b2b.nati011.github.com/internal/core/domain/paymentProcessor"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

/*Dependency Tree*/

// Container
// ├── CategoryService
// │   └── ProductService
// │     └── ConfigurableProductService
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
	PaymentProcessorService    payment_processor.Provider
	PaymentService             payment.Provider
}

func NewContainer(application_core application_core.Container, db *sql.DB) *Container {
	container := Container{}
	container.db = db
	container.ApplicationServices = application_core

	// ORDER ORDER!!

	//	messing up the order creates chaos

	container.InitCategoryService()
	container.InitProductService()
	container.InitConfigrableProductService()
	container.InitInvoiceService()
	container.InitRetailerService()
	container.InitOrderService()
	container.InitDistributorService()
	container.InitPaymentProcessorService()
	container.InitPaymentService()

	return &container
}

func (m *Container) InitCategoryService() {
	m.CategoryService = category.NewCategory(category_db_port.NewPostgres(m.db))
}

func (m *Container) InitProductService() {
	m.ProductService = product.NewProduct(product_db_port.NewPostgres(m.db), m.CategoryService)
}

func (m *Container) InitConfigrableProductService() {
	m.ConfigurableProductService = configurable_product.NewConfigurableProductService(configurable_product_db_port.NewPostgres(m.db),
		m.ProductService,
	)
}

func (m *Container) InitInvoiceService() {
	m.InvoiceService = invoice.NewInvoice(invoice_db_port.NewPostgres(m.db))
}

func (m *Container) InitOrderService() {
	m.OrderService = order.NewOrderService(order_db_port.NewPostgres(m.db), m.InvoiceService, m.ProductService, m.RetailerService)
}

func (m *Container) InitDistributorService() {
	m.DistributorService = distributor.NewDistributorService(m.ApplicationServices.UserService, distributor_db_port.NewPostgres(m.db))
}

func (m *Container) InitRetailerService() {
	m.RetailerService = retailer.NewRetailerService(m.ApplicationServices.UserService, retailer_db_port.NewPostgres(m.db))
}

func (m *Container) InitPaymentProcessorService() {
	m.PaymentProcessorService = payment_processor.NewProcessor(m.OrderService)
}

func (m *Container) InitPaymentService() {
	m.PaymentService = payment.NewPaymentService(
		m.ApplicationServices.PaymentPartnerService,
		m.ApplicationServices.TransactionService,
		m.PaymentProcessorService,
	)
}
