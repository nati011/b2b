package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	auth "b2b.nati011.github.com/internal/core/application/auth"
	"b2b.nati011.github.com/internal/core/application/resource"
	"b2b.nati011.github.com/internal/core/application/role"
	"b2b.nati011.github.com/internal/core/application/user"
)

type Claims struct {
	ExpirationTime time.Time
	IssuedAt       time.Time
	NotBefore      time.Time
	Issuer         string
	Subject        string
	Audience       string
	Email          string
}

type Auth struct {
	authService     auth.Provider
	roleService     role.Provider
	ResourceService resource.Provider
	userService     user.Provider
}

type Option func(*Auth)

func NewAuthMiddleware(
	authService auth.Provider,
	roleService role.Provider,
	ResourceService resource.Provider,
	userService user.Provider,
) *Auth {
	return &Auth{
		authService:     authService,
		roleService:     roleService,
		ResourceService: ResourceService,
		userService:     userService,
	}
}

func WithRole(roles []string) Option {
	return func(a *Auth) {}
}

func (am *Auth) RequireAuthentication(next http.Handler, options ...Option) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			util.UnauthorizedResponse(w)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.UnauthorizedResponse(w)
			return
		}

		token := parts[1]
		if strings.TrimSpace(token) == "" {
			util.UnauthorizedResponse(w)
			return
		}

		result, err := am.authService.RetrospectToken(r.Context(), token)
		if err != nil {
			util.UnauthorizedResponse(w)
			return
		}

		if !result.Active {
			util.UnauthorizedResponse(w)
			return
		}

		decodedToken, err := am.authService.DecodeToken(r.Context(), token)
		if err != nil {
			util.UnauthorizedResponse(w)
			return
		}
		var claims = Claims(decodedToken.Claims)

		u, err := am.userService.GetByParam(r.Context(), &user.GetByParam{
			Email: claims.Email,
		})
		if err != nil {
			log.Printf("failed to get claims: %v", err)
			return
		}
		assignedRoles, err := am.userService.GetAllAssignedRoles(r.Context(), u.List[0].Id)
		if err != nil {
			log.Printf("failed to get assigned roles: %v", err)
			return
		}

		resource, err := am.ResourceService.GetByName(r.Context(), r.RequestURI)
		if err != nil {
			log.Printf("failed to get resource: %v", err)
			return
		}

		var hasResource bool
		for _, ro := range assignedRoles.List {
			hasResource, err = am.roleService.HasResource(r.Context(), &role.HasResourceRequest{
				ResourceId: resource.Id,
				RoleId:     ro.Id,
			})
			if err != nil {
				log.Printf("failed to get assigned roles: %v", err)
				return
			}
		}
		if hasResource {
			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			util.UnauthorizedResponse(w)
			return
		}
	})
}

func (am *Auth) RequireNoAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
