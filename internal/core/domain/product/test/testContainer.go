package catalogue

import (
	"database/sql"

	"b2b.nati011.github.com/config"
	"b2b.nati011.github.com/internal/core/application/event"
	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	distributor_test "b2b.nati011.github.com/internal/core/domain/distributor/test"
	"b2b.nati011.github.com/internal/core/domain/product"

	category_db "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	configurableProduct_db "b2b.nati011.github.com/internal/adapter/secondary/domain/configurable_product/db"
	product_db "b2b.nati011.github.com/internal/adapter/secondary/domain/product/db"
)

type TestContainer struct {
	CategoryService            category.Provider
	ProductService             product.Provider
	ConfigurableProductService configurable_product.Provider
	DistributorService         distributor.Provider
	Event                      *event.Broker
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.CategoryService = category.NewCategory(
		category_db.NewMock(),
	)
	container.DistributorService = distributor.NewPackageIntegrationTestContainer().DistributorService
	container.ProductService = product.NewProduct(
		product_db.NewMock(),
		container.CategoryService,
		container.DistributorService,
		container.Event)

	container.ConfigurableProductService =
		configurable_product.NewConfigurableProductService(
			configurableProduct_db.NewMock(),
			container.ProductService,
		)
	return container
}

func NewDBIntegrationTestContainer(db *sql.DB) TestContainer {
	container := TestContainer{}
	container.CategoryService = category.NewCategory(
		category_db.NewMock(),
	)
	container.DistributorService = distributor_test.NewDBIntegrationTestContainer(db).DistributorService
	container.ProductService = product.NewProduct(
		product_db.NewPostgres(db, config.DefaultPaginationBuilder().Build()),
		container.CategoryService,
		container.DistributorService,
		container.Event)

	container.ConfigurableProductService =
		configurable_product.NewConfigurableProductService(
			configurableProduct_db.NewMock(),
			container.ProductService,
		)
	return container
}

func (t *TestContainer) Teardown(db *sql.DB) {
	t.ProductService = product.NewProduct(
		product_db.NewPostgres(db, config.DefaultPaginationBuilder().Build()),
		t.CategoryService,
		t.DistributorService,
		t.Event)
}
