package distributor

import "time"

type RegisterDistributorRequest struct {
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	DOB             time.Time `json:"date_of_birth"`
	PhoneNumber     string    `json:"phone_number"`
	ConfirmPassword string    `json:"confirmed_password"`
	FirstName       string    `json:"full_name"`
	Username        string    `json:"username"`
	ExternalId      string    `json:"external_id"`
}

type BusinessLocation struct {
	GeneralZone string `json:"general_zone"`
	Region      string `json:"region"`
	Woreda      string `json:"woreda"`
}

type UpdateBusinessRequest struct {
	DistributorId int              `json:"distributorId"`
	Name          string           `json:"name"`
	Tin           int              `json:"tin"`
	Region        BusinessLocation `json:"region"`
}

type RegisterDistributorResponse struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type GetResponse struct {
	Id         int       `json:"ID"`
	FirstName  string    `json:"full_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Username   string    `json:"username"`
	DOB        time.Time `json:"date_of_birth"`
	ExternalId string    `json:"external_id"`
}

type GetAllResponse struct {
	List []GetResponse
}

type GetByParamRequest struct {
	Name  string
	Email string
}
