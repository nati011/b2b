package app

import (
	"context"
	"database/sql"
	"marketplace/internal/config"
	customerhttp "marketplace/internal/core/customer/api/http"
	customerservice "marketplace/internal/core/customer/service"
	orderhttp "marketplace/internal/core/order/api/http"
	orderservice "marketplace/internal/core/order/service"
	producthttp "marketplace/internal/core/product/api/http"
	productservice "marketplace/internal/core/product/service"
	supplierhttp "marketplace/internal/core/supplier/api/http"
	supplierservice "marketplace/internal/core/supplier/service"
	basicauthservice "marketplace/internal/infra/auth/basic/service"
	permissionhttp "marketplace/internal/infra/authz/permission/api/http"
	permissionservice "marketplace/internal/infra/authz/permission/service"
	resourceservice "marketplace/internal/infra/authz/resource/service"
	role "marketplace/internal/infra/authz/role/api/http"
	roleservice "marketplace/internal/infra/authz/role/service"
	http "marketplace/internal/infra/http"
	userhttp "marketplace/internal/infra/user/api/http"
	userservice "marketplace/internal/infra/user/service"
)

type Application struct {
	Config           *config.Config
	DB               *sql.DB
	Server           http.Server
	UserService      *userservice.Service
	UserHandler      *userhttp.UserHandler
	RoleService      *roleservice.RoleService
	RoleHandler      *role.RoleHandler
	PermService      *permissionservice.PermissionService
	PermHandler      *permissionhttp.PermissionHandler
	ResourceService  *resourceservice.Service
	BasicAuthService *basicauthservice.Service
	CustomerService  *customerservice.CustomerService
	CustomerHandler  *customerhttp.CustomerHandler
	SupplierService  *supplierservice.SupplierService
	SupplierHandler  *supplierhttp.SupplierHandler
	ProductService   *productservice.Service
	ProductHandler   *producthttp.ProductHandler
	OrderService     *orderservice.Service
	OrderHandler     *orderhttp.OrderHandler
}

// Bootstrap seeds resources, creates default roles, and creates admin users.
// This method is idempotent and can be called multiple times safely.
// It should be called before Start() during application initialization.
func (a *Application) Bootstrap(ctx context.Context) error {
	return Bootstrap(ctx, a.Config, a.UserService, a.RoleService, a.PermService, a.ResourceService, a.BasicAuthService)
}

// Start starts the HTTP server.
// Bootstrap should be called before Start() during application initialization.
func (a *Application) Start() error {
	return a.Server.Start()
}

func (a *Application) Shutdown(ctx context.Context) error {
	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}
	a.DB.Close()
	return nil
}
