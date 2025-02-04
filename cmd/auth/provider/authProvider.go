package auth

type LoginAuthResonse struct {
}

type CreateClientAuthResonse struct {
}

type AuthProvider interface {
	CreateNewClient(firstName string, lastName string, email string, username string, password string) (CreateClientAuthResonse, error)
	ClientLogin(email, password string) (LoginAuthResonse, error)
}
