package email

type Mock struct {
}

func NewMock() Provider {
	return &Mock{}
}
func (m Mock) Create(CreateRequest) (CreateResponse, error) {
	return CreateResponse{}, nil
}

func (m Mock) Get(string) (GetResponse, error) {
	return GetResponse{}, nil
}

func (m Mock) GetAll() (GetAllResponse, error) {
	return GetAllResponse{}, nil
}
