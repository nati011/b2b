package util

var (
	EVENT_DISTRIBUTOR_DEACTIVATE = "event_distributor_disable"
	EVENT_RETAILER_SSO           = "event_retailer_signon"
)

type EventRetailerSSOPayload struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
