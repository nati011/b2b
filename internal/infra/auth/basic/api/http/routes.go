package http

import (
	"net/http"
	"strings"

	httputil "marketplace/pkg/http"
	"marketplace/pkg/http/middleware"
)

// HTTP route paths
const (
	RouteAuthBasicCreds    = "/auth/basic/credentials"
	RouteAuthBasicCredsSuf = "/auth/basic/credentials/"
)

// Resource code
const ResourceAuthCreds = "auth_credentials"

// Path segments used in auth routes
const (
	PathSegmentAuth        = "auth"
	PathSegmentBasic       = "basic"
	PathSegmentCredentials = "credentials"
	PathSegmentPassword    = "password"
	PathSegmentDeactivate  = "deactivate"
)

// Action names used in auth context
const (
	ActionCreate     = "create"
	ActionUpdate     = "update"
	ActionDeactivate = "deactivate"
)

// HTTP route paths
const (
	RouteAuthLogin = "/api/v1/auth/login"
)

// publicRoutesProvider implements middleware.PublicRoutesProvider
type publicRoutesProvider struct{}

func (p *publicRoutesProvider) PublicRoutes() []string {
	return []string{
		"POST " + RouteAuthBasicCreds, // Allow users to create their own login credentials
		"POST " + RouteAuthLogin,      // Allow users to login
		"/auth/error",                 // Auth error page (e.g. after failed login redirect), no auth required
		"/auth/callback",              // OAuth callback, no auth required
	}
}

// init registers this module's public routes with the middleware registry
func init() {
	middleware.RegisterPublicRoutesProvider(&publicRoutesProvider{})
}

// @resource code=auth_credentials service=authentication desc="Basic authentication credentials"
// RegisterHTTPRoutes wires all basic auth credential HTTP routes into the provided mux.
func RegisterHTTPRoutes(mux *http.ServeMux, handler *Handler) {
	// @action name=login desc="Login with email and password"
	// POST /api/v1/auth/login - Login
	mux.HandleFunc(RouteAuthLogin, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.Login(w, r)
			return
		}
		httputil.MethodNotAllowed(w)
	})

	// @action name=create desc="Create basic authentication credentials"
	// POST /auth/basic/credentials - Create credential
	mux.HandleFunc(RouteAuthBasicCreds, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateCredential(w, r)
			return
		}
		httputil.MethodNotAllowed(w)
	})

	// @action name=update desc="Update credential password"
	// @action name=deactivate desc="Deactivate credential"
	// PUT /auth/basic/credentials/:username/password - Update password
	// POST /auth/basic/credentials/:username/deactivate - Deactivate credential
	mux.HandleFunc(RouteAuthBasicCredsSuf, func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(strings.TrimPrefix(r.URL.Path, RouteAuthBasicCreds), "/")
		segments := strings.Split(path, "/")

		if len(segments) >= 2 {
			action := segments[1]
			switch action {
			case PathSegmentPassword:
				if r.Method == http.MethodPut {
					handler.UpdatePassword(w, r)
					return
				}
			case PathSegmentDeactivate:
				if r.Method == http.MethodPost {
					handler.DeactivateCredential(w, r)
					return
				}
			}
		}

		httputil.MethodNotAllowed(w)
	})
}

