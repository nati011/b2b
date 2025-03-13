package core

import (
	"database/sql"

	config "b2b.nati011.github.com/config"
	"b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
	distributor_db_port "b2b.nati011.github.com/internal/adapter/secondary/application/distributor/db"
	category_db_port "b2b.nati011.github.com/internal/adapter/secondary/domain/category"
	"b2b.nati011.github.com/internal/core/application/auth"
	"b2b.nati011.github.com/internal/core/domain/category"

	distributor "b2b.nati011.github.com/internal/core/domain/distributor/service"
	auth_port "b2b.nati011.github.com/internal/port/application/auth/provider"
)

type MasterContainer struct {
	AuthService        auth_port.Provider
	CategoryService    category.Provider
	DistributorService distributor.Provider
}

func NewMasterContainer(db *sql.DB) *MasterContainer {
	container := MasterContainer{}

	container.InitCategoryService(db)
	container.InitAuthService()

	return &container
}

func (m *MasterContainer) InitAuthService() {
	appConfig := config.AppConfig
	keyCloakProvider := provider.NewKeycloakProvider(
		appConfig.KeycloakInstanceURL,
		appConfig.KeycloakUsername,
		appConfig.KeycloakPassword,
		appConfig.KeycloakRealm,
		appConfig.KeycloakApplicationRealm,
		appConfig.KeycloakClientId,
	)
	m.AuthService = auth.NewAuthService(
		keyCloakProvider,
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

func (m *MasterContainer) InitResourceService(db *sql.DB) {
	m.CategoryService = category.NewCategory(category_db_port.NewPostgres(db))
}
