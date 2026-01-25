package http

import (
	"net/http"

	httputil "marketplace/pkg/http"
	"marketplace/pkg/http/middleware"
)

// HTTP route paths
const (
	RouteReferrals        = "/referral"
	RouteReferralAffiliate = "/referral/affiliate"
	RouteReferralCode     = "/referral/code"
	RouteReferralValidate = "/referral/validate"
	RouteReferralRelationship = "/referral/relationship"
	RouteReferralCommission = "/referral/commission"
)

// Resource code
const ResourceReferrals = "referral"

// Action names used in referral context
const (
	ActionView   = "view"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

// publicRoutesProvider registers public referral routes.
type publicRoutesProvider struct{}

func (p *publicRoutesProvider) PublicRoutes() []string {
	return []string{
		"POST " + RouteReferralValidate,
	}
}

func init() {
	middleware.RegisterPublicRoutesProvider(&publicRoutesProvider{})
}

// @resource code=referral service=referral-management desc="Referral and affiliate management"
// RegisterHTTPRoutes wires all referral endpoints.
func RegisterHTTPRoutes(mux *http.ServeMux, handler *ReferralHandler) {
	// Affiliate endpoints
	// @action name=view desc="List affiliate records"
	// @action name=create desc="Create affiliate records"
	mux.HandleFunc(RouteReferralAffiliate, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListAffiliates(w, r)
		case http.MethodPost:
			handler.CreateAffiliate(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// @action name=view desc="Read a specific affiliate record"
	mux.HandleFunc(RouteReferralAffiliate+"/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetAffiliate(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// Referral code endpoints
	// @action name=create desc="Create referral code"
	mux.HandleFunc(RouteReferralCode, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateReferralCode(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// @action name=view desc="Read a specific referral code"
	mux.HandleFunc(RouteReferralCode+"/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetReferralCode(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// Referral validation endpoint
	// @action name=create desc="Validate a referral code"
	mux.HandleFunc(RouteReferralValidate, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.ValidateReferralCode(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// Referral relationship endpoints
	// @action name=create desc="Create referral relationship"
	mux.HandleFunc(RouteReferralRelationship, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateReferralRelationship(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// Commission endpoints
	// @action name=create desc="Create commission"
	mux.HandleFunc(RouteReferralCommission, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateCommission(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})
}

