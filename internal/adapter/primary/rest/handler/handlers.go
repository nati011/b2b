package handler

import (
	"net/http"

	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	resource "b2b.nati011.github.com/internal/core/application/resource"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

var (
	ACTION_READ   = "READ"
	ACTION_WRITE  = "WRITE"
	ACTION_DELETE = "DELETE"
)

var handlers []Handler
var resources []resource.CreateRequest

type Handler interface {
	Routes(mux *http.ServeMux)
	Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainServices *domain_core.Container) error
}

func Register(h Handler) {
	handlers = append(handlers, h)
}

func RegisterResource(resource resource.CreateRequest) {
	resources = append(resources, resource)
}

func GetHandlers() []Handler {
	return handlers
}

func GetRoutes() []resource.CreateRequest {
	return resources
}
