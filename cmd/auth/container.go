package auth

import (
	provider "b2b.nati011.github.com/cmd/auth/provider"
	authService "b2b.nati011.github.com/cmd/auth/service"
)

type Container struct {
	AuthProvider provider.AuthProvider
	AuthService  authService.AuthService
}

func NewContainer(ap provider.AuthProvider) *Container {
	c := new(Container)
	c.initAuthProvider(ap)
	c.initAuthService()
	return c
}

func (c *Container) initAuthProvider(ap provider.AuthProvider) {
	c.AuthProvider = ap
}

func (c *Container) initAuthService() {
	c.AuthService = authService.NewAuthService(&c.AuthProvider)
}
