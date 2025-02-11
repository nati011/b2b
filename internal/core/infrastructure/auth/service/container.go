package auth

import (
	provider "b2b.nati011.github.com/internal/core/infrastructure/auth/provider"
)

type Container struct {
	//exportable
	AuthService *AuthService

	//non exportables
	authProvider provider.AuthProvider
}

func NewContainer(ap provider.AuthProvider) *Container {
	c := new(Container)
	c.initAuthProvider(ap)
	c.initAuthService()
	return c
}

func (c *Container) initAuthProvider(ap provider.AuthProvider) {
	c.authProvider = ap
}

func (c *Container) initAuthService() {
	c.AuthService = NewAuthService(c.authProvider)
}
