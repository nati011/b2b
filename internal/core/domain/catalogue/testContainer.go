package catalogue

import (
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/product"
)

type TestContainer struct {
	ProductService             product.Provider
	ConfigurableProductService configurable_product.Provider
}

func NewPackageTestContainer() TestContainer {
	container := TestContainer{}
	container.ProductService = product.NewPackageIntegrationTestContainer().ProductService
	container.ConfigurableProductService = configurable_product.NewPackageIntegrationTestContainer().ConfigurableProductService
	return container
}

func (t *TestContainer) Teardown() {
	t.ProductService = product.NewPackageIntegrationTestContainer().ProductService
	t.ConfigurableProductService = configurable_product.NewPackageIntegrationTestContainer().ConfigurableProductService
}
