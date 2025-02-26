package email

import (
	"context"
	"errors"

	db "b2b.nati011.github.com/internal/core/application/util/email/provider/db"
)

var (
	ErrSysInvalidTemplate = errors.New("invalid template")
	ErrSysInvalidName     = errors.New("invalid name")
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
	List []GetResponse
}

type Templer interface {
	Create(context.Context, CreateRequest) (CreateResponse, error)
	GetAll(context.Context) (GetAllResponse, error)
	Get(context.Context, string) (GetResponse, error)
}

type TemplateService struct {
	db db.ReaderWriter
}

func NewTemplateService(db_provider db.ReaderWriter) Templer {
	return &TemplateService{
		db: db_provider,
	}
}

func (t *TemplateService) Create(ctx context.Context, req CreateRequest) (CreateResponse, error) {
	isValid_Name := ValidateName(req.Name)
	if !isValid_Name {
		return CreateResponse{
			Message: ErrSysInvalidName.Error(),
		}, ErrSysInvalidName
	}
	isValid_HTML := ValidateHTML(req.HtmlTemplate)
	if !isValid_HTML {
		return CreateResponse{
			Message: ErrSysInvalidTemplate.Error(),
		}, ErrSysInvalidTemplate
	}
	resp, err := t.db.Create(ctx, &db.CreateRequest{
		Name:         req.Name,
		HtmlTemplate: req.HtmlTemplate,
	})
	if err != nil {
		switch err {
		case db.ErrSysDuplicateName_L1:
			return CreateResponse{
				Name:    req.Name,
				Message: ErrSysDuplicateName.Error(),
			}, ErrSysDuplicateName
		}
	}
	return CreateResponse{
		Name:    resp.Name,
		Message: SUCCESS_MESSAGE,
	}, nil
}

func (t TemplateService) Get(ctx context.Context, r string) (GetResponse, error) {
	rslt, err := t.db.Get(ctx, r)
	if err != nil {
		switch err {
		case db.ErrSysUnknown_L1:
			return GetResponse{}, ErrSysUnknown
		}
	}
	return GetResponse{
		Name:         rslt.Name,
		HtmlTemplate: rslt.HtmlTemplate,
	}, nil
}

func (t TemplateService) GetAll(ctx context.Context) (GetAllResponse, error) {
	rslt, err := t.db.GetAll(ctx)
	if err != nil {
		switch err {
		case db.ErrSysUnknown_L1:
			return GetAllResponse{}, db.ErrSysUnknown_L1
		}
	}
	result := []GetResponse{}
	for _, i := range rslt.List {
		result = append(result, GetResponse{
			Name:         i.Name,
			HtmlTemplate: i.HtmlTemplate,
		})
	}

	return GetAllResponse{
		List: result,
	}, nil
}
