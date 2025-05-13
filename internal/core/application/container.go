package core

import (
	"database/sql"

	"b2b.nati011.github.com/config"
	auth_provider_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
	template_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/email-template/db"
	email_provider_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/email/smtp"
	payment_db_port "b2b.nati011.github.com/internal/adapter/secondary/application/payment/db"
	payment_partner_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/payment_partner/db"
	resource_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	role_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/role/db"
	transaction_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/transaction/db"
	user_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/user/db"

	// sms_provider_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/sms/provider"

	"b2b.nati011.github.com/internal/core/application/auth"
	"b2b.nati011.github.com/internal/core/application/checkout"
	"b2b.nati011.github.com/internal/core/application/email"
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
	"b2b.nati011.github.com/internal/core/domain/payment_verification"
	"b2b.nati011.github.com/internal/core/domain/retailer"

	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
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
	db                         *sql.DB
	AuthService                auth.Provider
	AuthMiddleware             *util.AuthMiddleware
	DistributorService         distributor.Provider
	RetailerService            retailer.Provider
	EmailService               email.Provider
	PaymentPartnerService      payment_partner.Provider
	RenderService              render.Provider
	ResourceService            resource.Provider
	RoleService                role.Provider
	SmsService                 sms.Provider
	TemplateService            template.Provider
	TransactionService         transaction.Provider
	UserService                user.Provider
	Pagination                 config.Pagination
	CheckoutService            checkout.Provider
	PaymentVerificationService payment_verification.Provider
	MobileClient               mobileclient.Provider
}

func NewContainer(
	//database
	db *sql.DB,
	//auth
	keycloakInstanceURL string,
	keycloakUsername string,
	keycloakPassword string,
	keycloakRealm string,
	keycloakApplicationRealm string,
	keycloakClientId string,
	keycloakClientSecret string,
	email_address,
	smtp_port string,
	//mobile client version
	MinMobileClientCompatibleVersion string,
	//baseurl
	baseUrl string,
	frontendUrl string) *Container {

	container := Container{}
	container.db = db

	//utils
	container.InitPagination()

	//ORDER ORDER!!
	container.InitAuthService(keycloakInstanceURL, keycloakUsername, keycloakPassword, keycloakRealm, keycloakApplicationRealm, keycloakClientId, keycloakClientSecret)
	container.InitUserService()
	container.InitTemplateService()
	container.InitRenderService()
	container.InitEmailService(email_address, smtp_port)
	container.InitPaymentPartnerService()
	container.InitTransactionService()
	container.InitResourceService()
	container.InitRoleService()
	container.InitUserService()
	container.InitMobileClientService(MinMobileClientCompatibleVersion)
	container.InitCheckoutService(baseUrl, frontendUrl)
	// container.InitSMSService()

	return &container
}

func (m *Container) InitMobileClientService(minMobileClientCompatibleVersion string) {
	m.MobileClient = mobileclient.NewMobileClientProvider(
		minMobileClientCompatibleVersion)
}

func (m *Container) InitAuthService(keycloakInstanceURL string, keycloakUsername string, keycloakPassword string, keycloakRealm string, keycloakApplicationRealm string, keycloakClientId string, keycloakClientSecret string) {
	m.AuthService = auth.NewAuthService(auth_provider_adapter.NewKeycloakProvider(keycloakInstanceURL, keycloakUsername, keycloakPassword, keycloakRealm, keycloakApplicationRealm, keycloakClientId, keycloakClientSecret))
	m.AuthMiddleware = util.NewAuthMiddleware(
		keycloakInstanceURL, keycloakClientId, keycloakClientSecret, keycloakRealm, keycloakPassword,
	)
}

func (m *Container) InitEmailService(email_address, smtp_port string) {
	m.EmailService = email.NewEmailService(email_provider_adapter.NewInbucket(email_address, smtp_port), m.RenderService)
}

func (m *Container) InitPaymentPartnerService() {
	m.PaymentPartnerService = payment_partner.NewPartner(payment_partner_db_adapter.NewPostgres(m.db, &m.Pagination))
}

func (m *Container) InitRenderService() {
	m.RenderService = render.NewRenderService(m.TemplateService)
}

func (m *Container) InitResourceService() {
	m.ResourceService = resource.NewResource(resource_db_adapter.NewPostgres(m.db, &m.Pagination))
}

func (m *Container) InitRoleService() {
	//create superAdminRole
	//grant Access to all resources
	m.RoleService = role.NewRole(role_db_adapter.NewPostgres(m.db, &m.Pagination), m.ResourceService)
}

// func (m *Container) InitSMSService() {
// 	m.SmsService = sms.NewSMSService()
// }

func (m *Container) InitTemplateService() {
	m.TemplateService = template.NewTemplateService(template_db_adapter.NewPostgres(m.db))
}

func (m *Container) InitTransactionService() {
	m.TransactionService = transaction.NewTransactionService(transaction_db_adapter.NewPostgres(m.db, &m.Pagination))
}

func (m *Container) InitUserService() {
	//create user with superadmin role
	m.UserService = user.NewUser(user_db_adapter.NewPostgres(m.db, &m.Pagination), m.RoleService, m.AuthService)
}

func (m *Container) InitPagination() {
	m.Pagination = *config.NewPaginationBuilder().Build()
}

func (m *Container) InitCheckoutService(baseUrl, frontendUrl string) {
	m.CheckoutService = checkout.NewCheckoutService(payment_db_port.NewPostgres(m.db), m.PaymentPartnerService, m.TransactionService, frontendUrl, baseUrl)
}
