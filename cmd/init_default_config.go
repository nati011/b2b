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
	roleId := InitSuperadminRole(cfg, applicationService)
	InitSuperAdminUser(roleId, cfg, applicationService)
}

func InitSuperAdminUser(roleId int, cfg config.Config, applicationService *application_core.Container) {
	log.Print("# creating superadmin user...")
	ctx := context.Background()
	randomPassword, err := generateRandomPassword(10)
	if err != nil {
		panic("failed to create randomPassword for superadmin")
	}

	userId, err := applicationService.UserService.Create(ctx, &user.CreateRequest{
		FirstName: "superadmin",
		Email:     cfg.DefaultSuperAdminUserEmail,
		Password:  randomPassword,
	})
	if err != nil {
		panic(" failed to create superadmin user")
	}

	err = applicationService.UserService.AssignRole(ctx, userId, roleId)
	switch cfg.Env {
	case "development":
	case "staging", "production":
		err = applicationService.UserService.ResetPassword(ctx, cfg.DefaultSuperAdminUserEmail)
		if err != nil {
			panic("failed to perform reset password on superadmin")
		}
	default:
		panic("failed to identify env")
	}
}

func InitSuperadminRole(cfg config.Config, applicationService *application_core.Container) int {
	log.Print("# creating superadmin role...")
	ctx := context.Background()
	roleId, err := applicationService.RoleService.Create(ctx, &role.CreateRequest{
		Name: "superadmin",
		Desc: "superadmin",
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
