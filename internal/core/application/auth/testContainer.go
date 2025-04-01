package auth

import "b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"

func NewIntegrationAuthContainer() Provider {
	var mock = provider.NewMockAuthProvider()
	return NewAuthService(&mock)
}
