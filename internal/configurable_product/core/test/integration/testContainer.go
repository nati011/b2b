package test_container

import (
	"database/sql"

	"b2b.nati011.github.com/config"
	configurableProduct_db "b2b.nati011.github.com/internal/adapter/secondary/domain/configurable_product/db"
	product_db "b2b.nati011.github.com/internal/adapter/secondary/domain/product/db"
	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	distributor_test "b2b.nati011.github.com/internal/core/domain/distributor/test"
	"b2b.nati011.github.com/internal/core/domain/product"
)

type TestContainer struct {
	CategoryService            category.Provider
	ProductService             product.Provider
	ConfigurableProductService configurable_product.Provider
	DistributorService         distributor.Provider
}

func NewDBIntegrationTestContainer(db *sql.DB) TestContainer {
	container := TestContainer{}
	container.DistributorService = distributor_test.NewDBIntegrationTestContainer(db).DistributorService
	container.ProductService = product.NewProduct(
		product_db.NewPostgres(db, config.DefaultPaginationBuilder().Build()),
		container.CategoryService,
		container.DistributorService,
	)
	container.ConfigurableProductService = configurable_product.NewConfigurableProductService(
		configurableProduct_db.NewPostgres(db),
		container.ProductService,
	)
	return container
}

func (t *TestContainer) Teardown() {
	t.ConfigurableProductService = configurable_product.NewConfigurableProductService(
		configurableProduct_db.NewMock(),
		t.ProductService,
	)
}
