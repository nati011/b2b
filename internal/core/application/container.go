package core

import (
	"database/sql"

	"b2b.nati011.github.com/config"
	auth_provider_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
	template_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/email-template/db"
	email_provider_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/email/smtp"
	payment_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/payment/db"
	payment_partner_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/payment_partner/db"
	resource_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	role_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/role/db"
	transaction_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/transaction/db"
	user_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/user/db"

	payment "b2b.nati011.github.com/internal/core/application/payment"

	// sms_provider_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/sms/provider"

	auth "b2b.nati011.github.com/internal/core/application/auth"
	"b2b.nati011.github.com/internal/core/application/checkout"
	"b2b.nati011.github.com/internal/core/application/email"
	"b2b.nati011.github.com/internal/core/application/middleware"
	mobileclient "b2b.nati011.github.com/internal/core/application/mobile_client"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/render"
	"b2b.nati011.github.com/internal/core/application/resource"
	"b2b.nati011.github.com/internal/core/application/role"
	"b2b.nati011.github.com/internal/core/application/sms"
	"b2b.nati011.github.com/internal/core/application/template"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/application/user"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

/* Dependency Tree */

// Container
// ├── EmailService
// │   └── RenderService
// │       └── TemplateService
// |
// ├── PaymentService
// │   ├── PaymentPartnerService
// │   └── TransactionService
// |        └── UserService
// |
// ├── RetailerService
// |   └── UserService
// │    	└── RoleService
// │        	└── ResourceService
// │
// ├── DistributorService
// │	└── UserService
// │
// ├── UserService
// │	└── RoleService
// │
// ├── UserService
// │	└── AuthService
// │
// ├── RoleService
// │    └── ResourceService
// |
// ├── SMS-Service

type Container struct {
	db                    *sql.DB
	AuthService           auth.Provider
	AuthMiddleware        *middleware.Auth
	DistributorService    distributor.Provider
	RetailerService       retailer.Provider
	EmailService          email.Provider
	PaymentPartnerService payment_partner.Provider
	RenderService         render.Provider
	ResourceService       resource.Provider
	RoleService           role.Provider
	SmsService            sms.Provider
	TemplateService       template.Provider
	TransactionService    transaction.Provider
	UserService           user.Provider
	Pagination            *config.Pagination
	CheckoutService       checkout.Provider
	MobileClient          mobileclient.Provider
	PaymentService        payment.Provider
}

func NewContainer(
	db *sql.DB,
	cfg *config.Config,
	pagination *config.Pagination) *Container {

	container := Container{}
	container.db = db
	container.Pagination = pagination

	//ORDER ORDER!!
	container.InitTemplateService()
	container.InitRenderService()
	container.InitEmailService(cfg.Email, cfg.SMTP, cfg.EmailPassword)
	container.InitPaymentPartnerService()
	container.InitTransactionService()
	container.InitResourceService()
	container.InitRoleService()
	container.InitAuthService(
		cfg.KeycloakInstanceURL,
		cfg.KeycloakUsername,
		cfg.KeycloakPassword,
		cfg.KeycloakRealm,
		cfg.KeycloakApplicationRealm,
		cfg.KeycloakClientId,
		cfg.KeycloakClientSecret,
		cfg.JWTSecret,
	)
	container.InitUserService()
	container.InitPaymentService()
	container.InitCheckoutService(cfg.BaseUrl, cfg.FrontendUrl)
	container.InitAuthMiddleware()
	container.InitMobileClientService(cfg.MinMobileClientCompatibleVersion)
	// container.InitSMSService()

	return &container
}

func (m *Container) InitMobileClientService(minMobileClientCompatibleVersion string) {
	m.MobileClient = mobileclient.NewMobileClientProvider(
		minMobileClientCompatibleVersion)
}

func (m *Container) InitAuthService(keycloakInstanceURL string, keycloakUsername string, keycloakPassword string, keycloakRealm string, keycloakApplicationRealm string, keycloakClientId string, keycloakClientSecret string, jwtSecret string) {
	m.AuthService = auth.NewAuthService(
		auth_provider_adapter.NewKeycloakProvider(keycloakInstanceURL, keycloakUsername, keycloakPassword, keycloakRealm, keycloakApplicationRealm, keycloakClientId, keycloakClientSecret),
		m.EmailService,
		m.RoleService,
		jwtSecret,
	)
}

func (m *Container) InitAuthMiddleware() {
	m.AuthMiddleware = middleware.NewAuthMiddleware(
		m.AuthService,
		m.RoleService,
		m.ResourceService,
		m.UserService)
}

func (m *Container) InitEmailService(email_address, smtp_port, email_password string) {
	m.EmailService = email.NewEmailService(email_provider_adapter.NewGmail(email_address, smtp_port, email_password), m.RenderService)
}

func (m *Container) InitPaymentPartnerService() {
	m.PaymentPartnerService = payment_partner.NewPartner(payment_partner_db_adapter.NewPostgres(m.db, m.Pagination))
}

func (m *Container) InitRenderService() {
	m.RenderService = render.NewRenderService(m.TemplateService)
}

func (m *Container) InitResourceService() {
	m.ResourceService = resource.NewResource(resource_db_adapter.NewPostgres(m.db, m.Pagination))
}

func (m *Container) InitRoleService() {
	m.RoleService = role.NewRole(role_db_adapter.NewPostgres(m.db, m.Pagination), m.ResourceService)
}

// func (m *Container) InitSMSService() {
// 	m.SmsService = sms.NewSMSService()
// }

func (m *Container) InitTemplateService() {
	m.TemplateService = template.NewTemplateService(template_db_adapter.NewPostgres(m.db))
}

func (m *Container) InitTransactionService() {
	m.TransactionService = transaction.NewTransactionService(transaction_db_adapter.NewPostgres(m.db, m.Pagination))
}

func (m *Container) InitUserService() {
	m.UserService = user.NewUser(user_db_adapter.NewPostgres(m.db, m.Pagination), m.RoleService, m.AuthService)
}

func (m *Container) InitPaymentService() {
	m.PaymentService = payment.NewPaymentService(payment_db_adapter.NewPostgres(m.db))
}

func (m *Container) InitCheckoutService(baseUrl, frontendUrl string) {
	m.CheckoutService = checkout.NewCheckoutService(m.PaymentService, m.PaymentPartnerService, m.TransactionService, frontendUrl, baseUrl)
}
