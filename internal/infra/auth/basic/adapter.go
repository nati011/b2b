package basic

import (
	"context"

	basicauthservice "marketplace/internal/infra/auth/basic/service"
	"marketplace/pkg/http/middleware"
)

// Adapter adapts the basic auth service to the middleware BasicAuthService interface
type Adapter struct {
	service *basicauthservice.Service
}

// NewAdapter creates a new adapter for the basic auth service
func NewAdapter(service *basicauthservice.Service) *Adapter {
	return &Adapter{
		service: service,
	}
}

// ValidateCredentials validates credentials and returns a BasicCredential for the middleware
func (a *Adapter) ValidateCredentials(ctx context.Context, username, password string) (*middleware.BasicCredential, error) {
	cred, err := a.service.ValidateCredentials(ctx, username, password)
	if err != nil {
		return nil, err
	}

	return &middleware.BasicCredential{
		Username: cred.Username,
		UserID:   cred.UserID,
	}, nil
}
