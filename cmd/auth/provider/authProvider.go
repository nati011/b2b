package auth

type AuthProvider interface {
	CreateNewClient(firstName string, lastName string, email string, username string, password string) (string, error)
	ClientLogin(email, password string) (string, error)
}
