package core

import (
	"database/sql"

	category_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	configurable_product_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/configurable_product/db"
	invoice_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/invoice/db"
	order_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/order/db"
	product_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/product/db"
	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
)

/*Dependency Tree*/

// Container
// ├── DB
// │   ├── CategoryService
// │   │   └── ProductService
// │   │       └── ConfigurableProductService
// │   └── InvoiceService
// │       └── OrderService

type Container struct {
	db                         *sql.DB
	CategoryService            category.Provider
	ProductService             product.Provider
	ConfigurableProductService configurable_product.Provider
	InvoiceService             invoice.Provider
	OrderService               order.Provider
}

func NewContainer(db *sql.DB) *Container {
	container := Container{}
	container.db = db

	container.InitCategoryService()
	container.InitProductService()
	container.InitConfigrableProductService()
	container.InitInvoiceService()
	container.InitOrderService()

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
	m.OrderService = order.NewOrderService(order_db_port.NewPostgres(m.db), m.InvoiceService)
}
