package handler

import (
	"net/http"

	"b2b.nati011.github.com/internal/core"
)

var handlers []Handler

type Handler interface {
	Routes(mux *http.ServeMux)
	Init(services *core.MasterContainer) error
}

func Register(h Handler) {
	handlers = append(handlers, h)
}

func GetHandlers() []Handler {
	return handlers
}
