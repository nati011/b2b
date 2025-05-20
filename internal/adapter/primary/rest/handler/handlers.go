package handler

import (
	"net/http"

	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

var handlers []Handler

type Handler interface {
	Routes(mux *http.ServeMux)
	Init(authMiddleWare *util.AuthMiddleware, applicationServices *application_core.Container, domainServices *domain_core.Container) error
}

func Register(h Handler) {
	handlers = append(handlers, h)
}

func GetHandlers() []Handler {
	return handlers
}
