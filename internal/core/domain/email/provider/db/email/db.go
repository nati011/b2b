package email

type CreateRequest struct {
}

type CreateResponse struct {
	Id string
}

type GetRequest struct {
}

type GetResponse struct {
}

type GetAllResponse struct {
}

type Provider interface {
	Create(CreateRequest) (CreateResponse, error)
	Get(string) (GetResponse, error)
	GetAll() (GetAllResponse, error)
}
