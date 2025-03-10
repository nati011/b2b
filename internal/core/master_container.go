package core

import (
	"database/sql"

	category_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/category"
	"b2b.nati011.github.com/internal/core/domain/category"
)

type MasterContainer struct {
	CategoryService category.Provider
}

func NewMasterContainer(db *sql.DB) *MasterContainer {
	return &MasterContainer{}
}

func (m *MasterContainer) InitCategoryService(db *sql.DB) {
	m.CategoryService = category.NewCategory(category_db_port.NewPostgres(db))
}
