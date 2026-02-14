package main

import (
	"errors"

	event "b2b.nati011.github.com/internal/adapter/primary/event/handler"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

var (
	ErrFailedToBuildEvent = errors.New(" failed to init event lister")
)

func InitEventLister(applicationServices *application_core.Container, domainServices *domain_core.Container) {
	err := event.BuildEventHandler(applicationServices, domainServices)
	if err != nil {
		panic("failed to init event listener")
	}
}
