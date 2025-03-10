package catalogue

import (
	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/product"

	category_db "b2b.nati011.github.com/internal/adapter/secondary/domain/category"
	configurableProduct_db "b2b.nati011.github.com/internal/adapter/secondary/domain/configurable_product"
	product_db "b2b.nati011.github.com/internal/adapter/secondary/domain/product"
)

type TestContainer struct {
	CategoryService            category.Provider
	ProductService             product.Provider
	ConfigurableProductService configurable_product.Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.CategoryService = category.NewCategory(
		category_db.NewMock(),
	)
	container.ProductService = product.NewProduct(
		product_db.NewMock(),
		container.CategoryService,
	)
	container.ConfigurableProductService =
		configurable_product.NewConfigurableProductService(
			configurableProduct_db.NewMock(),
			container.ProductService,
		)
	return container
}
