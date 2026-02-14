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
	ErrRetailerSuperAdminRoleAlreadyCreated = errors.New(" retailer superadmin role already present")
)

var (
	RETAILER_SUPERADMIN_ROLE_NAME = "retailer"
	RETAILER_SUPERADMIN_ROLE_DESC = "retailer superadmin"
)

func checkIfRetailerSuperAdminRoleExists(applicationService *application_core.Container) (bool, error) {
	ctx := context.Background()
	_, err := applicationService.RoleService.Get(ctx, &role.GetRequest{
		Name: RETAILER_SUPERADMIN_ROLE_NAME,
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

func createRetailerSuperadminRole(applicationService *application_core.Container) (int, error) {
	ctx := context.Background()

	roleId, err := applicationService.RoleService.Create(ctx, &role.CreateRequest{
		Name: RETAILER_SUPERADMIN_ROLE_NAME,
		Desc: RETAILER_SUPERADMIN_ROLE_DESC,
	})
	if err != nil {
		log.Printf(" failed to create retailer superadmin role")
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

func InitRetailerSuperadminRole(cfg config.Config, applicationService *application_core.Container) int {
	superadminRoleExists, err := checkIfRetailerSuperAdminRoleExists(applicationService)
	if err != nil {
		panic(" failed to check retailer superadmin role")
	}
	roleId := 0
	if superadminRoleExists {
		log.Print("# retailer superadmin role already created...")
	} else {
		log.Print("# creating retailer superadmin role...")
		roleId, err = createRetailerSuperadminRole(applicationService)
		if err != nil {
			panic("failed to add resources to retailer")
		}
	}
	return roleId
}
