package test_container

import (
	"database/sql"

	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/product"
	product_test_container "b2b.nati011.github.com/internal/core/domain/product/test"

	configurableProduct_db "b2b.nati011.github.com/internal/adapter/secondary/domain/configurable_product/db"
)

type TestContainer struct {
	CategoryService            category.Provider
	ProductService             product.Provider
	ConfigurableProductService configurable_product.Provider
}

func NewDBIntegrationTestContainer(db *sql.DB) TestContainer {
	container := TestContainer{}
	container.ProductService = product_test_container.NewDBIntegrationTestContainer(db).ProductService
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
