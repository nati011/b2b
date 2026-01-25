//go:build wireinject
// +build wireinject

package app

import (
	"database/sql"
	stdhttp "net/http"

	"github.com/google/wire"

	"marketplace/internal/config"
	customerhttp "marketplace/internal/core/customer/api/http"
	customerrepo "marketplace/internal/core/customer/repository"
	customerservice "marketplace/internal/core/customer/service"
	orderhttp "marketplace/internal/core/order/api/http"
	supplierhttp "marketplace/internal/core/supplier/api/http"
	supplierrepo "marketplace/internal/core/supplier/repository"
	supplierservice "marketplace/internal/core/supplier/service"
	orderrepo "marketplace/internal/core/order/repository"
	orderservice "marketplace/internal/core/order/service"
	producthttp "marketplace/internal/core/product/api/http"
	productrepo "marketplace/internal/core/product/repository"
	productservice "marketplace/internal/core/product/service"
	referralhttp "marketplace/internal/core/referral/api/http"
	referralrepo "marketplace/internal/core/referral/repository"
	referralservice "marketplace/internal/core/referral/service"
	authDomain "marketplace/internal/infra/auth"
	basicauth "marketplace/internal/infra/auth/basic"
	basicauthhttp "marketplace/internal/infra/auth/basic/api/http"
	basicauthrepo "marketplace/internal/infra/auth/basic/repository"
	basicauthservice "marketplace/internal/infra/auth/basic/service"
	permissionhttp "marketplace/internal/infra/authz/permission/api/http"
	permissionrepo "marketplace/internal/infra/authz/permission/repository"
	permissionservice "marketplace/internal/infra/authz/permission/service"
	resourcerepo "marketplace/internal/infra/authz/resource/repository"
	resourceservice "marketplace/internal/infra/authz/resource/service"
	rolehttp "marketplace/internal/infra/authz/role/api/http"
	rolerepo "marketplace/internal/infra/authz/role/repository"
	roleservice "marketplace/internal/infra/authz/role/service"
	"marketplace/internal/infra/db"
	http "marketplace/internal/infra/http"
	"marketplace/internal/infra/idempotency"
	userhttp "marketplace/internal/infra/user/api/http"
	userrepository "marketplace/internal/infra/user/repository"
	userservice "marketplace/internal/infra/user/service"
	"marketplace/pkg/http/middleware"
)

// InitializeApp wires up the application dependencies.
func InitializeApp(cfg *config.Config) (*Application, error) {
	wire.Build(
		provideDB,
		provideTxManager,
		provideHTTPServer,
		provideResourceRepository,
		provideResourceService,
		providePermissionRepository,
		providePermissionService,
		providePermissionHandler,
		provideRoleRepository,
		provideRoleService,
		provideRoleHandler,
		provideUserRepository,
		providePhoneValidationRepository,
		providePhoneValidationService,
		provideRegistrationTokenRepository,
		provideRegistrationTokenService,
		provideUserService,
		provideRegistrationTokenAdapter,
		provideUserPermissionChecker,
		provideUserHandler,
		provideCustomerRepository,
		provideCustomerService,
		provideCustomerHandler,
		provideSupplierRepository,
		provideSupplierService,
		provideSupplierHandler,
		provideProductRepository,
		provideProductService,
		provideProductHandler,
		provideOrderRepository,
		provideOrderService,
		provideOrderHandler,
		provideReferralRepository,
		provideReferralService,
		provideReferralHandler,
		provideBasicAuthRepository,
		provideBasicAuthService,
		provideBasicAuthAdapter,
		provideBasicAuthHandler,
		provideIdempotencyStore,
		provideIdempotencyMiddleware,
		provideAuthMiddleware,
		provideAuthorizationMiddleware,
		provideRegistrationTokenMiddleware,
		provideHTTPMiddleware,
		provideHTTPHandler,
		provideApplicationInfrastructure,
		provideApplicationServices,
		provideApplicationHandlers,
		provideApplication,
	)
	return nil, nil
}

// provideDB creates a DB connection.
func provideDB(cfg *config.Config) (*sql.DB, error) {
	return db.NewPostgres(cfg.DB.DSN())
}

// provideTxManager creates a transaction manager for the provided DB.
func provideTxManager(dbConn *sql.DB) db.TxManager {
	return db.NewTxManager(dbConn)
}

// provideHTTPServer creates an HTTP server.
func provideHTTPServer(cfg *config.Config, handler stdhttp.Handler) http.Server {
	srv := http.New(http.Options{
		Address:      cfg.HttpServer.Address(),
		Handler:      handler,
		ReadTimeout:  cfg.HttpServer.ReadTimeout,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	})
	return *srv
}

// provideUserRepository creates a user repository.
func provideUserRepository(dbConn *sql.DB) *userrepository.Repository {
	return userrepository.NewRepository(dbConn)
}

// providePhoneValidationRepository creates a phone validation repository.
func providePhoneValidationRepository(dbConn *sql.DB) userservice.PhoneValidationRepository {
	return userrepository.NewPhoneValidationRepository(dbConn)
}

// providePhoneValidationService creates a phone validation service.
func providePhoneValidationService(repo userservice.PhoneValidationRepository) *userservice.PhoneValidationService {
	return userservice.NewPhoneValidationService(repo)
}

// provideRegistrationTokenRepository creates a registration token repository.
func provideRegistrationTokenRepository(dbConn *sql.DB) userservice.RegistrationTokenRepository {
	return userrepository.NewRegistrationTokenRepository(dbConn)
}

// provideRegistrationTokenService creates a registration token service.
func provideRegistrationTokenService(repo userservice.RegistrationTokenRepository) *userservice.RegistrationTokenService {
	return userservice.NewRegistrationTokenService(repo)
}

// provideRegistrationTokenAdapter creates a registration token adapter.
func provideRegistrationTokenAdapter(service *userservice.RegistrationTokenService) *basicauth.RegistrationTokenAdapter {
	return basicauth.NewRegistrationTokenAdapter(service)
}

// provideUserService creates a user service.
func provideUserService(repo *userrepository.Repository, roleRepo *rolerepo.RoleRepository, phoneValidation *userservice.PhoneValidationService, registrationTokenService *userservice.RegistrationTokenService) *userservice.Service {
	return userservice.NewService(repo, roleRepo, phoneValidation, registrationTokenService)
}

// provideUserHandler creates a user HTTP handler.
func provideUserHandler(userService *userservice.Service, permissionChecker *userservice.UserPermissionChecker) *userhttp.UserHandler {
	return userhttp.NewUserHandler(userService, permissionChecker)
}

// provideCustomerRepository creates a customer repository.
func provideCustomerRepository(dbConn *sql.DB) *customerrepo.CustomerRepository {
	return customerrepo.NewCustomerRepository(dbConn)
}

// provideCustomerService creates a customer service.
func provideCustomerService(repo *customerrepo.CustomerRepository) *customerservice.CustomerService {
	return customerservice.NewCustomerService(repo)
}

// provideCustomerHandler creates a customer HTTP handler.
func provideCustomerHandler(service *customerservice.CustomerService) *customerhttp.CustomerHandler {
	return customerhttp.NewCustomerHandler(service)
}

// provideSupplierRepository creates a supplier repository.
func provideSupplierRepository(dbConn *sql.DB) *supplierrepo.SupplierRepository {
	return supplierrepo.NewSupplierRepository(dbConn)
}

// provideSupplierService creates a supplier service.
func provideSupplierService(repo *supplierrepo.SupplierRepository) *supplierservice.SupplierService {
	return supplierservice.NewSupplierService(repo)
}

// provideSupplierHandler creates a supplier HTTP handler.
func provideSupplierHandler(service *supplierservice.SupplierService) *supplierhttp.SupplierHandler {
	return supplierhttp.NewSupplierHandler(service)
}

// provideProductRepository creates a product repository.
func provideProductRepository(dbConn *sql.DB) *productrepo.Repository {
	return productrepo.NewRepository(dbConn)
}

// provideProductService creates a product service.
func provideProductService(repo *productrepo.Repository) *productservice.Service {
	return productservice.NewService(repo)
}

// provideProductHandler creates a product HTTP handler.
func provideProductHandler(service *productservice.Service) *producthttp.ProductHandler {
	return producthttp.NewProductHandler(service)
}

// provideOrderRepository creates an order repository.
func provideOrderRepository(dbConn *sql.DB) *orderrepo.Repository {
	return orderrepo.NewRepository(dbConn)
}

// provideOrderService creates an order service.
func provideOrderService(repo *orderrepo.Repository) *orderservice.Service {
	return orderservice.NewService(repo)
}

// provideOrderHandler creates an order HTTP handler.
func provideOrderHandler(service *orderservice.Service) *orderhttp.OrderHandler {
	return orderhttp.NewOrderHandler(service)
}

// provideReferralRepository creates a referral repository.
func provideReferralRepository(dbConn *sql.DB) *referralrepo.ReferralRepository {
	return referralrepo.NewReferralRepository(dbConn)
}

// provideReferralService creates a referral service.
func provideReferralService(repo *referralrepo.ReferralRepository) *referralservice.ReferralService {
	return referralservice.NewReferralService(repo)
}

// provideReferralHandler creates a referral HTTP handler.
func provideReferralHandler(service *referralservice.ReferralService) *referralhttp.ReferralHandler {
	return referralhttp.NewReferralHandler(service)
}

// provideRoleRepository creates a role repository.
func provideRoleRepository(dbConn *sql.DB, txManager db.TxManager) *rolerepo.RoleRepository {
	return rolerepo.NewRoleRepository(dbConn, txManager)
}

// provideRoleService creates a role service.
func provideRoleService(repo *rolerepo.RoleRepository, permRepo *permissionrepo.PermissionRepository) *roleservice.RoleService {
	return roleservice.NewRoleService(repo, permRepo)
}

// provideRoleHandler creates a role HTTP handler.
func provideRoleHandler(roleService *roleservice.RoleService) *rolehttp.RoleHandler {
	return rolehttp.NewRoleHandler(roleService)
}

// providePermissionRepository creates a permission repository.
func providePermissionRepository(dbConn *sql.DB) *permissionrepo.PermissionRepository {
	return permissionrepo.NewPermissionRepository(dbConn)
}

// providePermissionService creates a permission service.
func providePermissionService(repo *permissionrepo.PermissionRepository, resourceService *resourceservice.Service) *permissionservice.PermissionService {
	return permissionservice.NewPermissionService(repo, resourceService)
}

// providePermissionHandler creates a permission HTTP handler.
func providePermissionHandler(service *permissionservice.PermissionService) *permissionhttp.PermissionHandler {
	return permissionhttp.NewPermissionHandler(service)
}

// provideResourceRepository creates a resource repository.
func provideResourceRepository(dbConn *sql.DB, txManager db.TxManager) *resourcerepo.Repository {
	return resourcerepo.NewRepository(dbConn, txManager)
}

// provideResourceService creates a resource service.
func provideResourceService(repo *resourcerepo.Repository) *resourceservice.Service {
	return resourceservice.NewService(repo)
}

// provideIdempotencyStore wires the idempotency store.
func provideIdempotencyStore(dbConn *sql.DB) *idempotency.Store {
	return idempotency.NewStore(dbConn)
}

// provideIdempotencyMiddleware wires the idempotency middleware.
func provideIdempotencyMiddleware(store *idempotency.Store) *middleware.IdempotencyMiddleware {
	return middleware.NewIdempotencyMiddleware(store)
}

// provideBasicAuthRepository creates a basic auth credential repository.
func provideBasicAuthRepository(dbConn *sql.DB) *basicauthrepo.CredentialRepository {
	return basicauthrepo.NewCredentialRepository(dbConn)
}

// provideUserPermissionChecker creates a permission checker for auth service.
func provideUserPermissionChecker(userRepo *userrepository.Repository, roleService *roleservice.RoleService) *userservice.UserPermissionChecker {
	return userservice.NewUserPermissionChecker(userRepo, roleService)
}

// provideBasicAuthService creates a basic auth service.
func provideBasicAuthService(repo *basicauthrepo.CredentialRepository, userService *userservice.Service, permissionChecker *userservice.UserPermissionChecker) *basicauthservice.Service {
	return basicauthservice.NewService(repo, userService, permissionChecker)
}

// provideBasicAuthAdapter creates an adapter for the basic auth service.
func provideBasicAuthAdapter(service *basicauthservice.Service) *basicauth.Adapter {
	return basicauth.NewAdapter(service)
}

// provideBasicAuthHandler creates a basic auth credential handler.
func provideBasicAuthHandler(service *basicauthservice.Service) *basicauthhttp.Handler {
	return basicauthhttp.NewHandler(service)
}

// provideRegistrationTokenMiddleware creates middleware for validating registration tokens.
func provideRegistrationTokenMiddleware(adapter *basicauth.RegistrationTokenAdapter) *middleware.RegistrationTokenMiddleware {
	return middleware.NewRegistrationTokenMiddleware(adapter, []string{"/auth/basic/credentials"})
}

// collectPublicRoutes collects public routes from all registered modules.
func collectPublicRoutes() []string {
	return middleware.CollectPublicRoutes()
}

// collectAuthenticatedRoutes collects authenticated routes from all registered modules.
func collectAuthenticatedRoutes() []string {
	return middleware.CollectAuthenticatedRoutes()
}

// provideAuthMiddleware wires the authentication middleware.
func provideAuthMiddleware(cfg *config.Config, basicAuthAdapter *basicauth.Adapter, userService *userservice.Service) *middleware.AuthMiddleware {
	realm := "Restricted"
	if cfg.Auth.Basic != nil && cfg.Auth.Basic.Realm != "" {
		realm = cfg.Auth.Basic.Realm
	}

	var basicAuthService middleware.BasicAuthService
	if cfg.Auth.Mode == string(authDomain.AuthModeBasic) {
		basicAuthService = basicAuthAdapter
	}

	publicRoutes := collectPublicRoutes()
	return middleware.NewAuthMiddleware(authDomain.AuthMode(cfg.Auth.Mode), realm, basicAuthService, userService, publicRoutes)
}

// provideAuthorizationMiddleware wires the authorization middleware.
func provideAuthorizationMiddleware(resourceService *resourceservice.Service, permissionChecker *userservice.UserPermissionChecker) *middleware.AuthorizationMiddleware {
	publicRoutes := collectPublicRoutes()
	allowedAuthenticatedRoutes := collectAuthenticatedRoutes()
	return middleware.NewAuthorizationMiddleware(resourceService, permissionChecker, publicRoutes, allowedAuthenticatedRoutes)
}

type httpMiddleware struct {
	idempotencyMiddleware       *middleware.IdempotencyMiddleware
	authMiddleware              *middleware.AuthMiddleware
	authorizationMiddleware     *middleware.AuthorizationMiddleware
	registrationTokenMiddleware *middleware.RegistrationTokenMiddleware
}

// provideHTTPMiddleware creates the httpMiddleware struct.
func provideHTTPMiddleware(idempotencyMiddleware *middleware.IdempotencyMiddleware, authMiddleware *middleware.AuthMiddleware, authorizationMiddleware *middleware.AuthorizationMiddleware, registrationTokenMiddleware *middleware.RegistrationTokenMiddleware) httpMiddleware {
	return httpMiddleware{
		idempotencyMiddleware:       idempotencyMiddleware,
		authMiddleware:              authMiddleware,
		authorizationMiddleware:     authorizationMiddleware,
		registrationTokenMiddleware: registrationTokenMiddleware,
	}
}

// provideHTTPHandler registers HTTP routes and returns a handler.
func provideHTTPHandler(userHandler *userhttp.UserHandler, roleHandler *rolehttp.RoleHandler, permHandler *permissionhttp.PermissionHandler, basicAuthHandler *basicauthhttp.Handler, customerHandler *customerhttp.CustomerHandler, supplierHandler *supplierhttp.SupplierHandler, productHandler *producthttp.ProductHandler, orderHandler *orderhttp.OrderHandler, referralHandler *referralhttp.ReferralHandler, mw httpMiddleware) stdhttp.Handler {
	mux := stdhttp.NewServeMux()
	userhttp.RegisterHTTPRoutes(mux, userHandler)
	rolehttp.RegisterHTTPRoutes(mux, roleHandler)
	permissionhttp.RegisterHTTPRoutes(mux, permHandler)
	basicauthhttp.RegisterHTTPRoutes(mux, basicAuthHandler)
	customerhttp.RegisterHTTPRoutes(mux, customerHandler)
	supplierhttp.RegisterHTTPRoutes(mux, supplierHandler)
	producthttp.RegisterHTTPRoutes(mux, productHandler)
	orderhttp.RegisterHTTPRoutes(mux, orderHandler)
	referralhttp.RegisterHTTPRoutes(mux, referralHandler)

	handler := mw.idempotencyMiddleware.Handle(mux)
	handler = mw.authorizationMiddleware.Handle(handler)
	handler = mw.authMiddleware.Handle(handler)
	handler = mw.registrationTokenMiddleware.Handle(handler)
	loggingMiddleware := middleware.NewLoggingMiddleware()
	handler = loggingMiddleware.Log(handler)
	handler = middleware.RequestIDMiddleware(handler)
	return handler
}

type applicationServices struct {
	userService      *userservice.Service
	roleService      *roleservice.RoleService
	permService      *permissionservice.PermissionService
	resourceService  *resourceservice.Service
	basicAuthService *basicauthservice.Service
	customerService  *customerservice.CustomerService
	supplierService  *supplierservice.SupplierService
	productService   *productservice.Service
	orderService     *orderservice.Service
	referralService  *referralservice.ReferralService
}

type applicationHandlers struct {
	userHandler     *userhttp.UserHandler
	roleHandler     *rolehttp.RoleHandler
	permHandler      *permissionhttp.PermissionHandler
	customerHandler *customerhttp.CustomerHandler
	supplierHandler *supplierhttp.SupplierHandler
	productHandler  *producthttp.ProductHandler
	orderHandler    *orderhttp.OrderHandler
	referralHandler *referralhttp.ReferralHandler
}

type applicationInfrastructure struct {
	cfg    *config.Config
	db     *sql.DB
	server http.Server
}

func provideApplicationInfrastructure(cfg *config.Config, dbConn *sql.DB, server http.Server) applicationInfrastructure {
	return applicationInfrastructure{
		cfg:    cfg,
		db:     dbConn,
		server: server,
	}
}

func provideApplicationServices(userService *userservice.Service, roleService *roleservice.RoleService, permService *permissionservice.PermissionService, resourceService *resourceservice.Service, basicAuthService *basicauthservice.Service, customerService *customerservice.CustomerService, supplierService *supplierservice.SupplierService, productService *productservice.Service, orderService *orderservice.Service, referralService *referralservice.ReferralService) applicationServices {
	return applicationServices{
		userService:      userService,
		roleService:      roleService,
		permService:      permService,
		resourceService:  resourceService,
		basicAuthService: basicAuthService,
		customerService:  customerService,
		supplierService:  supplierService,
		productService:   productService,
		orderService:     orderService,
		referralService:  referralService,
	}
}

func provideApplicationHandlers(userHandler *userhttp.UserHandler, roleHandler *rolehttp.RoleHandler, permHandler *permissionhttp.PermissionHandler, customerHandler *customerhttp.CustomerHandler, supplierHandler *supplierhttp.SupplierHandler, productHandler *producthttp.ProductHandler, orderHandler *orderhttp.OrderHandler, referralHandler *referralhttp.ReferralHandler) applicationHandlers {
	return applicationHandlers{
		userHandler:     userHandler,
		roleHandler:     roleHandler,
		permHandler:     permHandler,
		customerHandler: customerHandler,
		supplierHandler: supplierHandler,
		productHandler:  productHandler,
		orderHandler:    orderHandler,
		referralHandler: referralHandler,
	}
}

func provideApplication(infra applicationInfrastructure, services applicationServices, handlers applicationHandlers) *Application {
	return &Application{
		Config:           infra.cfg,
		DB:               infra.db,
		Server:           infra.server,
		UserService:      services.userService,
		UserHandler:      handlers.userHandler,
		RoleService:      services.roleService,
		RoleHandler:      handlers.roleHandler,
		PermService:      services.permService,
		PermHandler:      handlers.permHandler,
		ResourceService:  services.resourceService,
		BasicAuthService: services.basicAuthService,
		CustomerService:  services.customerService,
		CustomerHandler:  handlers.customerHandler,
		SupplierService:  services.supplierService,
		SupplierHandler:  handlers.supplierHandler,
		ProductService:   services.productService,
		ProductHandler:   handlers.productHandler,
		OrderService:     services.orderService,
		OrderHandler:     handlers.orderHandler,
		ReferralService:  services.referralService,
		ReferralHandler:  handlers.referralHandler,
	}
}
