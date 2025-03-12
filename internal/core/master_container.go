package core

import (
	"database/sql"

	distributor_db_port "b2b.nati011.github.com/internal/adapter/secondary/application/distributor/db"
	category_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/category"
	auth_port "b2b.nati011.github.com/internal/port/application/auth/provider"

	"b2b.nati011.github.com/internal/core/application/auth"
	"b2b.nati011.github.com/internal/core/domain/category"
	distributor "b2b.nati011.github.com/internal/core/domain/distributor/service"
)

type MasterContainer struct {
	AuthService        auth_port.Provider
	CategoryService    category.Provider
	DistributorService distributor.Provider
}

func NewMasterContainer(db *sql.DB) *MasterContainer {
	container := MasterContainer{}
	container.InitCategoryService(db)

	return &container
}

func (m *MasterContainer) InitAuthService(ap auth_port.Provider) {
	m.AuthService = auth.NewAuthService(
		ap,
	)
}

func (m *MasterContainer) InitCategoryService(db *sql.DB) {
	m.CategoryService = category.NewCategory(category_db_port.NewPostgres(db))
}

func (m *MasterContainer) InitDistributorService(db *sql.DB) {
	m.DistributorService = distributor.NewDistributorService(
		distributor_db_port.NewPostgres(db),
		m.AuthService,
	)
}
