package invoice

import "context"

type GetResponse struct {
}

type GetByExternalIdResponse struct {
}

type GetAllResponse struct {
	List []GetResponse
}

type CreateRequest struct {
}

type UpdateExternalIdRequest struct {
}

type UpdateStatusRequest struct {
}

type Reader interface {
	Get(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByExternalId(ctx context.Context, extId string) (GetAllResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
	UpdateExternalId(ctx context.Context, req *UpdateExternalIdRequest) error
	UpdateImages(ctx context.Context, req *UpdateStatusRequest) error
}

type DB interface {
	Reader
	Writer
}
