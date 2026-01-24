package service

import (
	"context"
	"marketplace/internal/infra/user/domain"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
)

// RegistrationTokenRepository defines the interface for registration token storage
type RegistrationTokenRepository interface {
	Create(ctx context.Context, token *domain.RegistrationToken) error
	FindByToken(ctx context.Context, token string) (*domain.RegistrationToken, error)
	MarkAsUsed(ctx context.Context, tokenID string) error
}

const (
	// RegistrationTokenExpiration is the duration before a registration token expires
	RegistrationTokenExpiration = 24 * time.Hour
	// RegistrationTokenLength is the length of the random token in bytes (before base64 encoding)
	RegistrationTokenLength = 32
)

// RegistrationTokenService handles business logic for registration tokens
type RegistrationTokenService struct {
	repo RegistrationTokenRepository
}

// NewRegistrationTokenService creates a new registration token service
func NewRegistrationTokenService(repo RegistrationTokenRepository) *RegistrationTokenService {
	return &RegistrationTokenService{
		repo: repo,
	}
}

// GenerateToken creates a new registration token for a user
func (s *RegistrationTokenService) GenerateToken(ctx context.Context, userID string) (*domain.RegistrationToken, error) {
	// Generate secure random token
	tokenBytes := make([]byte, RegistrationTokenLength)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, errors.New("failed to generate token")
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	now := time.Now()
	regToken := &domain.RegistrationToken{
		ID:        uuid.NewString(),
		UserID:    userID,
		Token:     token,
		Used:      false,
		ExpiresAt: now.Add(RegistrationTokenExpiration),
		CreatedAt: now,
	}

	if err := s.repo.Create(ctx, regToken); err != nil {
		return nil, err
	}

	return regToken, nil
}

// ValidateAndConsume validates a registration token and marks it as used
// Returns the user ID if valid, or an error if invalid/expired/used
func (s *RegistrationTokenService) ValidateAndConsume(ctx context.Context, token string) (string, error) {
	regToken, err := s.repo.FindByToken(ctx, token)
	if err != nil {
		// Assume any error means token not found
		return "", domain.ErrRegistrationTokenInvalid
	}

	if !regToken.IsValid() {
		if regToken.Used {
			return "", domain.ErrRegistrationTokenUsed
		}
		if time.Now().After(regToken.ExpiresAt) {
			return "", domain.ErrRegistrationTokenExpired
		}
		return "", domain.ErrRegistrationTokenInvalid
	}

	// Mark token as used
	if err := s.repo.MarkAsUsed(ctx, regToken.ID); err != nil {
		return "", err
	}

	return regToken.UserID, nil
}
