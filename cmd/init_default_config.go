package main

import (
	"context"
	"errors"
	"log"

	"b2b.nati011.github.com/config"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/role"
	"b2b.nati011.github.com/internal/core/application/user"
)

var (
	ErrSuperAdminAlreadyCreated     = errors.New(" superadmin already present")
	ErrSuperAdminRoleAlreadyCreated = errors.New(" superadmin role already present")
	ErrUnknown                      = errors.New(" unknown error")
)

func InitDefaultConfig(cfg config.Config, applicationService *application_core.Container) {

	log.Print("# initializing default configs...")
	superadminRoleId := InitSuperadminRole(cfg, applicationService)
	InitSuperAdminUser(superadminRoleId, cfg, applicationService)
	// InitRetailerRole(cfg, applicationService)
	// InitDistributorSuperadminRole(cfg, applicationService)
}

func checkIfSuperAdminExists(cfg config.Config, applicationService *application_core.Container) (bool, error) {
	ctx := context.Background()
	_, err := applicationService.UserService.GetByParam(ctx, &user.GetByParam{
		Username: cfg.DefaultSuperAdminUserUsername,
	})

	switch err {
	case nil:
		return true, nil
	case user.ErrEmptyGetContent:
		return false, nil
	default:
		return false, ErrUnknown
	}
}

func createSuperAdminUser(roleId int, cfg config.Config, applicationService *application_core.Container) (int, error) {
	ctx := context.Background()
	log.Print("# creating superadmin user...")
	userId, err := applicationService.UserService.Create(ctx, &user.CreateRequest{
		FirstName: cfg.DefaultSuperAdminUserUsername,
		LastName:  cfg.DefaultSuperAdminUserUsername,
		Username:  cfg.DefaultSuperAdminUserUsername,
		Email:     cfg.DefaultSuperAdminUserEmail,
		Password:  cfg.DefaultSuperAdminUserPassword,
	})
	if err != nil {
		panic(" failed to create superadmin user")
	}

	err = applicationService.UserService.AssignRole(ctx, userId, roleId)
	if err != nil {

	}
	return userId, nil
}

func InitSuperAdminUser(roleId int, cfg config.Config, applicationService *application_core.Container) {
	ctx := context.Background()
	superAdminExists, err := checkIfSuperAdminExists(cfg, applicationService)
	if err != nil {
		panic(" failed to check if superadmin exists")
	}

	if superAdminExists {
		log.Print("# superadmin user already created...")
	} else {
		createSuperAdminUser(roleId, cfg, applicationService)
		switch cfg.Env {
		case "development":
		case "staging", "production":
			err = applicationService.UserService.ResetPassword(ctx, cfg.DefaultSuperAdminUserEmail)
			if err != nil {
				panic(" failed to perform reset password on superadmin")
			}
		default:
			panic(" failed to identify env")
		}
	}
}

func checkIfSuperAdminRoleExists(applicationService *application_core.Container) (bool, error) {
	ctx := context.Background()
	_, err := applicationService.RoleService.Get(ctx, &role.GetRequest{
		Name: SUPERADMIN_ROLE_NAME,
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

var (
	SUPERADMIN_ROLE_NAME = "superadmin"
	SUPERADMIN_ROLE_DESC = "benevolent dictator"
)

func createSuperadminRole(applicationService *application_core.Container) (int, error) {
	ctx := context.Background()
	roleId, err := applicationService.RoleService.Create(ctx, &role.CreateRequest{
		Name: SUPERADMIN_ROLE_NAME,
		Desc: SUPERADMIN_ROLE_DESC,
	})
	if err != nil {
		log.Printf(" failed to create superadmin role")
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

func InitSuperadminRole(cfg config.Config, applicationService *application_core.Container) int {
	superadminRoleExists, err := checkIfSuperAdminRoleExists(applicationService)
	if err != nil {
		panic(" failed to check superadmin role")
	}

	if superadminRoleExists {
		log.Print("# superadmin role already created...")
		return 0
	} else {
		log.Print("# creating superadmin role...")
		roleId, err := createSuperadminRole(applicationService)
		if err != nil {
			panic("failed to create superadmin role")
		}
		return roleId
	}
}

// func InitDistributorSuperadminRole(cfg config.Config, applicationService *application_core.Container) int {
// 	ctx := context.Background()
// 	_, err := applicationService.RoleService.Get(ctx, &role.GetRequest{
// 		Name: "distributor",
// 	})
// 	wantErr := role.ErrEmptyGetContent
// 	if err != wantErr {
// 		switch err {
// 		case nil:
// 			log.Print("# distributor role already created...")
// 			return 0
// 		default:
// 			panic("failed to get role")
// 		}
// 	}

// 	log.Print("# creating distributor superadmin role...")
// 	roleId, err := applicationService.RoleService.Create(ctx, &role.CreateRequest{
// 		Name: "distributor",
// 		Desc: "distributor superadmin",
// 	})
// 	if err != nil {
// 		panic(" failed to create distributor role")
// 	}

// 	resources, err := applicationService.ResourceService.GetAll(ctx)
// 	if err != nil {
// 		panic(" failed to fetch all resources")
// 	}

// 	for _, r := range resources.List {
// 		applicationService.RoleService.AddResource(ctx, &role.AddResourceRequest{
// 			ResourceId: r.Id,
// 			RoleId:     roleId,
// 		})
// 	}
// 	return roleId
// }

// func InitRetailerRole(cfg config.Config, applicationService *application_core.Container) int {
// 	ctx := context.Background()
// 	_, err := applicationService.RoleService.Get(ctx, &role.GetRequest{
// 		Name: "retailer",
// 	})
// 	wantErr := role.ErrEmptyGetContent
// 	if err != wantErr {
// 		switch err {
// 		case nil:
// 			log.Print("# retailer role already created...")
// 			return 0
// 		default:
// 			panic("HB5a158TiZ6r2tM8OonAgqIBpPSyFXF3failed to get role")
// 		}
// 	}

// 	log.Print("# creating retailer role...")
// 	roleId, err := applicationService.RoleService.Create(ctx, &role.CreateRequest{
// 		Name: "retailer",
// 		Desc: "retailer superadmin",
// 	})
// 	if err != nil {
// 		panic("failed to create retailer role")
// 	}

// 	resources, err := applicationService.ResourceService.GetAll(ctx)
// 	if err != nil {
// 		panic("failed to fetch all resources")
// 	}

// 	for _, r := range resources.List {
// 		err = applicationService.RoleService.AddResource(ctx, &role.AddResourceRequest{
// 			ResourceId: r.Id,
// 			RoleId:     roleId,
// 		})
// 		if err != nil {
// 			panic("failed to add resourceId")
// 		}
// 	}
// 	return roleId
// }

// func generateRandomPassword(length int) (string, error) {
// 	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
// 	password := make([]byte, length)

// 	for i := 0; i < length; i++ {
// 		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
// 		if err != nil {
// 			return "", err
// 		}
// 		password[i] = charset[randIndex.Int64()]
// 	}
// 	return string(password), nil
// }
