package core

import (
	"database/sql"

	category_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	"b2b.nati011.github.com/internal/core/domain/category"
)

type Container struct {
	CategoryService category.Provider
}

func NewContainer(db *sql.DB) *Container {
	container := Container{}

	container.InitCategoryService(db)

	return &container
}

func (m *Container) InitCategoryService(db *sql.DB) {
	m.CategoryService = category.NewCategory(category_db_port.NewPostgres(db))
}
