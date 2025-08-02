package main

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"

	"b2b.nati011.github.com/config"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/role"
	"b2b.nati011.github.com/internal/core/application/user"
)

func InitDefaultConfig(cfg config.Config, applicationService *application_core.Container) {

	log.Print("# initializing default configs...")
	superadminRoleId := InitSuperadminRole(cfg, applicationService)
	InitSuperAdminUser(superadminRoleId, cfg, applicationService)
	InitRetailerRole(cfg, applicationService)
	InitDistributorSuperadminRole(cfg, applicationService)
}

func InitSuperAdminUser(roleId int, cfg config.Config, applicationService *application_core.Container) {
	ctx := context.Background()
	randomPassword, err := generateRandomPassword(10)
	if err != nil {
		panic("failed to generate password for superadmin")
	}

	_, err = applicationService.UserService.GetByParam(ctx, &user.GetByParam{
		Username: "superadmin",
	})

	wantErr := user.ErrEmptyGetContent
	if err != wantErr {
		switch err {
		case nil:
			log.Print("# superadmin user already created...")
			return
		default:
			panic("failed to get user")
		}
	}
	var userId int
	log.Print("# creating superadmin user...")
	switch cfg.Env {
	case "development":
		userId, err = applicationService.UserService.Create(ctx, &user.CreateRequest{
			FirstName: "superadmin",
			LastName:  "superadmin",
			Username:  "superadmin",
			Email:     cfg.DefaultSuperAdminUserEmail,
			Password:  "superadmin",
		})
		if err != nil {
			panic(" failed to create superadmin user")
		}

		err = applicationService.UserService.AssignRole(ctx, userId, roleId)
		if err != nil {
			panic(" failed to assign role to superadmin")
		}
	case "staging", "production":
		userId, err = applicationService.UserService.Create(ctx, &user.CreateRequest{
			FirstName: "superadmin",
			LastName:  "superadmin",
			Username:  "superadmin",
			Email:     cfg.DefaultSuperAdminUserEmail,
			Password:  randomPassword,
		})
		if err != nil {
			switch err {
			case user.ErrEmailTaken:
				return
			default:
				panic(" failed to create superadmin user")
			}
		}

		err = applicationService.UserService.AssignRole(ctx, userId, roleId)
		if err != nil {
			panic(" failed to assign role to superadmin")
		}

		err = applicationService.UserService.ResetPassword(ctx, cfg.DefaultSuperAdminUserEmail)
		if err != nil {
			panic("failed to perform reset password on superadmin")
		}
	default:
		panic("failed to identify env")
	}
}

func InitSuperadminRole(cfg config.Config, applicationService *application_core.Container) int {
	ctx := context.Background()
	_, err := applicationService.RoleService.Get(ctx, &role.GetRequest{
		Name: "superadmin",
	})
	wantErr := role.ErrEmptyGetContent
	if err != wantErr {
		switch err {
		case nil:
			log.Print("# superadmin role already created...")
			return 0
		default:
			panic("failed to get role")
		}
	}

	log.Print("# creating superadmin role...")
	roleId, err := applicationService.RoleService.Create(ctx, &role.CreateRequest{
		Name: "superadmin",
		Desc: "benevolent_dictator",
	})
	if err != nil {
		panic("failed to create superadmin role")
	}

	resources, err := applicationService.ResourceService.GetAll(ctx)
	if err != nil {
		panic("failed to fetch all resources")
	}

	for _, r := range resources.List {
		applicationService.RoleService.AddResource(ctx, &role.AddResourceRequest{
			ResourceId: r.Id,
			RoleId:     roleId,
		})
	}
	return roleId
}

func InitDistributorSuperadminRole(cfg config.Config, applicationService *application_core.Container) int {
	ctx := context.Background()
	_, err := applicationService.RoleService.Get(ctx, &role.GetRequest{
		Name: "distributor",
	})
	wantErr := role.ErrEmptyGetContent
	if err != wantErr {
		switch err {
		case nil:
			log.Print("# distributor role already created...")
			return 0
		default:
			panic("failed to get role")
		}
	}

	log.Print("# creating distributor superadmin role...")
	roleId, err := applicationService.RoleService.Create(ctx, &role.CreateRequest{
		Name: "distributor",
		Desc: "distributor superadmin",
	})
	if err != nil {
		panic("failed to create distributor role")
	}

	resources, err := applicationService.ResourceService.GetAll(ctx)
	if err != nil {
		panic("failed to fetch all resources")
	}

	for _, r := range resources.List {
		applicationService.RoleService.AddResource(ctx, &role.AddResourceRequest{
			ResourceId: r.Id,
			RoleId:     roleId,
		})
	}
	return roleId
}

func InitRetailerRole(cfg config.Config, applicationService *application_core.Container) int {
	ctx := context.Background()
	_, err := applicationService.RoleService.Get(ctx, &role.GetRequest{
		Name: "retailer",
	})
	wantErr := role.ErrEmptyGetContent
	if err != wantErr {
		switch err {
		case nil:
			log.Print("# retailer role already created...")
			return 0
		default:
			panic("failed to get role")
		}
	}

	log.Print("# creating retailer role...")
	roleId, err := applicationService.RoleService.Create(ctx, &role.CreateRequest{
		Name: "retailer",
		Desc: "retailer superadmin",
	})
	if err != nil {
		panic("failed to create retailer role")
	}

	resources, err := applicationService.ResourceService.GetAll(ctx)
	if err != nil {
		panic("failed to fetch all resources")
	}

	for _, r := range resources.List {
		applicationService.RoleService.AddResource(ctx, &role.AddResourceRequest{
			ResourceId: r.Id,
			RoleId:     roleId,
		})
	}
	return roleId
}

func generateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[randIndex.Int64()]
	}
	return string(password), nil
}
