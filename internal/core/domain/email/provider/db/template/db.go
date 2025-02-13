package email

import "errors"

var (
	ErrSysDuplicateName_L1 = errors.New("duplicate name")
	ErrSysUnknown_L1       = errors.New("unknown")
)

type CreateRequest struct {
	Name         string
	HtmlTemplate string
}

type CreateResponse struct {
	Name string
}

type GetResponse struct {
	Name         string
	HtmlTemplate string
}

type GetAllResponse struct {
	List []base
}
type base struct {
	Name         string
	HtmlTemplate string
}

type template struct {
	name string
	html string
}

type Provider interface {
	Create(*CreateRequest) (CreateResponse, error)
	Get(string) (GetResponse, error)
	GetAll() (GetAllResponse, error)
}
