package core

import (
	"database/sql"

	auth_provider_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
	distributor_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/distributor/db"
	template_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/email-template/db"
	email_provider_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/email/smtp"
	payment_partner_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/payment-partner/db"
	resource_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	retailer_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/retailer/db"
	role_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/role/db"
	transaction_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/transaction/db"
	user_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/user/db"

	// sms_provider_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/sms/provider"

	"b2b.nati011.github.com/internal/core/application/auth"
	"b2b.nati011.github.com/internal/core/application/email"
	"b2b.nati011.github.com/internal/core/application/payment"
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
// ├── AuthService
// |
// ├── EmailService
// │   └── RenderService
// │       └── TemplateService
// |
// ├── PaymentService
// │   ├── PaymentPartnerService
// │   └── TransactionService
// |        └── UserService
// |
// ├──  UserService
// │   └── RoleService
// │       └── ResourceService
// |
// ├── SMSService

type Container struct {
	db                    *sql.DB
	AuthService           auth.Provider
	DistributorService    distributor.Provider
	RetailerService       retailer.Provider
	EmailService          email.Provider
	PaymentService        payment.Provider
	PaymentPartnerService payment_partner.Provider
	RenderService         render.Renderer
	ResourceService       resource.Provider
	RoleService           role.Provider
	SmsService            sms.Provider
	TemplateService       template.Provider
	TransactionService    transaction.Provider
	UserService           user.Provider
}

func NewContainer(db *sql.DB, keycloakInstanceURL string, keycloakUsername string, keycloakPassword string, keycloakRealm string, keycloakApplicationRealm string, keycloakClientId string, email_address, smtp_port string) *Container {
	container := Container{}
	container.db = db

	//ORDER ORDER!!
	container.InitAuthService(keycloakInstanceURL, keycloakUsername, keycloakPassword, keycloakRealm, keycloakApplicationRealm, keycloakClientId)
	container.InitDistributorService()
	container.InitResourceService()
	container.InitRoleService()
	container.InitTemplateService()
	container.InitRenderService()
	container.InitEmailService(email_address, smtp_port)
	container.InitPaymentPartnerService()
	container.InitTransactionService()
	container.InitPaymentPartnerService()
	container.InitPaymentService()
	container.InitResourceService()
	container.InitRetailerService()
	container.InitRoleService()
	container.InitUserService()
	// container.InitSMSService()

	return &container
}

func (m *Container) InitAuthService(keycloakInstanceURL string, keycloakUsername string, keycloakPassword string, keycloakRealm string, keycloakApplicationRealm string, keycloakClientId string) {
	m.AuthService = auth.NewAuthService(auth_provider_adapter.NewKeycloakProvider(keycloakInstanceURL, keycloakUsername, keycloakPassword, keycloakRealm, keycloakApplicationRealm, keycloakClientId))
}

func (m *Container) InitDistributorService() {
	m.DistributorService = distributor.NewDistributorService(distributor_db_adapter.NewPostgres(m.db), m.AuthService)
}

func (m *Container) InitEmailService(email_address, smtp_port string) {
	m.EmailService = email.NewEmailService(email_provider_adapter.NewInbucket(email_address, smtp_port), m.RenderService)
}

func (m *Container) InitPaymentService() {
	m.PaymentService = payment.NewPaymentService()
}

func (m *Container) InitPaymentPartnerService() {
	m.PaymentPartnerService = payment_partner.NewPartner(payment_partner_db_adapter.NewPostgres(m.db))
}

func (m *Container) InitRenderService() {
	m.RenderService = render.NewRenderService(m.TemplateService)
}

func (m *Container) InitResourceService() {
	m.ResourceService = resource.NewResource(resource_db_adapter.NewPostgres(m.db))
}

func (m *Container) InitRetailerService() {
	m.RetailerService = retailer.NewRetailerService(retailer_db_adapter.NewPostgres(m.db), m.AuthService)
}

func (m *Container) InitRoleService() {
	m.RoleService = role.NewRole(role_db_adapter.NewPostgres(m.db), m.ResourceService)
}

// func (m *Container) InitSMSService() {
// 	m.SmsService = sms.NewSMSService()
// }

func (m *Container) InitTemplateService() {
	m.TemplateService = template.NewTemplateService(template_db_adapter.NewPostgres(m.db))
}

func (m *Container) InitTransactionService() {
	m.TransactionService = transaction.NewTransactionService(transaction_db_adapter.NewPostgres(m.db), m.PaymentPartnerService, m.UserService)
}

func (m *Container) InitUserService() {
	m.UserService = user.NewUser(user_db_adapter.NewPostgres(m.db), m.RoleService)
}
