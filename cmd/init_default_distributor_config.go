package main

import (
	"context"
	"errors"
	"log"

	"b2b.nati011.github.com/config"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/role"
)

var (
	ErrDistributorSuperAdminRoleAlreadyCreated = errors.New(" distributor superadmin role already present")
)

var (
	DISTRIBUTOR_SUPERADMIN_ROLE_NAME = "distributor"
	DISTRIBUTOR_SUPERADMIN_ROLE_DESC = "distributor superadmin"
)

func checkIfDistributorSuperAdminRoleExists(applicationService *application_core.Container) (bool, error) {
	ctx := context.Background()
	_, err := applicationService.RoleService.Get(ctx, &role.GetRequest{
		Name: DISTRIBUTOR_SUPERADMIN_ROLE_NAME,
	})

	switch err {
	case nil:
		return true, nil
	case role.ErrEmptyGetContent:
		return false, nil
	default:
		return false, ErrUnknown
	}
}

func createDistributorSuperadminRole(applicationService *application_core.Container) (int, error) {
	ctx := context.Background()

	roleId, err := applicationService.RoleService.Create(ctx, &role.CreateRequest{
		Name: DISTRIBUTOR_SUPERADMIN_ROLE_NAME,
		Desc: DISTRIBUTOR_SUPERADMIN_ROLE_DESC,
	})
	if err != nil {
		log.Printf(" failed to create distributor superadmin role")
		return 0, ErrUnknown
	}

	resources, err := applicationService.ResourceService.GetAll(ctx)
	if err != nil {
		log.Printf(" failed to get all resources err: %v", err)
		return 0, ErrUnknown
	}

	for _, r := range resources.List {
		if err = applicationService.RoleService.AddResource(ctx,
			&role.AddResourceRequest{
				ResourceId: r.Id,
				RoleId:     roleId}); err != nil {
			log.Printf(" failed to add resourceId: %v to roleId: %v err: %v", r.Id, roleId, err)
			return 0, ErrUnknown
		}
	}
	return roleId, nil
}

func InitDistributorSuperadminRole(cfg config.Config, applicationService *application_core.Container) int {
	superadminRoleExists, err := checkIfDistributorSuperAdminRoleExists(applicationService)
	if err != nil {
		panic(" failed to check distributor superadmin role")
	}
	roleId := 0
	if superadminRoleExists {
		log.Print("# distributor superadmin role already created...")
	} else {
		log.Print("# creating distributor superadmin role...")
		roleId, err = createDistributorSuperadminRole(applicationService)
		if err != nil {
			panic("failed to add resources to distributor")
		}
	}
	return roleId
}
