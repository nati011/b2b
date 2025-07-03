package distributor

import (
	"context"
)

type CreateRequest struct {
	Name        string
	Tin         string
	Latitude    string
	Longitude   string
	GeneralZone string
	Region      string
	Woreda      string
	UserId      int
}

type UpdateNameRequest struct {
	Name string
	Id   int
}

type UpdateTinRequest struct {
	Tin string
	Id  int
}

type GetResponse struct {
	Id          int
	Name        string
	Tin         string
	Latitude    string
	Longitude   string
	GeneralZone string
	Region      string
	Woreda      string
	IsActive    bool
	Verdict     string
}

type GetUserDetailResponse struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Username  string
}

type GetAllUserDetailResponse struct {
	List       []GetUserDetailResponse
	TotalCount int
}

type GetAllResponse struct {
	List       []GetResponse
	TotalCount int64
}

type CreateUserAgentRequest struct {
	User_id        int
	Distributor_Id int
}

type GetUserResponse struct {
	Id int
}

type GetAllUserResponse struct {
	List []GetUserResponse
}

type Reader interface {
	Get(ctx context.Context, id int) (GetResponse, error)
	GetByUserId(ctx context.Context, user_id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByName(ctx context.Context, name string) (GetAllResponse, error)
	GetByTin(ctx context.Context, tin string) (GetResponse, error)
	GetByStatus(ctx context.Context, status bool) (GetAllResponse, error)
	GetAllUserAgents(ctx context.Context, id int) (GetAllUserResponse, error)
	GetAllUserDetail(ctx context.Context, id int) (GetAllUserDetailResponse, error)
	GetByApprovalStatus(ctx context.Context, status string) (GetAllResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req CreateRequest) (int, error)
	CreateDistributorUser(ctx context.Context, req *CreateUserAgentRequest) (int, error)
	UpdateName(ctx context.Context, req *UpdateNameRequest) error
	UpdateTin(ctx context.Context, tin *UpdateTinRequest) error
	Activate(ctx context.Context, id int) error
	Dectivate(ctx context.Context, id int) error
}

type DB interface {
	Reader
	Writer
}
