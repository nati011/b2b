package auth

import (
	provider "b2b.nati011.github.com/pkg/auth/provider"
)

type Container struct {
	AuthProvider provider.AuthProvider
	AuthService  AuthService
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
	c.AuthService = NewAuthService(&c.AuthProvider)
}
