package application

import (
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	mobileclient "b2b.nati011.github.com/internal/core/application/mobile_client"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type MobileClientHandler struct {
	service mobileclient.Provider
}

func InitMobileClient() {
	handler.Register(new(MobileClientHandler))

	handler.RegisterResource("/api/v1/minCompatibleClientVersion")
}

func (m *MobileClientHandler) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainService *domain_core.Container) error {
	m.service = applicationServices.MobileClient
	return nil
}

func (m *MobileClientHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/minCompatibleClientVersion", m.CheckminimumCompatibleVersion)
}

func (m *MobileClientHandler) CheckminimumCompatibleVersion(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(w, util.Envelope{"min_client_version": m.service.GetMinCompatibleVersion(r.Context())}, http.StatusAccepted)
}
