package email

import (
	"errors"

	db "b2b.nati011.github.com/internal/core/domain/email/provider/db/template"
)

var (
	ErrSysInvalidTemplate = errors.New("template empty")
	ErrSysInvalidName     = errors.New("empty name")
	ErrSysUnknown         = errors.New("unknown")
	ErrSysDuplicateName   = errors.New("duplicate name")
)

const (
	SUCCESS_MESSAGE = "Ahoy, template created!"
)

type CreateRequest struct {
	Name         string
	HtmlTemplate string
}

type CreateResponse struct {
	Name    string
	Message string
}

type GetResponse struct {
	Name         string
	HtmlTemplate string
}

type GetAllResponse struct {
	list []base
}
type base struct {
	Name         string
	HtmlTemplate string
}

type Templer interface {
	Create(CreateRequest) (CreateResponse, error)
	GetAll() (GetAllResponse, error)
	Get(string) (GetResponse, error)
}

type TemplateService struct {
	db db.Provider
}

func NewTemplateService(db_provider db.Provider) Templer {
	return &TemplateService{
		db: db_provider,
	}
}

func (t *TemplateService) Create(req CreateRequest) (CreateResponse, error) {
	isValid_Name := ValidateName(req.Name)
	if !isValid_Name {
		return CreateResponse{}, ErrSysInvalidName
	}
	isValid_HTML := ValidateHTML(req.HtmlTemplate)
	if !isValid_HTML {
		return CreateResponse{}, ErrSysInvalidTemplate
	}
	resp, err := t.db.Create(&db.CreateRequest{})
	if err != nil {
		switch err {
		case db.ErrSysDuplicateName_L1:
			return CreateResponse{}, ErrSysDuplicateName
		}
	}
	return CreateResponse{
		Name:    resp.Name,
		Message: SUCCESS_MESSAGE,
	}, nil
}

func (t TemplateService) Get(r string) (GetResponse, error) {
	rslt, err := t.db.Get(r)
	if err != nil {
		switch err {
		case db.ErrSysUnknown_L1:
			return GetResponse{}, nil
		}
	}
	return GetResponse{
		Name:         rslt.Name,
		HtmlTemplate: rslt.HtmlTemplate,
	}, nil
}

func (t TemplateService) GetAll() (GetAllResponse, error) {
	rslt, err := t.db.GetAll()
	if err != nil {
		switch err {
		case db.ErrSysUnknown_L1:
			return GetAllResponse{}, db.ErrSysUnknown_L1
		}
	}
	result := []base{}
	for _, i := range rslt.List

	return GetAllResponse{}, nil
}
