package cors

import (
	"net/http"
	"strings"
)

// Middleware sets CORS headers and handles preflight. Use allowedOrigins ["*"] to allow all origins.
type Middleware struct {
	allowedOrigins []string
	allowAll       bool
}

func NewMiddleware(allowedOrigins []string) *Middleware {
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}
	allowAll := false
	for _, o := range allowedOrigins {
		if strings.TrimSpace(o) == "*" {
			allowAll = true
			break
		}
	}
	return &Middleware{
		allowedOrigins: allowedOrigins,
		allowAll:       allowAll,
	}
}

func (c *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))

		allowed := c.allowAll
		if !allowed {
			for _, o := range c.allowedOrigins {
				if strings.TrimSpace(o) == origin {
					allowed = true
					break
				}
			}
		}

		if allowed {
			allowOrigin := origin
			if allowOrigin == "" {
				allowOrigin = "*"
			}
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			if allowOrigin != "*" {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Accept-Language, Origin, X-User-ID, Cache-Control, Pragma")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
