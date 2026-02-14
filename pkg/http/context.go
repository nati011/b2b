package http

import (
	"context"

	userDomain "marketplace/internal/infra/user/domain"
)

type contextKey string

const (
	userContextKey        contextKey = "user"
	registrationUserIDKey contextKey = "registration_user_id"
)

// WithUser stores a user in the request context.
func WithUser(ctx context.Context, user *userDomain.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext retrieves the user from the request context.
// Returns nil if no user is present in the context.
func UserFromContext(ctx context.Context) *userDomain.User {
	user, ok := ctx.Value(userContextKey).(*userDomain.User)
	if !ok {
		return nil
	}
	return user
}

// WithRegistrationUserID stores a user ID from a validated registration token in the request context.
func WithRegistrationUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, registrationUserIDKey, userID)
}

// RegistrationUserIDFromContext retrieves the user ID from a validated registration token from the request context.
// Returns empty string if no registration user ID is present in the context.
func RegistrationUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(registrationUserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}
