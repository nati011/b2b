package product

import (
	"context"
	"time"

	"math/rand"

	port "b2b.nati011.github.com/internal/port/domain/product"
)

type MockProduct struct {
	Id            int
	Name          string
	Desc          string
	ExternalID    string
	Images        []string
	Price         float64
	Attributes    map[string]string
	DistributorId int
	CategoryId    []int
	Stock         int
	IsActive      bool
}

type Mock struct {
	products []MockProduct
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	id := rand.Intn(1000-10+1) + 10
	m.products = append(m.products, MockProduct{
		Id:            id,
		Name:          req.Name,
		Desc:          req.Desc,
		ExternalID:    req.ExternalID,
		Images:        req.Images,
		Price:         req.Price,
		Attributes:    req.Attributes,
		DistributorId: req.DistributorId,
		CategoryId:    req.CategoryId,
		Stock:         0,
		IsActive:      false,
	})
	return id, nil
}

func (m *Mock) Get(ctx context.Context, id int) (port.GetResponse, error) {
	for _, i := range m.products {
		if i.Id == id {
			return port.GetResponse{
				Id:            id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalID:    i.ExternalID,
				Images:        i.Images,
				Price:         i.Price,
				Attributes:    i.Attributes,
				DistributorId: i.DistributorId,
				CategoryId:    i.CategoryId,
				Stock:         i.Stock,
				IsActive:      i.IsActive,
			}, nil
		}
	}
	return port.GetResponse{}, port.ErrSysNoRows
}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	responses := port.GetAllResponse{}
	for _, i := range m.products {
		responses.List = append(responses.List,
			port.GetResponse{
				Id:            i.Id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalID:    i.ExternalID,
				Images:        i.Images,
				Price:         i.Price,
				Attributes:    i.Attributes,
				DistributorId: i.DistributorId,
				CategoryId:    i.CategoryId,
				Stock:         i.Stock,
				IsActive:      i.IsActive,
			},
		)
	}
	if len(responses.List) == 0 {
		return responses, port.ErrSysNoRows
	}
	return responses, nil
}

func (m *Mock) GetByName(ctx context.Context, req *port.GetByNameRequest) (port.GetAllResponse, error) {
	responses := port.GetAllResponse{}
	for _, i := range m.products {
		if req.Name == i.Name {
			responses.List = append(responses.List,
				port.GetResponse{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		}

	}
	if len(responses.List) == 0 {
		return responses, port.ErrSysNoRows
	}
	return responses, nil
}

func (m *Mock) GetByExternalId(ctx context.Context, req *port.GetByExternalIdRequest) (port.GetAllResponse, error) {
	responses := port.GetAllResponse{}
	for _, i := range m.products {
		if req.ExternalId == i.ExternalID {
			responses.List = append(responses.List,
				port.GetResponse{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		}

	}
	if len(responses.List) == 0 {
		return responses, port.ErrSysNoRows
	}
	return responses, nil
}

func (m *Mock) GetByDistributorId(ctx context.Context, req *port.GetByDistributorIdRequest) (port.GetAllResponse, error) {
	responses := port.GetAllResponse{}
	for _, i := range m.products {
		if req.DistributorId == i.DistributorId {
			responses.List = append(responses.List,
				port.GetResponse{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		}

	}
	if len(responses.List) == 0 {
		return responses, port.ErrSysNoRows
	}
	return responses, nil
}

func (m *Mock) GetByCategory(ctx context.Context, req *port.GetByCategoryRequest) (port.GetAllResponse, error) {
	responses := port.GetAllResponse{}
	for _, i := range m.products {
		for _, j := range i.CategoryId {
			for _, k := range req.CategoryId {
				if j == k {
					responses.List = append(responses.List,
						port.GetResponse{
							Id:            i.Id,
							Name:          i.Name,
							Desc:          i.Desc,
							ExternalID:    i.ExternalID,
							Images:        i.Images,
							Price:         i.Price,
							Attributes:    i.Attributes,
							DistributorId: i.DistributorId,
							CategoryId:    i.CategoryId,
							Stock:         i.Stock,
							IsActive:      i.IsActive,
						},
					)
				}
			}

		}

	}
	if len(responses.List) == 0 {
		return responses, port.ErrSysNoRows
	}
	return responses, nil
}

func (m *Mock) GetByPriceRange(ctx context.Context, req *port.GetByPriceRangeRequest) (port.GetAllResponse, error) {
	responses := port.GetAllResponse{}
	for _, i := range m.products {
		if req.PriceMin <= int(i.Price) && req.PriceMax >= int(i.Price) {
			responses.List = append(responses.List,
				port.GetResponse{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		}
	}

	if len(responses.List) == 0 {
		return responses, port.ErrSysNoRows
	}
	return responses, nil
}

func (m *Mock) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	responses := []MockProduct{}
	for _, i := range m.products {
		if req.Id == i.Id {
			responses = append(responses,
				MockProduct{
					Id:            i.Id,
					Name:          req.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		} else {
			responses = append(responses,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		}

	}
	m.products = responses
	return nil
}

func (m *Mock) UpdateExternalID(ctx context.Context, req *port.UpdateExternalIDRequest) error {
	responses := []MockProduct{}
	for _, i := range m.products {
		if req.Id == i.Id {
			responses = append(responses,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    req.ExternalId,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		} else {
			responses = append(responses,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		}

	}
	m.products = responses
	return nil
}

func (m *Mock) UpdatePrice(ctx context.Context, req *port.UpdatePriceRequest) error {
	responses := []MockProduct{}
	for _, i := range m.products {
		if req.Id == i.Id {
			responses = append(responses,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         float64(req.Price),
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		} else {
			responses = append(responses,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		}

	}
	m.products = responses
	return nil
}

func (m *Mock) UpdateDesc(ctx context.Context, req *port.UpdateDescRequest) error {
	responses := []MockProduct{}
	for _, i := range m.products {
		if req.Id == i.Id {
			responses = append(responses,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          req.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		} else {
			responses = append(responses,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		}

	}
	m.products = responses
	return nil
}

func (m *Mock) UpdateImages(ctx context.Context, req *port.UpdateImagesRequest) error {
	responses := []MockProduct{}
	for _, i := range m.products {
		if req.Id == i.Id {
			responses = append(responses,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        req.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		} else {
			responses = append(responses,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		}

	}
	m.products = responses
	return nil
}

func (m *Mock) UpdateCategoryId(ctx context.Context, req *port.UpdateCategoryIdRequest) error {
	responses := []MockProduct{}
	for _, i := range m.products {
		if req.Id == i.Id {
			responses = append(responses,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    req.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		} else {
			responses = append(responses,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		}

	}
	m.products = responses
	return nil
}

func (m *Mock) GoodsReceiving(ctx context.Context, req *port.GoodsReceivingRequest) error {
	resp := []MockProduct{}
	for _, i := range m.products {
		if req.Id == i.Id {
			resp = append(resp,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock + req.Amount,
					IsActive:      i.IsActive,
				},
			)
		}
	}
	m.products = resp
	return nil
}

func (m *Mock) Dispatch(ctx context.Context, req *port.DispatchRequest) error {
	for _, i := range m.products {
		if req.Id == i.Id {
			m.products = append(m.products,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock - req.Amount,
					IsActive:      i.IsActive,
				},
			)
		}
	}
	return nil
}

func (m *Mock) UpdateActiveStatus(ctx context.Context, req *port.UpdateActiveStatusRequest) error {
	responses := []MockProduct{}
	for _, i := range m.products {
		if req.Id == i.Id {
			responses = append(responses,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      req.Status,
				},
			)
		} else {
			responses = append(responses,
				MockProduct{
					Id:            i.Id,
					Name:          i.Name,
					Desc:          i.Desc,
					ExternalID:    i.ExternalID,
					Images:        i.Images,
					Price:         i.Price,
					Attributes:    i.Attributes,
					DistributorId: i.DistributorId,
					CategoryId:    i.CategoryId,
					Stock:         i.Stock,
					IsActive:      i.IsActive,
				},
			)
		}

	}
	m.products = responses
	return nil
}
