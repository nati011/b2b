package main

import (
	"errors"

	"b2b.nati011.github.com/internal/adapter/primary/cron"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"github.com/go-co-op/gocron/v2"
)

var (
	ErrFailedToBuildCrons = errors.New("¯\\_(ツ)_/¯, failed to init cron jobs")
)

func InitCron(s gocron.Scheduler, applicationServices *application_core.Container, domainServices *domain_core.Container) {
	err := cron.BuildCrons(s, applicationServices, domainServices)
	if err != nil {
		panic(ErrFailedToBuildCrons)
	}
}
