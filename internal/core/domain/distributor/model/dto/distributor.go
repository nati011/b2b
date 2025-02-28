package distributor

type RegisterDistributorRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmed_password"`
	FullName        string `json:"full_name"`
	Username        string `json:"username"`
}

type RegisterDistributorResponse struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}
