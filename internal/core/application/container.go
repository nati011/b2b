package core

import (
	"database/sql"

	resource_db_port "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	"b2b.nati011.github.com/internal/adapter/secondary/application/sms"
	"b2b.nati011.github.com/internal/core/application/auth"
	"b2b.nati011.github.com/internal/core/application/email"
	"b2b.nati011.github.com/internal/core/application/payment"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/render"
	"b2b.nati011.github.com/internal/core/application/resource"
	"b2b.nati011.github.com/internal/core/application/role"
	"b2b.nati011.github.com/internal/core/application/template"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/application/user"
)

type Container struct {
	db                    *sql.DB
	AuthService           auth.Provider
	EmailService          email.Provider
	PaymentService        payment.Provider
	paymentPartnerService payment_partner.Provider
	RenderService         render.Renderer
	ResourceService       resource.Provider
	RoleService           role.Provider
	SmsService            sms.Provider
	templateService       template.Provider
	transactionService    transaction.Provider
	userService           user.Provider
}

func NewContainer(db *sql.DB) *Container {
	container := Container{}
	container.db = db

	container.InitResourceService()

	return &container
}

func (m *Container) InitAuthService() {
}

func (m *Container) InitEmailService() {
}

func (m *Container) InitPaymentService() {
}

func (m *Container) InitPaymentPartnerService() {
}

func (m *Container) InitRenderService() {
}

func (m *Container) InitResourceService() {
	m.ResourceService = resource.NewResource(resource_db_port.NewPostgres(m.db))
}

func (m *Container) InitRoleService() {
}

func (m *Container) InitSmsService() {
}

func (m *Container) InitTemplateService() {
}

func (m *Container) InitTransactionService() {
}

func (m *Container) InitUserService() {
}
