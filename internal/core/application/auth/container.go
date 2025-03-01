package auth

import (
	port "b2b.nati011.github.com/internal/port/application/auth/provider"
)

type Container struct {
	//exportable
	AuthService *AuthService

	//non exportables
	authProvider port.AuthProvider
}

func NewContainer(ap port.AuthProvider) *Container {
	c := new(Container)
	c.initAuthProvider(ap)
	c.initAuthService()
	return c
}

func (c *Container) initAuthProvider(ap port.AuthProvider) {
	c.authProvider = ap
}

func (c *Container) initAuthService() {
	c.AuthService = NewAuthService(c.authProvider)
}
