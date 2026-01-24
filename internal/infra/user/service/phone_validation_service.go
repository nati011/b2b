package service

import (
	"context"
	"regexp"
	"sync"
)

// PhoneValidationService handles phone number validation using database-stored regex patterns
type PhoneValidationService struct {
	repository PhoneValidationRepository
	cache      map[string]*regexp.Regexp
	cacheMutex sync.RWMutex
}

// NewPhoneValidationService creates a new phone validation service
func NewPhoneValidationService(repository PhoneValidationRepository) *PhoneValidationService {
	return &PhoneValidationService{
		repository: repository,
		cache:      make(map[string]*regexp.Regexp),
	}
}

// ValidatePhoneNumber validates a phone number using the default validation scheme
func (s *PhoneValidationService) ValidatePhoneNumber(ctx context.Context, phoneNumber string) error {
	scheme, err := s.repository.FindDefault(ctx)
	if err != nil {
		return err
	}

	return s.validateWithScheme(ctx, phoneNumber, scheme)
}

// ValidatePhoneNumberWithScheme validates a phone number using a specific validation scheme
func (s *PhoneValidationService) ValidatePhoneNumberWithScheme(ctx context.Context, phoneNumber string, schemeCode string) error {
	scheme, err := s.repository.FindByCode(ctx, schemeCode)
	if err != nil {
		return err
	}

	return s.validateWithScheme(ctx, phoneNumber, scheme)
}

// validateWithScheme validates a phone number against a specific scheme
func (s *PhoneValidationService) validateWithScheme(ctx context.Context, phoneNumber string, scheme *PhoneValidationScheme) error {
	// Check cache first
	s.cacheMutex.RLock()
	compiledRegex, exists := s.cache[scheme.RegexPattern]
	s.cacheMutex.RUnlock()

	if !exists {
		// Compile regex and cache it
		var err error
		compiledRegex, err = regexp.Compile(scheme.RegexPattern)
		if err != nil {
			return err
		}

		s.cacheMutex.Lock()
		s.cache[scheme.RegexPattern] = compiledRegex
		s.cacheMutex.Unlock()
	}

	if !compiledRegex.MatchString(phoneNumber) {
		return &PhoneValidationError{
			PhoneNumber: phoneNumber,
			Scheme:      scheme.Code,
			Pattern:     scheme.RegexPattern,
		}
	}

	return nil
}

// PhoneValidationError represents a phone number validation error
type PhoneValidationError struct {
	PhoneNumber string
	Scheme      string
	Pattern     string
}

func (e *PhoneValidationError) Error() string {
	return "phone number does not match validation pattern"
}
