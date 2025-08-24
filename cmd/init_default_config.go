package main

import (
	"log"

	"b2b.nati011.github.com/config"
	application_core "b2b.nati011.github.com/internal/core/application"
)

func InitDefaultConfig(cfg config.Config, applicationService *application_core.Container) {
	log.Print("# initializing default configs...")
	superadminRoleId := InitSuperadminRole(cfg, applicationService)
	InitSuperAdminUser(superadminRoleId, cfg, applicationService)
	InitDistributorSuperadminRole(cfg, applicationService)
}
