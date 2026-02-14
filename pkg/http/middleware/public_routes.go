package middleware

import "sync"

// PublicRoutesProvider is an interface for modules that provide public routes
type PublicRoutesProvider interface {
	PublicRoutes() []string
}

var (
	publicRoutesProviders []PublicRoutesProvider
	publicRoutesMutex     sync.RWMutex
)

// RegisterPublicRoutesProvider registers a module that provides public routes
func RegisterPublicRoutesProvider(provider PublicRoutesProvider) {
	publicRoutesMutex.Lock()
	defer publicRoutesMutex.Unlock()
	publicRoutesProviders = append(publicRoutesProviders, provider)
}

// CollectPublicRoutes collects all public routes from registered providers
func CollectPublicRoutes() []string {
	publicRoutesMutex.RLock()
	defer publicRoutesMutex.RUnlock()

	var routes []string
	for _, provider := range publicRoutesProviders {
		routes = append(routes, provider.PublicRoutes()...)
	}
	return routes
}

// AuthenticatedRoutesProvider is an interface for modules that provide authenticated routes
// (routes that require authentication but bypass permission checks)
type AuthenticatedRoutesProvider interface {
	AuthenticatedRoutes() []string
}

var (
	authenticatedRoutesProviders []AuthenticatedRoutesProvider
	authenticatedRoutesMutex     sync.RWMutex
)

// RegisterAuthenticatedRoutesProvider registers a module that provides authenticated routes
func RegisterAuthenticatedRoutesProvider(provider AuthenticatedRoutesProvider) {
	authenticatedRoutesMutex.Lock()
	defer authenticatedRoutesMutex.Unlock()
	authenticatedRoutesProviders = append(authenticatedRoutesProviders, provider)
}

// CollectAuthenticatedRoutes collects all authenticated routes from registered providers
func CollectAuthenticatedRoutes() []string {
	authenticatedRoutesMutex.RLock()
	defer authenticatedRoutesMutex.RUnlock()

	var routes []string
	for _, provider := range authenticatedRoutesProviders {
		routes = append(routes, provider.AuthenticatedRoutes()...)
	}
	return routes
}
