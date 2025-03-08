package handler

import (
	"net/http"
)

var handlers []Handler

type Handler interface {
	Routes(mux *http.ServeMux)
	Init() error
}

func Register(h Handler) {
	handlers = append(handlers, h)
}

func GetHandlers() []Handler {
	return handlers
}
