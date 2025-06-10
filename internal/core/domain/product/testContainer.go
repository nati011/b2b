package product

import (
	category_db "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	db "b2b.nati011.github.com/internal/adapter/secondary/domain/product/db"
	category "b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/distributor"
)

type TestContainer struct {
	CategoryService    category.Provider
	ProductService     Provider
	DistributorService distributor.Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.CategoryService = category.NewCategory(
		category_db.NewMock(),
	)
	container.DistributorService = distributor.NewPackageIntegrationTestContainer().DistributorService
	container.ProductService = NewProduct(
		db.NewMock(),
		container.CategoryService,
		container.DistributorService,
	)
	return container
}

func (t *TestContainer) Teardown() {
	t.ProductService = NewProduct(
		db.NewMock(),
		t.CategoryService,
		t.DistributorService,
	)
}
