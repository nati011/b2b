package product

import (
	category_db "b2b.nati011.github.com/internal/adapter/secondary/domain/catalogue/category"
	db "b2b.nati011.github.com/internal/adapter/secondary/domain/catalogue/product"
	category "b2b.nati011.github.com/internal/core/domain/category"
)

type TestContainer struct {
	CategoryService category.Provider
	ProductService  Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.CategoryService = category.NewCategory(
		category_db.NewMock(),
	)
	container.ProductService = NewProduct(
		db.NewMock(),
		container.CategoryService,
	)
	return container
}
