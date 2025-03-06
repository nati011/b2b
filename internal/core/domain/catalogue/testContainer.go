package catalogue

import (
	"b2b.nati011.github.com/internal/core/domain/catalogue/category"
	"b2b.nati011.github.com/internal/core/domain/catalogue/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/catalogue/product"

	category_db "b2b.nati011.github.com/internal/adapter/secondary/domain/catalogue/category"
	configurableProduct_db "b2b.nati011.github.com/internal/adapter/secondary/domain/catalogue/configurable_product"
	product_db "b2b.nati011.github.com/internal/adapter/secondary/domain/catalogue/product"
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
