//go:build wirelegacy
// +build wirelegacy

package app

import (
	"context"
	"marketplace/internal/config"
	ledgerhttp "marketplace/internal/core/accounting/ledger/api/http"
	ledgerrepo "marketplace/internal/core/accounting/ledger/repository"
	ledgerservice "marketplace/internal/core/accounting/ledger/service"
	mappinghttp "marketplace/internal/core/accounting/mapping/api/http"
	mappingrepo "marketplace/internal/core/accounting/mapping/repository"
	mappingservice "marketplace/internal/core/accounting/mapping/service"
	periodhttp "marketplace/internal/core/accounting/period/api/http"
	periodrepo "marketplace/internal/core/accounting/period/repository"
	periodservice "marketplace/internal/core/accounting/period/service"
	provisionhttp "marketplace/internal/core/accounting/provision/api/http"
	provisionrepo "marketplace/internal/core/accounting/provision/repository"
	provisionservice "marketplace/internal/core/accounting/provision/service"
	clienthttp "marketplace/internal/core/member/client/api/http"
	clientrepo "marketplace/internal/core/member/client/repository"
	clientservice "marketplace/internal/core/member/client/service"
	branchhttp "marketplace/internal/core/organization/branch/api/http"
	branchrepo "marketplace/internal/core/organization/branch/repository"
	branchservice "marketplace/internal/core/organization/branch/service"
	staffhttp "marketplace/internal/core/organization/staff/api/http"
	staffrepo "marketplace/internal/core/organization/staff/repository"
	staffservice "marketplace/internal/core/organization/staff/service"

	"database/sql"
	customerhttp "marketplace/internal/core/customer/api/http"
	customerrepo "marketplace/internal/core/customer/repository"
	customerservice "marketplace/internal/core/customer/service"
	supplierrepo "marketplace/internal/core/supplier/repository"
	supplierservice "marketplace/internal/core/supplier/service"
	chargehttp "marketplace/internal/core/portfolio/charge/api/http"
	chargerepo "marketplace/internal/core/portfolio/charge/repository"
	chargeservice "marketplace/internal/core/portfolio/charge/service"
	depositaccounthttp "marketplace/internal/core/portfolio/deposit/account/api/http"
	depositaccountrepo "marketplace/internal/core/portfolio/deposit/account/repository"
	depositaccountservice "marketplace/internal/core/portfolio/deposit/account/service"
	depositproducthttp "marketplace/internal/core/portfolio/deposit/product/api/http"
	depositproductrepo "marketplace/internal/core/portfolio/deposit/product/repository"
	depositproductservice "marketplace/internal/core/portfolio/deposit/product/service"
	institutionalaccounthttp "marketplace/internal/core/portfolio/institutional/account/api/http"
	institutionalaccountrepo "marketplace/internal/core/portfolio/institutional/account/repository"
	institutionalaccountservice "marketplace/internal/core/portfolio/institutional/account/service"
	loanaccounthttp "marketplace/internal/core/portfolio/loan/account/api/http"
	loanaccountdomain "marketplace/internal/core/portfolio/loan/account/domain"
	loanaccountrepo "marketplace/internal/core/portfolio/loan/account/repository"
	loanaccountservice "marketplace/internal/core/portfolio/loan/account/service"
	collateralhttp "marketplace/internal/core/portfolio/loan/collateral/api/http"
	collateralrepo "marketplace/internal/core/portfolio/loan/collateral/repository"
	collateralservice "marketplace/internal/core/portfolio/loan/collateral/service"
	loanproducthttp "marketplace/internal/core/portfolio/loan/product/api/http"
	loanproductrepo "marketplace/internal/core/portfolio/loan/product/repository"
	loanproductservice "marketplace/internal/core/portfolio/loan/product/service"
	savingsaccounthttp "marketplace/internal/core/portfolio/saving/account/api/http"
	savingsaccountrepo "marketplace/internal/core/portfolio/saving/account/repository"
	savingsaccountservice "marketplace/internal/core/portfolio/saving/account/service"
	savingsproducthttp "marketplace/internal/core/portfolio/saving/product/api/http"
	savingsproductrepo "marketplace/internal/core/portfolio/saving/product/repository"
	savingsproductservice "marketplace/internal/core/portfolio/saving/product/service"
	shareaccounthttp "marketplace/internal/core/portfolio/share/account/api/http"
	shareaccountrepo "marketplace/internal/core/portfolio/share/account/repository"
	shareaccountservice "marketplace/internal/core/portfolio/share/account/service"
	shareproducthttp "marketplace/internal/core/portfolio/share/product/api/http"
	shareproductrepo "marketplace/internal/core/portfolio/share/product/repository"
	shareproductservice "marketplace/internal/core/portfolio/share/product/service"
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
	stdhttp "net/http"

	"github.com/google/wire"
)

// InitializeApp wires up the application dependencies
func InitializeApp(cfg *config.Config) (*Application, error) {
	wire.Build(
		// Database provider
		provideDB,
		provideTxManager,
		// Resource catalog (service only, no HTTP handler)
		provideResourceRepository,
		provideResourceService,
		// Role module providers (needed by user service)
		provideRoleRepository,
		provideRoleService,
		provideRoleHandler,
		// User module providers
		provideUserRepository,
		providePhoneValidationRepository,
		providePhoneValidationService,
		provideRegistrationTokenRepository,
		provideRegistrationTokenService,
		provideUserService,
		provideRegistrationTokenAdapter,
		provideUserHandler,
		// Customer module providers
		provideCustomerRepository,
		provideCustomerService,
		provideCustomerHandler,
		// Loan product module providers
		provideLoanProductRepository,
		provideLoanProductService,
		provideLoanProductHandler,
		// Savings product module providers
		provideSavingsProductRepository,
		provideSavingsProductService,
		provideSavingsProductDropdownReadService,
		provideSavingsProductReadService,
		provideSavingsProductHandler,
		// Deposit product module providers
		provideDepositProductRepository,
		provideDepositProductService,
		provideDepositProductHandler,
		// Loan account module providers
		provideLoanAccountRepository,
		provideLoanAccountTransactionRepository,
		provideLoanAccountScheduleRepository,
		provideLoanAccountChargeRepository,
		provideLoanAccountScheduleCalculator,
		provideLoanAccountTransactionProcessor,
		provideLoanAccountService,
		provideLoanAccountHandler,
		// Savings account module providers
		provideSavingsAccountRepository,
		provideSavingsAccountTransactionRepository,
		provideSavingsAccountService,
		provideSavingsAccountHandler,
		// Deposit account module providers
		provideDepositAccountRepository,
		provideDepositAccountService,
		provideDepositAccountHandler,
		// Share product module providers
		provideShareProductRepository,
		provideShareProductService,
		provideShareProductHandler,
		// Share account module providers
		provideShareAccountRepository,
		provideShareAccountTransactionRepository,
		provideShareAccountService,
		provideShareAccountHandler,
		// Institutional account module providers
		provideInstitutionalAccountRepository,
		provideInstitutionalAccountTransactionRepository,
		provideInstitutionalAccountService,
		provideInstitutionalAccountHandler,
		// Branch module providers
		provideBranchRepository,
		provideBranchService,
		provideBranchHandler,
		provideStaffRepository,
		provideStaffService,
		provideStaffHandler,
		// GL account module providers
		provideGLAccountRepository,
		provideGLAccountService,
		provideGLAccountHandler,
		// Account-to-GL mapping module providers
		provideAccountGLMappingRepository,
		provideAccountGLMappingService,
		provideAccountGLMappingHandler,
		// Period module providers
		providePeriodRepository,
		providePeriodService,
		providePeriodHandler,
		// Provision module providers
		provideProvisionRepository,
		provideProvisionService,
		provideProvisionHandler,
		// Loan collateral module providers
		provideCollateralRepository,
		provideCodeValueRepository,
		provideCollateralService,
		provideCollateralHandler,
		// Permission module providers
		providePermissionRepository,
		providePermissionService,
		providePermissionHandler,
		// Charge module providers
		provideChargeRepository,
		provideChargeService,
		provideChargeHandler,
		// Client module providers
		provideClientRepository,
		provideClientService,
		provideClientHandler,
		// Idempotency infrastructure
		provideIdempotencyStore,
		provideIdempotencyMiddleware,
		// Basic auth infrastructure
		provideBasicAuthRepository,
		provideUserPermissionChecker,
		provideBasicAuthService,
		provideBasicAuthAdapter,
		provideBasicAuthHandler,
		// Registration token middleware
		provideRegistrationTokenMiddleware,
		// Auth and authorization middleware
		provideAuthMiddleware,
		provideAuthorizationMiddleware,
		// HTTP handler grouping
		provideAuthHandlers,
		provideProductHandlers,
		provideAccountHandlers,
		provideOrganizationHandlers,
		provideHTTPHandlerParams,
		provideHTTPHandlers,
		provideHTTPMiddleware,
		// HTTP handler + server
		provideHTTPHandler,
		// HTTP server provider
		provideHTTPServer,
		// Application grouping
		provideApplicationInfrastructure,
		provideApplicationServices,
		provideApplicationHandlers,
		// Application provider
		provideApplication,
	)
	return nil, nil
}

// provideDB creates a database connection
func provideDB(cfg *config.Config) (*sql.DB, error) {
	return db.NewPostgres(cfg.DB.DSN())
}

// provideTxManager creates a transaction manager for the provided DB.
func provideTxManager(dbConn *sql.DB) db.TxManager {
	return db.NewTxManager(dbConn)
}

// provideHTTPServer creates an HTTP server
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

// provideUserRepository creates a user repository
func provideUserRepository(db *sql.DB) *userrepository.Repository {
	return userrepository.NewRepository(db)
}

// providePhoneValidationRepository creates a phone validation repository
func providePhoneValidationRepository(db *sql.DB) userservice.PhoneValidationRepository {
	return userrepository.NewPhoneValidationRepository(db)
}

// providePhoneValidationService creates a phone validation service
func providePhoneValidationService(repo userservice.PhoneValidationRepository) *userservice.PhoneValidationService {
	return userservice.NewPhoneValidationService(repo)
}

// provideUserService creates a user service
func provideUserService(repo *userrepository.Repository, roleRepo *rolerepo.RoleRepository, phoneValidation *userservice.PhoneValidationService, registrationTokenService *userservice.RegistrationTokenService) *userservice.Service {
	return userservice.NewService(repo, roleRepo, phoneValidation, registrationTokenService)
}

// provideRegistrationTokenRepository creates a registration token repository
func provideRegistrationTokenRepository(db *sql.DB) userservice.RegistrationTokenRepository {
	return userrepository.NewRegistrationTokenRepository(db)
}

// provideRegistrationTokenService creates a registration token service
func provideRegistrationTokenService(repo userservice.RegistrationTokenRepository) *userservice.RegistrationTokenService {
	return userservice.NewRegistrationTokenService(repo)
}

// provideRegistrationTokenAdapter creates a registration token adapter that translates
// user domain errors to auth module errors
func provideRegistrationTokenAdapter(service *userservice.RegistrationTokenService) *basicauth.RegistrationTokenAdapter {
	return basicauth.NewRegistrationTokenAdapter(service)
}

// provideUserHandler creates a user HTTP handler
func provideUserHandler(userService *userservice.Service, permissionChecker *userservice.UserPermissionChecker, supplierService *supplierservice.SupplierService) *userhttp.UserHandler {
	return userhttp.NewUserHandler(userService, permissionChecker, supplierService)
}

// provideCustomerRepository creates a customer repository
func provideCustomerRepository(db *sql.DB) *customerrepo.CustomerRepository {
	return customerrepo.NewCustomerRepository(db)
}

// provideCustomerService creates a customer service
func provideCustomerService(repo *customerrepo.CustomerRepository) *customerservice.CustomerService {
	return customerservice.NewCustomerService(repo)
}

// provideCustomerHandler creates a customer HTTP handler
func provideCustomerHandler(service *customerservice.CustomerService) *customerhttp.CustomerHandler {
	return customerhttp.NewCustomerHandler(service)
}

// httpHandlers groups all HTTP handlers together
type httpHandlers struct {
	userHandler                 *userhttp.UserHandler
	roleHandler                 *rolehttp.RoleHandler
	permissionHandler           *permissionhttp.PermissionHandler
	basicAuthHandler            *basicauthhttp.Handler
	customerHandler             *customerhttp.CustomerHandler
	chargeHandler               *chargehttp.ChargeHandler
	clientHandler               *clienthttp.ClientHandler
	loanProductHandler          *loanproducthttp.LoanProductHandler
	savingsProductHandler       *savingsproducthttp.SavingsProductHandler
	depositProductHandler       *depositproducthttp.FixedDepositProductHandler
	shareProductHandler         *shareproducthttp.ShareProductHandler
	loanAccountHandler          *loanaccounthttp.LoanAccountHandler
	savingsAccountHandler       *savingsaccounthttp.SavingsAccountHandler
	depositAccountHandler       *depositaccounthttp.FixedDepositAccountHandler
	shareAccountHandler         *shareaccounthttp.ShareAccountHandler
	institutionalAccountHandler *institutionalaccounthttp.InstitutionalAccountHandler
	branchHandler               *branchhttp.BranchHandler
	// stuffHandler removed - functionality moved to asset module
	staffHandler            *staffhttp.StaffHandler
	collateralHandler       *collateralhttp.CollateralHandler
	glAccountHandler        *ledgerhttp.GLAccountHandler
	accountGLMappingHandler *mappinghttp.AccountGLMappingHandler
	periodHandler           *periodhttp.PeriodHandler
	provisionHandler        *provisionhttp.ProvisionHandler
}

// httpMiddleware groups all HTTP middleware together
type httpMiddleware struct {
	idempotencyMiddleware       *middleware.IdempotencyMiddleware
	authMiddleware              *middleware.AuthMiddleware
	authorizationMiddleware     *middleware.AuthorizationMiddleware
	registrationTokenMiddleware *middleware.RegistrationTokenMiddleware
}

// authHandlers groups authentication and authorization handlers
type authHandlers struct {
	userHandler       *userhttp.UserHandler
	roleHandler       *rolehttp.RoleHandler
	permissionHandler *permissionhttp.PermissionHandler
	basicAuthHandler  *basicauthhttp.Handler
}

// productHandlers groups product domain handlers
type productHandlers struct {
	loanProductHandler    *loanproducthttp.LoanProductHandler
	savingsProductHandler *savingsproducthttp.SavingsProductHandler
	depositProductHandler *depositproducthttp.FixedDepositProductHandler
	shareProductHandler   *shareproducthttp.ShareProductHandler
}

// accountHandlers groups account domain handlers
type accountHandlers struct {
	loanAccountHandler          *loanaccounthttp.LoanAccountHandler
	savingsAccountHandler       *savingsaccounthttp.SavingsAccountHandler
	depositAccountHandler       *depositaccounthttp.FixedDepositAccountHandler
	shareAccountHandler         *shareaccounthttp.ShareAccountHandler
	institutionalAccountHandler *institutionalaccounthttp.InstitutionalAccountHandler
}

// organizationHandlers groups organization domain handlers
type organizationHandlers struct {
	branchHandler *branchhttp.BranchHandler
	// stuffHandler removed - functionality moved to asset module
	staffHandler *staffhttp.StaffHandler
}

// httpHandlerParams groups handler parameters to reduce function parameter count
type httpHandlerParams struct {
	userHandler                 *userhttp.UserHandler
	roleHandler                 *rolehttp.RoleHandler
	permissionHandler           *permissionhttp.PermissionHandler
	basicAuthHandler            *basicauthhttp.Handler
	customerHandler             *customerhttp.CustomerHandler
	chargeHandler               *chargehttp.ChargeHandler
	clientHandler               *clienthttp.ClientHandler
	loanProductHandler          *loanproducthttp.LoanProductHandler
	savingsProductHandler       *savingsproducthttp.SavingsProductHandler
	depositProductHandler       *depositproducthttp.FixedDepositProductHandler
	shareProductHandler         *shareproducthttp.ShareProductHandler
	loanAccountHandler          *loanaccounthttp.LoanAccountHandler
	savingsAccountHandler       *savingsaccounthttp.SavingsAccountHandler
	depositAccountHandler       *depositaccounthttp.FixedDepositAccountHandler
	shareAccountHandler         *shareaccounthttp.ShareAccountHandler
	institutionalAccountHandler *institutionalaccounthttp.InstitutionalAccountHandler
	branchHandler               *branchhttp.BranchHandler
	staffHandler                *staffhttp.StaffHandler
	glAccountHandler            *ledgerhttp.GLAccountHandler
	accountGLMappingHandler     *mappinghttp.AccountGLMappingHandler
	periodHandler               *periodhttp.PeriodHandler
	provisionHandler            *provisionhttp.ProvisionHandler
}

// provideAuthHandlers creates the authHandlers struct
func provideAuthHandlers(
	userHandler *userhttp.UserHandler,
	roleHandler *rolehttp.RoleHandler,
	permissionHandler *permissionhttp.PermissionHandler,
	basicAuthHandler *basicauthhttp.Handler,
) authHandlers {
	return authHandlers{
		userHandler:       userHandler,
		roleHandler:       roleHandler,
		permissionHandler: permissionHandler,
		basicAuthHandler:  basicAuthHandler,
	}
}

// provideProductHandlers creates the productHandlers struct
func provideProductHandlers(
	loanProductHandler *loanproducthttp.LoanProductHandler,
	savingsProductHandler *savingsproducthttp.SavingsProductHandler,
	depositProductHandler *depositproducthttp.FixedDepositProductHandler,
	shareProductHandler *shareproducthttp.ShareProductHandler,
) productHandlers {
	return productHandlers{
		loanProductHandler:    loanProductHandler,
		savingsProductHandler: savingsProductHandler,
		depositProductHandler: depositProductHandler,
		shareProductHandler:   shareProductHandler,
	}
}

// provideAccountHandlers creates the accountHandlers struct
func provideAccountHandlers(
	loanAccountHandler *loanaccounthttp.LoanAccountHandler,
	savingsAccountHandler *savingsaccounthttp.SavingsAccountHandler,
	depositAccountHandler *depositaccounthttp.FixedDepositAccountHandler,
	shareAccountHandler *shareaccounthttp.ShareAccountHandler,
	institutionalAccountHandler *institutionalaccounthttp.InstitutionalAccountHandler,
) accountHandlers {
	return accountHandlers{
		loanAccountHandler:          loanAccountHandler,
		savingsAccountHandler:       savingsAccountHandler,
		depositAccountHandler:       depositAccountHandler,
		shareAccountHandler:         shareAccountHandler,
		institutionalAccountHandler: institutionalAccountHandler,
	}
}

// provideOrganizationHandlers creates the organizationHandlers struct
func provideOrganizationHandlers(
	branchHandler *branchhttp.BranchHandler,
	// stuffHandler removed - functionality moved to asset module
	staffHandler *staffhttp.StaffHandler,
) organizationHandlers {
	return organizationHandlers{
		branchHandler: branchHandler,
		// stuffHandler removed
		staffHandler: staffHandler,
	}
}

// provideHTTPHandlerParams creates the httpHandlerParams struct
func provideHTTPHandlerParams(
	auth authHandlers,
	customerHandler *customerhttp.CustomerHandler,
	products productHandlers,
	accounts accountHandlers,
	org organizationHandlers,
	glAccountHandler *ledgerhttp.GLAccountHandler,
	accountGLMappingHandler *mappinghttp.AccountGLMappingHandler,
	periodHandler *periodhttp.PeriodHandler,
	provisionHandler *provisionhttp.ProvisionHandler,
	chargeHandler *chargehttp.ChargeHandler,
	clientHandler *clienthttp.ClientHandler,
) httpHandlerParams {
	return httpHandlerParams{
		userHandler:                 auth.userHandler,
		roleHandler:                 auth.roleHandler,
		permissionHandler:           auth.permissionHandler,
		basicAuthHandler:            auth.basicAuthHandler,
		customerHandler:             customerHandler,
		chargeHandler:               chargeHandler,
		clientHandler:               clientHandler,
		loanProductHandler:          products.loanProductHandler,
		savingsProductHandler:       products.savingsProductHandler,
		depositProductHandler:       products.depositProductHandler,
		shareProductHandler:         products.shareProductHandler,
		loanAccountHandler:          accounts.loanAccountHandler,
		savingsAccountHandler:       accounts.savingsAccountHandler,
		depositAccountHandler:       accounts.depositAccountHandler,
		shareAccountHandler:         accounts.shareAccountHandler,
		institutionalAccountHandler: accounts.institutionalAccountHandler,
		branchHandler:               org.branchHandler,
		// stuffHandler removed
		staffHandler:            org.staffHandler,
		glAccountHandler:        glAccountHandler,
		accountGLMappingHandler: accountGLMappingHandler,
		periodHandler:           periodHandler,
		provisionHandler:        provisionHandler,
	}
}

// provideHTTPHandlers creates the httpHandlers struct
func provideHTTPHandlers(params httpHandlerParams, collateralHandler *collateralhttp.CollateralHandler) httpHandlers {
	return httpHandlers{
		userHandler:                 params.userHandler,
		roleHandler:                 params.roleHandler,
		permissionHandler:           params.permissionHandler,
		basicAuthHandler:            params.basicAuthHandler,
		customerHandler:             params.customerHandler,
		chargeHandler:               params.chargeHandler,
		clientHandler:               params.clientHandler,
		loanProductHandler:          params.loanProductHandler,
		savingsProductHandler:       params.savingsProductHandler,
		depositProductHandler:       params.depositProductHandler,
		shareProductHandler:         params.shareProductHandler,
		loanAccountHandler:          params.loanAccountHandler,
		savingsAccountHandler:       params.savingsAccountHandler,
		depositAccountHandler:       params.depositAccountHandler,
		shareAccountHandler:         params.shareAccountHandler,
		institutionalAccountHandler: params.institutionalAccountHandler,
		branchHandler:               params.branchHandler,
		// stuffHandler removed
		staffHandler:            params.staffHandler,
		collateralHandler:       collateralHandler,
		glAccountHandler:        params.glAccountHandler,
		accountGLMappingHandler: params.accountGLMappingHandler,
		periodHandler:           params.periodHandler,
		provisionHandler:        params.provisionHandler,
	}
}

// provideHTTPMiddleware creates the httpMiddleware struct
func provideHTTPMiddleware(idempotencyMiddleware *middleware.IdempotencyMiddleware, authMiddleware *middleware.AuthMiddleware, authorizationMiddleware *middleware.AuthorizationMiddleware, registrationTokenMiddleware *middleware.RegistrationTokenMiddleware) httpMiddleware {
	return httpMiddleware{
		idempotencyMiddleware:       idempotencyMiddleware,
		authMiddleware:              authMiddleware,
		authorizationMiddleware:     authorizationMiddleware,
		registrationTokenMiddleware: registrationTokenMiddleware,
	}
}

// provideHTTPHandler registers all HTTP routes and returns a handler
func provideHTTPHandler(handlers httpHandlers, mw httpMiddleware) stdhttp.Handler {
	mux := stdhttp.NewServeMux()
	userhttp.RegisterHTTPRoutes(mux, handlers.userHandler)
	rolehttp.RegisterHTTPRoutes(mux, handlers.roleHandler)
	permissionhttp.RegisterHTTPRoutes(mux, handlers.permissionHandler)
	basicauthhttp.RegisterHTTPRoutes(mux, handlers.basicAuthHandler)
	customerhttp.RegisterHTTPRoutes(mux, handlers.customerHandler)
	chargehttp.RegisterHTTPRoutes(mux, handlers.chargeHandler)
	clienthttp.RegisterHTTPRoutes(mux, handlers.clientHandler)
	loanproducthttp.RegisterHTTPRoutes(mux, handlers.loanProductHandler)
	savingsproducthttp.RegisterHTTPRoutes(mux, handlers.savingsProductHandler)
	depositproducthttp.RegisterFixedDepositProductHTTPRoutes(mux, handlers.depositProductHandler)
	shareproducthttp.RegisterHTTPRoutes(mux, handlers.shareProductHandler)
	// Loan account routes also handle collateral routes (nested under accounts)
	loanaccounthttp.RegisterHTTPRoutes(mux, handlers.loanAccountHandler, handlers.collateralHandler)
	savingsaccounthttp.RegisterHTTPRoutes(mux, handlers.savingsAccountHandler)
	depositaccounthttp.RegisterFixedDepositAccountHTTPRoutes(mux, handlers.depositAccountHandler)
	shareaccounthttp.RegisterHTTPRoutes(mux, handlers.shareAccountHandler)
	institutionalaccounthttp.RegisterHTTPRoutes(mux, handlers.institutionalAccountHandler)
	branchhttp.RegisterHTTPRoutes(mux, handlers.branchHandler)
	staffhttp.RegisterHTTPRoutes(mux, handlers.staffHandler)
	ledgerhttp.RegisterHTTPRoutes(mux, handlers.glAccountHandler)
	mappinghttp.RegisterHTTPRoutes(mux, handlers.accountGLMappingHandler)
	periodhttp.RegisterHTTPRoutes(mux, handlers.periodHandler)
	provisionhttp.RegisterHTTPRoutes(mux, handlers.provisionHandler)

	// Chain middleware in order: request_id -> logging -> registration_token -> auth -> authorization -> idempotency -> handler
	// Request ID must be first so all subsequent middleware have access to it
	handler := mw.idempotencyMiddleware.Handle(mux)
	handler = mw.authorizationMiddleware.Handle(handler)
	handler = mw.authMiddleware.Handle(handler)
	handler = mw.registrationTokenMiddleware.Handle(handler)
	loggingMiddleware := middleware.NewLoggingMiddleware()
	handler = loggingMiddleware.Log(handler)
	handler = middleware.RequestIDMiddleware(handler)
	return handler
}

// applicationServices groups all application services together
type applicationServices struct {
	userService      *userservice.Service
	roleService      *roleservice.RoleService
	permService      *permissionservice.PermissionService
	resourceService  *resourceservice.Service
	basicAuthService *basicauthservice.Service
	customerService  *customerservice.CustomerService
}

// applicationHandlers groups all application HTTP handlers together
type applicationHandlers struct {
	userHandler     *userhttp.UserHandler
	roleHandler     *rolehttp.RoleHandler
	permHandler     *permissionhttp.PermissionHandler
	customerHandler *customerhttp.CustomerHandler
}

// applicationInfrastructure groups infrastructure components
type applicationInfrastructure struct {
	cfg    *config.Config
	db     *sql.DB
	server http.Server
}

// provideApplicationInfrastructure creates the applicationInfrastructure struct
func provideApplicationInfrastructure(cfg *config.Config, db *sql.DB, server http.Server) applicationInfrastructure {
	return applicationInfrastructure{
		cfg:    cfg,
		db:     db,
		server: server,
	}
}

// provideApplicationServices creates the applicationServices struct
func provideApplicationServices(userService *userservice.Service, roleService *roleservice.RoleService, permService *permissionservice.PermissionService, resourceService *resourceservice.Service, basicAuthService *basicauthservice.Service, customerService *customerservice.CustomerService) applicationServices {
	return applicationServices{
		userService:      userService,
		roleService:      roleService,
		permService:      permService,
		resourceService:  resourceService,
		basicAuthService: basicAuthService,
		customerService:  customerService,
	}
}

// provideApplicationHandlers creates the applicationHandlers struct
func provideApplicationHandlers(userHandler *userhttp.UserHandler, roleHandler *rolehttp.RoleHandler, permHandler *permissionhttp.PermissionHandler, customerHandler *customerhttp.CustomerHandler) applicationHandlers {
	return applicationHandlers{
		userHandler:     userHandler,
		roleHandler:     roleHandler,
		permHandler:     permHandler,
		customerHandler: customerHandler,
	}
}

// provideApplication creates the Application struct
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
	}
}

// provideRoleRepository creates a role repository
func provideRoleRepository(db *sql.DB, txManager db.TxManager) *rolerepo.RoleRepository {
	return rolerepo.NewRoleRepository(db, txManager)
}

// provideRoleService creates a role service
func provideRoleService(repo *rolerepo.RoleRepository, permRepo *permissionrepo.PermissionRepository) *roleservice.RoleService {
	return roleservice.NewRoleService(repo, permRepo)
}

// provideRoleHandler creates a role HTTP handler
func provideRoleHandler(roleService *roleservice.RoleService) *rolehttp.RoleHandler {
	return rolehttp.NewRoleHandler(roleService)
}

// providePermissionRepository creates a permission repository.
func providePermissionRepository(db *sql.DB) *permissionrepo.PermissionRepository {
	return permissionrepo.NewPermissionRepository(db)
}

// providePermissionService creates a permission service.
func providePermissionService(repo *permissionrepo.PermissionRepository, resourceService *resourceservice.Service) *permissionservice.PermissionService {
	return permissionservice.NewPermissionService(repo, resourceService)
}

// providePermissionHandler creates a permission HTTP handler.
func providePermissionHandler(service *permissionservice.PermissionService) *permissionhttp.PermissionHandler {
	return permissionhttp.NewPermissionHandler(service)
}

func provideResourceRepository(db *sql.DB, txManager db.TxManager) *resourcerepo.Repository {
	return resourcerepo.NewRepository(db, txManager)
}

func provideResourceService(repo *resourcerepo.Repository) *resourceservice.Service {
	return resourceservice.NewService(repo)
}

// provideIdempotencyStore wires the idempotency store.
func provideIdempotencyStore(db *sql.DB) *idempotency.Store {
	return idempotency.NewStore(db)
}

// provideIdempotencyMiddleware wires the idempotency middleware.
func provideIdempotencyMiddleware(store *idempotency.Store) *middleware.IdempotencyMiddleware {
	return middleware.NewIdempotencyMiddleware(store)
}

// provideBasicAuthRepository creates a basic auth credential repository
func provideBasicAuthRepository(dbConn *sql.DB) *basicauthrepo.CredentialRepository {
	return basicauthrepo.NewCredentialRepository(dbConn)
}

// provideUserPermissionChecker creates a permission checker for auth service
func provideUserPermissionChecker(userRepo *userrepository.Repository, roleService *roleservice.RoleService) *userservice.UserPermissionChecker {
	return userservice.NewUserPermissionChecker(userRepo, roleService)
}

// provideBasicAuthService creates a basic auth service
func provideBasicAuthService(repo *basicauthrepo.CredentialRepository, userService *userservice.Service, permissionChecker *userservice.UserPermissionChecker) *basicauthservice.Service {
	return basicauthservice.NewService(repo, userService, permissionChecker)
}

// provideBasicAuthAdapter creates an adapter for the basic auth service
func provideBasicAuthAdapter(service *basicauthservice.Service) *basicauth.Adapter {
	return basicauth.NewAdapter(service)
}

// provideBasicAuthHandler creates a basic auth credential handler
func provideBasicAuthHandler(service *basicauthservice.Service, userService *userservice.Service) *basicauthhttp.Handler {
	return basicauthhttp.NewHandler(service, userService, userService)
}

// provideRegistrationTokenMiddleware creates middleware for validating registration tokens
func provideRegistrationTokenMiddleware(adapter *basicauth.RegistrationTokenAdapter) *middleware.RegistrationTokenMiddleware {
	return middleware.NewRegistrationTokenMiddleware(adapter, []string{"/auth/basic/credentials"})
}

// collectPublicRoutes collects public routes from all registered modules.
// Modules automatically register their public routes via init() functions.
func collectPublicRoutes() []string {
	return middleware.CollectPublicRoutes()
}

// collectAuthenticatedRoutes collects authenticated routes from all registered modules.
// Modules automatically register their authenticated routes via init() functions.
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

	// Collect public routes from all modules
	publicRoutes := collectPublicRoutes()

	return middleware.NewAuthMiddleware(authDomain.AuthMode(cfg.Auth.Mode), realm, basicAuthService, userService, publicRoutes)
}

// provideAuthorizationMiddleware wires the authorization middleware.
func provideAuthorizationMiddleware(resourceService *resourceservice.Service, permissionChecker *userservice.UserPermissionChecker) *middleware.AuthorizationMiddleware {
	// Collect public routes from all modules (same as auth middleware)
	publicRoutes := collectPublicRoutes()
	// Collect authenticated routes from all modules (routes accessible to all authenticated users, bypass permission checks)
	allowedAuthenticatedRoutes := collectAuthenticatedRoutes()
	return middleware.NewAuthorizationMiddleware(resourceService, permissionChecker, publicRoutes, allowedAuthenticatedRoutes)
}

// provideLoanProductRepository creates a loan product repository
func provideLoanProductRepository(db *sql.DB) *loanproductrepo.Repository {
	return loanproductrepo.NewRepository(db)
}

// provideLoanProductService creates a loan product service
func provideLoanProductService(repo *loanproductrepo.Repository) *loanproductservice.Service {
	return loanproductservice.NewService(repo)
}

// provideLoanProductHandler creates a loan product HTTP handler
func provideLoanProductHandler(service *loanproductservice.Service) *loanproducthttp.LoanProductHandler {
	return loanproducthttp.NewLoanProductHandler(service)
}

// provideSavingsProductRepository creates a savings product repository
func provideSavingsProductRepository(db *sql.DB) *savingsproductrepo.Repository {
	return savingsproductrepo.NewRepository(db)
}

// provideSavingsProductService creates a savings product service
func provideSavingsProductService(repo *savingsproductrepo.Repository) *savingsproductservice.Service {
	return savingsproductservice.NewService(repo)
}

// provideSavingsProductDropdownReadService creates a dropdown read service for savings products
func provideSavingsProductDropdownReadService() savingsproductservice.DropdownReadService {
	return savingsproductservice.NewDropdownReadService()
}

// provideSavingsProductReadService creates a read service for savings products
func provideSavingsProductReadService(
	repo *savingsproductrepo.Repository,
	dropdownService savingsproductservice.DropdownReadService,
) savingsproductservice.ReadService {
	return savingsproductservice.NewReadService(repo, dropdownService)
}

// provideSavingsProductHandler creates a savings product HTTP handler
func provideSavingsProductHandler(
	service *savingsproductservice.Service,
	readService savingsproductservice.ReadService,
) *savingsproducthttp.SavingsProductHandler {
	return savingsproducthttp.NewSavingsProductHandler(service, readService)
}

// provideLoanAccountRepository creates a loan account repository
func provideLoanAccountRepository(db *sql.DB) *loanaccountrepo.Repository {
	return loanaccountrepo.NewRepository(db)
}

// provideLoanAccountTransactionRepository creates a loan account transaction repository
func provideLoanAccountTransactionRepository(db *sql.DB) *loanaccountrepo.TransactionRepository {
	return loanaccountrepo.NewTransactionRepository(db)
}

// provideLoanAccountScheduleRepository creates a loan account schedule repository
func provideLoanAccountScheduleRepository(db *sql.DB) *loanaccountrepo.ScheduleRepository {
	return loanaccountrepo.NewScheduleRepository(db)
}

// provideLoanAccountChargeRepository creates a loan account charge repository
func provideLoanAccountChargeRepository(db *sql.DB) *loanaccountrepo.ChargeRepository {
	return loanaccountrepo.NewChargeRepository(db)
}

// provideLoanAccountScheduleCalculator creates a loan account schedule calculator
func provideLoanAccountScheduleCalculator() *loanaccountdomain.DefaultScheduleCalculator {
	return loanaccountdomain.NewDefaultScheduleCalculator()
}

// provideLoanAccountTransactionProcessor creates a loan account transaction processor
func provideLoanAccountTransactionProcessor() *loanaccountdomain.DefaultTransactionProcessor {
	return loanaccountdomain.NewDefaultTransactionProcessor()
}

// provideLoanAccountService creates a loan account service
func provideLoanAccountService(
	repo *loanaccountrepo.Repository,
	transactionRepo *loanaccountrepo.TransactionRepository,
	scheduleRepo *loanaccountrepo.ScheduleRepository,
	chargeRepo *loanaccountrepo.ChargeRepository,
	scheduleCalculator *loanaccountdomain.DefaultScheduleCalculator,
	transactionProcessor *loanaccountdomain.DefaultTransactionProcessor,
) *loanaccountservice.Service {
	return loanaccountservice.NewService(
		repo,
		transactionRepo,
		scheduleRepo,
		chargeRepo,
		scheduleCalculator,
		transactionProcessor,
	)
}

// provideLoanAccountHandler creates a loan account HTTP handler
func provideLoanAccountHandler(service *loanaccountservice.Service) *loanaccounthttp.LoanAccountHandler {
	return loanaccounthttp.NewLoanAccountHandler(service)
}

// provideSavingsAccountRepository creates a savings account repository
func provideSavingsAccountRepository(db *sql.DB) *savingsaccountrepo.Repository {
	return savingsaccountrepo.NewRepository(db)
}

// provideSavingsAccountTransactionRepository creates a savings account transaction repository
func provideSavingsAccountTransactionRepository(db *sql.DB) *savingsaccountrepo.TransactionRepository {
	return savingsaccountrepo.NewTransactionRepository(db)
}

// provideSavingsAccountInterestCalculator creates a savings account interest calculator
func provideSavingsAccountInterestCalculator(transactionRepo *savingsaccountrepo.TransactionRepository) *savingsaccountservice.InterestCalculator {
	return savingsaccountservice.NewInterestCalculator(transactionRepo)
}

// provideSavingsAccountInterestService creates a savings account interest service
func provideSavingsAccountInterestService(
	service *savingsaccountservice.Service,
	interestCalc *savingsaccountservice.InterestCalculator,
	productRepo *savingsproductrepo.Repository,
) *savingsaccountservice.InterestService {
	return savingsaccountservice.NewInterestService(service, interestCalc, productRepo)
}

// provideSavingsAccountService creates a savings account service
func provideSavingsAccountService(
	repo *savingsaccountrepo.Repository,
	transactionRepo *savingsaccountrepo.TransactionRepository,
	productRepo *savingsproductrepo.Repository,
	txManager db.TxManager,
) *savingsaccountservice.Service {
	return savingsaccountservice.NewService(repo, transactionRepo, productRepo, txManager)
}

// provideSavingsAccountHandler creates a savings account HTTP handler
func provideSavingsAccountHandler(service *savingsaccountservice.Service) *savingsaccounthttp.SavingsAccountHandler {
	return savingsaccounthttp.NewSavingsAccountHandler(service)
}

// provideInstitutionalAccountRepository creates an institutional account repository
func provideInstitutionalAccountRepository(db *sql.DB) *institutionalaccountrepo.Repository {
	return institutionalaccountrepo.NewRepository(db)
}

// provideInstitutionalAccountTransactionRepository creates an institutional account transaction repository
func provideInstitutionalAccountTransactionRepository(db *sql.DB) *institutionalaccountrepo.TransactionRepository {
	return institutionalaccountrepo.NewTransactionRepository(db)
}

// provideInstitutionalAccountService creates an institutional account service
func provideInstitutionalAccountService(
	repo *institutionalaccountrepo.Repository,
	transactionRepo *institutionalaccountrepo.TransactionRepository,
	txManager db.TxManager,
) *institutionalaccountservice.Service {
	return institutionalaccountservice.NewService(repo, transactionRepo, txManager)
}

// provideInstitutionalAccountHandler creates an institutional account HTTP handler
func provideInstitutionalAccountHandler(service *institutionalaccountservice.Service) *institutionalaccounthttp.InstitutionalAccountHandler {
	return institutionalaccounthttp.NewInstitutionalAccountHandler(service)
}

// provideBranchRepository creates a branch repository
func provideBranchRepository(db *sql.DB) *branchrepo.Repository {
	return branchrepo.NewRepository(db)
}

// provideBranchService creates a branch service
func provideBranchService(
	repo *branchrepo.Repository,
	txManager db.TxManager,
) *branchservice.Service {
	return branchservice.NewService(repo, txManager)
}

// provideBranchHandler creates a branch HTTP handler
func provideBranchHandler(service *branchservice.Service) *branchhttp.BranchHandler {
	return branchhttp.NewBranchHandler(service)
}

// provideStuffRepository - removed (functionality moved to asset module)
// func provideStuffRepository(db *sql.DB) *stufrepo.Repository {
// 	return stufrepo.NewRepository(db)
// }

// provideStuffService - removed (functionality moved to asset module)
// func provideStuffService(
// 	repo *stufrepo.Repository,
// 	txManager db.TxManager,
// ) *stufservice.Service {
// 	return stufservice.NewService(repo, txManager)
// }

// provideStuffHandler - removed (functionality moved to asset module)
// func provideStuffHandler(service *stufservice.Service) *stufhttp.StuffHandler {
// 	return stufhttp.NewStuffHandler(service)
// }

// provideStaffRepository creates a staff repository
func provideStaffRepository(db *sql.DB) *staffrepo.Repository {
	return staffrepo.NewRepository(db)
}

// provideStaffService creates a staff service
func provideStaffService(
	repo *staffrepo.Repository,
	txManager db.TxManager,
) *staffservice.Service {
	return staffservice.NewService(repo, txManager)
}

// provideStaffHandler creates a staff HTTP handler
func provideStaffHandler(service *staffservice.Service) *staffhttp.StaffHandler {
	return staffhttp.NewStaffHandler(service)
}

// provideCollateralRepository creates a collateral repository
func provideCollateralRepository(db *sql.DB) *collateralrepo.Repository {
	return collateralrepo.NewRepository(db)
}

// provideCodeValueRepository creates a code value repository
// Note: This is a placeholder. In production, this should be provided by the organization/code module
func provideCodeValueRepository(db *sql.DB) collateralservice.CodeValueRepository {
	return collateralrepo.NewCodeValueRepository(db)
}

// provideCollateralService creates a collateral service
func provideCollateralService(
	collateralRepo *collateralrepo.Repository,
	codeValueRepo collateralservice.CodeValueRepository,
	loanRepo *loanaccountrepo.Repository,
) *collateralservice.Service {
	// Create an adapter that implements LoanRepository interface
	loanRepoAdapter := &loanRepositoryAdapter{repo: loanRepo}
	return collateralservice.NewService(collateralRepo, codeValueRepo, loanRepoAdapter)
}

// loanRepositoryAdapter adapts loan account repository to collateral service's LoanRepository interface
type loanRepositoryAdapter struct {
	repo *loanaccountrepo.Repository
}

func (a *loanRepositoryAdapter) FindByID(ctx context.Context, id string) (*collateralservice.LoanInfo, error) {
	loan, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if loan == nil {
		return nil, nil
	}
	// Convert domain object to minimal LoanInfo to maintain domain boundaries
	return &collateralservice.LoanInfo{
		ID:                loan.ID,
		Status:            string(loan.Status),
		ApprovedPrincipal: loan.ApprovedPrincipal,
	}, nil
}

// provideCollateralHandler creates a collateral HTTP handler
func provideCollateralHandler(service *collateralservice.Service) *collateralhttp.CollateralHandler {
	return collateralhttp.NewCollateralHandler(service)
}

// provideShareProductRepository creates a share product repository
func provideShareProductRepository(db *sql.DB) *shareproductrepo.Repository {
	return shareproductrepo.NewRepository(db)
}

// provideShareProductService creates a share product service
func provideShareProductService(repo *shareproductrepo.Repository) *shareproductservice.Service {
	return shareproductservice.NewService(repo)
}

// provideShareProductHandler creates a share product HTTP handler
func provideShareProductHandler(
	service *shareproductservice.Service,
) *shareproducthttp.ShareProductHandler {
	return shareproducthttp.NewShareProductHandler(service)
}

// provideShareAccountRepository creates a share account repository
func provideShareAccountRepository(db *sql.DB) *shareaccountrepo.Repository {
	return shareaccountrepo.NewRepository(db)
}

// provideShareAccountTransactionRepository creates a share account transaction repository
func provideShareAccountTransactionRepository(db *sql.DB) *shareaccountrepo.TransactionRepository {
	return shareaccountrepo.NewTransactionRepository(db)
}

// provideShareAccountService creates a share account service
func provideShareAccountService(
	repo *shareaccountrepo.Repository,
	transactionRepo *shareaccountrepo.TransactionRepository,
	productRepo *shareproductrepo.Repository,
	txManager db.TxManager,
) *shareaccountservice.Service {
	return shareaccountservice.NewService(repo, transactionRepo, productRepo, txManager)
}

// provideShareAccountHandler creates a share account HTTP handler
func provideShareAccountHandler(service *shareaccountservice.Service) *shareaccounthttp.ShareAccountHandler {
	return shareaccounthttp.NewShareAccountHandler(service)
}

// provideDepositProductRepository creates a deposit product repository
func provideDepositProductRepository(db *sql.DB) depositproductservice.Repository {
	return depositproductrepo.NewRepository(db)
}

// provideDepositProductService creates a deposit product service
func provideDepositProductService(
	repo depositproductservice.Repository,
	savingsProductRepo *savingsproductrepo.Repository,
) *depositproductservice.Service {
	return depositproductservice.NewService(repo, savingsProductRepo)
}

// provideDepositProductHandler creates a deposit product HTTP handler
func provideDepositProductHandler(service *depositproductservice.Service) *depositproducthttp.FixedDepositProductHandler {
	return depositproducthttp.NewFixedDepositProductHandler(service)
}

// provideDepositAccountRepository creates a deposit account repository
func provideDepositAccountRepository(db *sql.DB) depositaccountservice.Repository {
	return depositaccountrepo.NewRepository(db)
}

// provideDepositAccountService creates a deposit account service
func provideDepositAccountService(
	repo depositaccountservice.Repository,
	productRepo depositproductservice.Repository,
	savingsAccountRepo *savingsaccountrepo.Repository,
) *depositaccountservice.Service {
	return depositaccountservice.NewService(repo, productRepo, savingsAccountRepo)
}

// provideDepositAccountHandler creates a deposit account HTTP handler
func provideDepositAccountHandler(service *depositaccountservice.Service) *depositaccounthttp.FixedDepositAccountHandler {
	return depositaccounthttp.NewFixedDepositAccountHandler(service)
}

// provideChargeRepository creates a charge repository
func provideChargeRepository(db *sql.DB) chargeservice.Repository {
	return chargerepo.NewRepository(db)
}

// provideChargeService creates a charge service
func provideChargeService(
	repo chargeservice.Repository,
) *chargeservice.Service {
	return chargeservice.NewService(repo)
}

// provideChargeHandler creates a charge HTTP handler
func provideChargeHandler(service *chargeservice.Service) *chargehttp.ChargeHandler {
	return chargehttp.NewChargeHandler(service)
}

// provideClientRepository creates a client repository
func provideClientRepository(db *sql.DB) clientservice.Repository {
	return clientrepo.NewRepository(db)
}

// provideClientService creates a client service
func provideClientService(repo clientservice.Repository) *clientservice.Service {
	return clientservice.NewService(repo)
}

// provideClientHandler creates a client HTTP handler
func provideClientHandler(service *clientservice.Service) *clienthttp.ClientHandler {
	return clienthttp.NewClientHandler(service)
}

// provideGLAccountRepository creates a GL account repository
func provideGLAccountRepository(db *sql.DB) *ledgerrepo.Repository {
	return ledgerrepo.NewRepository(db)
}

// provideGLAccountService creates a GL account service
func provideGLAccountService(repo *ledgerrepo.Repository) *ledgerservice.Service {
	return ledgerservice.NewService(repo)
}

// provideGLAccountHandler creates a GL account HTTP handler
func provideGLAccountHandler(service *ledgerservice.Service) *ledgerhttp.GLAccountHandler {
	return ledgerhttp.NewGLAccountHandler(service)
}

// provideAccountGLMappingRepository creates an account-to-GL mapping repository
func provideAccountGLMappingRepository(db *sql.DB) *mappingrepo.Repository {
	return mappingrepo.NewRepository(db)
}

// provideAccountGLMappingService creates an account-to-GL mapping service
func provideAccountGLMappingService(repo *mappingrepo.Repository, glAccountRepo *ledgerrepo.Repository) *mappingservice.Service {
	return mappingservice.NewService(repo, glAccountRepo)
}

// provideAccountGLMappingHandler creates an account-to-GL mapping HTTP handler
func provideAccountGLMappingHandler(service *mappingservice.Service) *mappinghttp.AccountGLMappingHandler {
	return mappinghttp.NewAccountGLMappingHandler(service)
}

// providePeriodRepository creates a period repository
func providePeriodRepository(db *sql.DB) *periodrepo.Repository {
	return periodrepo.NewRepository(db)
}

// providePeriodService creates a period service
func providePeriodService(repo *periodrepo.Repository) *periodservice.Service {
	return periodservice.NewService(repo)
}

// providePeriodHandler creates a period HTTP handler
func providePeriodHandler(service *periodservice.Service) *periodhttp.PeriodHandler {
	return periodhttp.NewPeriodHandler(service)
}

// provideProvisionRepository creates a provision repository
func provideProvisionRepository(db *sql.DB) *provisionrepo.Repository {
	return provisionrepo.NewRepository(db)
}

// provideProvisionService creates a provision service
// Note: journalService is nil for now - will be wired when journal module is integrated
func provideProvisionService(repo *provisionrepo.Repository) *provisionservice.Service {
	return provisionservice.NewService(repo, nil)
}

// provideProvisionHandler creates a provision HTTP handler
func provideProvisionHandler(service *provisionservice.Service) *provisionhttp.ProvisionHandler {
	return provisionhttp.NewProvisionHandler(service)
}
