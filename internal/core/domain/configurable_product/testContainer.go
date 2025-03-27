package configurable_product

import (
	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/product"

	configurableProduct_db "b2b.nati011.github.com/internal/adapter/secondary/domain/configurable_product/db"
)

type TestContainer struct {
	CategoryService            category.Provider
	ProductService             product.Provider
	ConfigurableProductService Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.ProductService = product.NewPackageIntegrationTestContainer().ProductService
	container.ConfigurableProductService = NewConfigurableProductService(
		configurableProduct_db.NewMock(),
		container.ProductService,
	)
	return container
}

func (t *TestContainer) cleanup() {
	t.ProductService = product.NewPackageIntegrationTestContainer().ProductService
	t.ConfigurableProductService = NewConfigurableProductService(
		configurableProduct_db.NewMock(),
		t.ProductService,
	)
}
