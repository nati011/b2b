package catalogue

import (
	category_db "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	configurableProduct_db "b2b.nati011.github.com/internal/adapter/secondary/domain/configurable_product/db"
	product_db "b2b.nati011.github.com/internal/adapter/secondary/domain/product/db"
	"b2b.nati011.github.com/internal/core/application/event"
	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/product"
)

type TestContainer struct {
	CategoryService            category.Provider
	ProductService             product.Provider
	DistributorService         distributor.Provider
	ConfigurableProductService configurable_product.Provider
	Event                      *event.Broker
	CatalogueService           Provider
}

func NewPackageTestContainer() TestContainer {
	container := TestContainer{}
	container.Event = event.NewBroker()
	container.CategoryService = category.NewCategory(
		category_db.NewMock(),
	)
	container.DistributorService = distributor.NewPackageIntegrationTestContainer().DistributorService
	container.ProductService = product.NewProduct(
		product_db.NewMock(),
		container.CategoryService,
		container.DistributorService,
		container.Event)
	container.ConfigurableProductService = configurable_product.NewConfigurableProductService(
		configurableProduct_db.NewMock(),
		container.ProductService,
	)
	container.CatalogueService = NewCatalogueService(
		container.ProductService,
		container.ConfigurableProductService)

	return container
}

func (t *TestContainer) Teardown() {
	t.ProductService = product.NewPackageIntegrationTestContainer().ProductService
	t.ConfigurableProductService = configurable_product.NewPackageIntegrationTestContainer().ConfigurableProductService
}
