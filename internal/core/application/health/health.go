package health

import (
	"errors"
	"time"
)

var (
	ErrUnknown = errors.New(" unknown error")
)

type Provider interface {
	GetSystemHealth() (HealthReportResponse, error)
}

type HealthService struct {
}

func NewHealthService() Provider {
	return &HealthService{}
}

type HealthReportResponse struct {
	Status       string
	Timestamp    string
	Service      string
	Version      string
	Dependencies map[string]string
}

func (h *HealthService) GetSystemHealth() (HealthReportResponse, error) {
	// TODO: implement proper check health !!!
	return HealthReportResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Service:   "b2b-backend",
		Version:   "1.0.0",
		Dependencies: map[string]string{
			"postgres": "connected",
			"keycloak": "connected",
		},
	}, nil
}
