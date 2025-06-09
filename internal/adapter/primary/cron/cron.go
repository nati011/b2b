package cron

import (
	domain_cron "b2b.nati011.github.com/internal/adapter/primary/cron/domain"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"github.com/go-co-op/gocron/v2"
)

func BuildCrons(scheduler gocron.Scheduler, applicationServices *application_core.Container, domainServices *domain_core.Container) error {
	domain_cron.InitOrderExpiry(scheduler, domainServices)
	return nil
}
