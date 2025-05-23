package handler

import (
	"net/http"

	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

var handlers []Handler

type Handler interface {
	Routes(mux *http.ServeMux)
	Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainServices *domain_core.Container) error
}

func Register(h Handler) {
	handlers = append(handlers, h)
}

func GetHandlers() []Handler {
	return handlers
}
