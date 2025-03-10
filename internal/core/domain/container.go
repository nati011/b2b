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
// ├── CategoryService
// │   └── ProductService
// │       └── ConfigurableProductService
// ├── InvoiceService
// │   └── OrderService

type Container struct {
	CategoryService            category.Provider
	ProductService             product.Provider
	ConfigurableProductService configurable_product.Provider
	InvoiceService             invoice.Provider
	OrderService               order.Provider
}

func NewContainer(db *sql.DB) *Container {
	container := Container{}

	container.InitCategoryService(db)

	return &container
}

func (m *Container) InitCategoryService(db *sql.DB) {
	m.CategoryService = category.NewCategory(category_db_port.NewPostgres(db))
}

func (m *Container) InitProductService(db *sql.DB) {
	m.ProductService = product.NewProduct(product_db_port.NewPostgres(db), m.CategoryService)
}

func (m *Container) InitConfigrableProductService(db *sql.DB) {
	m.ConfigurableProductService = configurable_product.NewConfigurableProductService(configurable_product_db_port.NewPostgres(db),
		m.ProductService,
	)
}

func (m *Container) InitInvoiceService(db *sql.DB) {
	m.InvoiceService = invoice.NewInvoice(invoice_db_port.NewPostgres(db))
}

func (m *Container) InitOrderService(db *sql.DB) {
	m.OrderService = order.NewOrderService(order_db_port.NewPostgres(db), m.InvoiceService)
}
