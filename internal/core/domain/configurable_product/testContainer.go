package configurable_product

import (
	configurableProduct_db "b2b.nati011.github.com/internal/adapter/secondary/domain/configurable_product/db"
	product_db "b2b.nati011.github.com/internal/adapter/secondary/domain/product/db"
	"b2b.nati011.github.com/internal/core/application/event"
	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/product"
)

type TestContainer struct {
	CategoryService            category.Provider
	ProductService             product.Provider
	ConfigurableProductService Provider
	DistributorService         distributor.Provider
	Event                      *event.Broker
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.DistributorService = distributor.NewPackageIntegrationTestContainer().DistributorService
	container.ProductService = product.NewProduct(
		product_db.NewMock(),
		container.CategoryService,
		container.DistributorService,
		container.Event)

	container.ConfigurableProductService = NewConfigurableProductService(
		configurableProduct_db.NewMock(),
		container.ProductService,
	)
	return container
}

func (t *TestContainer) Teardown() {
	t.ConfigurableProductService = NewConfigurableProductService(
		configurableProduct_db.NewMock(),
		t.ProductService,
	)
}
