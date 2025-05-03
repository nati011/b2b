package configurable_product

import (
	"context"
	"math/rand"
	"time"

	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/configurable_product"
)

type MockConfigurableProduct struct {
	Id                int
	Name              string
	Desc              string
	ExternalId        string
	IsAvailableStatus bool
	Products          []int
	Images            []string
	AttributeKeys     []string
	CategoryId        []int
	DistributorId     int
}

type Mock struct {
	configurables []MockConfigurableProduct
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	id := rand.Intn(1000-10+1) + 10
	m.configurables = append(m.configurables, MockConfigurableProduct{
		Id:                id,
		Name:              req.Name,
		Desc:              req.Desc,
		ExternalId:        req.ExternalId,
		IsAvailableStatus: req.IsAvailableStatus,
		Products:          req.Products,
		Images:            req.Images,
		AttributeKeys:     req.AttributeKeys,
	})
	return id, nil
}

func (m *Mock) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	new_list := []MockConfigurableProduct{}
	for _, i := range m.configurables {
		if i.Id != req.Id {
			new_list = append(new_list, i)
		} else {
			new_list = append(new_list, MockConfigurableProduct{
				Id:                i.Id,
				Name:              req.Name,
				Desc:              i.Desc,
				ExternalId:        i.ExternalId,
				IsAvailableStatus: i.IsAvailableStatus,
				Products:          i.Products,
				Images:            i.Images,
				AttributeKeys:     i.AttributeKeys,
			})
		}
	}
	m.configurables = new_list
	return nil
}

func (m *Mock) UpdateDesc(ctx context.Context, req *port.UpdateDescRequest) error {
	new_list := []MockConfigurableProduct{}
	for _, i := range m.configurables {
		if i.Id != req.Id {
			new_list = append(new_list, i)
		} else {
			new_list = append(new_list, MockConfigurableProduct{
				Id:                i.Id,
				Name:              i.Name,
				Desc:              req.Desc,
				ExternalId:        i.ExternalId,
				IsAvailableStatus: i.IsAvailableStatus,
				Products:          i.Products,
				Images:            i.Images,
				AttributeKeys:     i.AttributeKeys,
			})
		}
	}
	m.configurables = new_list
	return nil
}

func (m *Mock) UpdateProducts(ctx context.Context, req *port.UpdateProductRequest) error {
	new_list := []MockConfigurableProduct{}
	for _, i := range m.configurables {
		if i.Id != req.Id {
			new_list = append(new_list, i)
		} else {
			new_list = append(new_list, MockConfigurableProduct{
				Id:                i.Id,
				Name:              i.Name,
				Desc:              i.Desc,
				ExternalId:        i.ExternalId,
				IsAvailableStatus: i.IsAvailableStatus,
				Products:          req.ProductIds,
				Images:            i.Images,
				AttributeKeys:     i.AttributeKeys,
			})
		}
	}
	m.configurables = new_list
	return nil
}

func (m *Mock) UpdateImages(ctx context.Context, req *port.UpdateImagesRequest) error {
	new_list := []MockConfigurableProduct{}
	for _, i := range m.configurables {
		if i.Id != req.Id {
			new_list = append(new_list, i)
		} else {
			new_list = append(new_list, MockConfigurableProduct{
				Id:                i.Id,
				Name:              i.Name,
				Desc:              i.Desc,
				ExternalId:        i.ExternalId,
				IsAvailableStatus: i.IsAvailableStatus,
				Products:          i.Products,
				Images:            req.Images,
				AttributeKeys:     i.AttributeKeys,
			})
		}
	}
	m.configurables = new_list
	return nil
}

func (m *Mock) UpdateExternalId(ctx context.Context, req *port.UpdateExternalIdRequest) error {
	new_list := []MockConfigurableProduct{}
	for _, i := range m.configurables {
		if i.Id != req.Id {
			new_list = append(new_list, i)
		} else {
			new_list = append(new_list, MockConfigurableProduct{
				Id:                i.Id,
				Name:              i.Name,
				Desc:              i.Desc,
				ExternalId:        req.ExternalId,
				IsAvailableStatus: i.IsAvailableStatus,
				Products:          i.Products,
				Images:            i.Images,
				AttributeKeys:     i.AttributeKeys,
			})
		}
	}
	m.configurables = new_list
	return nil
}

func (m *Mock) UpdateIsAvailableStatus(ctx context.Context, req *port.UpdateIsAvailableStatusRequest) error {
	new_list := []MockConfigurableProduct{}
	for _, i := range m.configurables {
		if i.Id != req.Id {
			new_list = append(new_list, i)
		} else {
			new_list = append(new_list, MockConfigurableProduct{
				Id:                i.Id,
				Name:              i.Name,
				Desc:              i.Desc,
				ExternalId:        i.ExternalId,
				IsAvailableStatus: req.Status,
				Products:          i.Products,
				Images:            i.Images,
				AttributeKeys:     i.AttributeKeys,
			})
		}
	}
	m.configurables = new_list
	return nil
}

func (m *Mock) UpdateAttributes(ctx context.Context, req *port.UpdateAttributes) error {
	new_list := []MockConfigurableProduct{}
	for _, i := range m.configurables {
		if i.Id != req.Id {
			new_list = append(new_list, i)
		} else {
			new_list = append(new_list, MockConfigurableProduct{
				Id:                i.Id,
				Name:              i.Name,
				Desc:              i.Desc,
				ExternalId:        i.ExternalId,
				IsAvailableStatus: i.IsAvailableStatus,
				AttributeKeys:     req.AttributeKeys,
				Products:          i.Products,
				Images:            i.Images,
			})
		}
	}
	m.configurables = new_list
	return nil
}

func (m *Mock) Disable(ctx context.Context, id int) error {
	new_list := []MockConfigurableProduct{}
	for _, i := range m.configurables {
		if i.Id != id {
			new_list = append(new_list, i)
		} else {
			new_list = append(new_list, MockConfigurableProduct{
				Id:                i.Id,
				Name:              i.Name,
				Desc:              i.Desc,
				ExternalId:        i.ExternalId,
				IsAvailableStatus: false,
				Products:          i.Products,
				Images:            i.Images,
				AttributeKeys:     i.AttributeKeys,
			})
		}
	}
	m.configurables = new_list
	return nil
}

func (m *Mock) Get(ctx context.Context, id int) (port.GetResponse, error) {
	for _, i := range m.configurables {

		if i.Id == id {
			return port.GetResponse{
				Id:            i.Id,
				Name:          i.Name,
				Desc:          i.Desc,
				ExternalId:    i.ExternalId,
				Attributes:    i.AttributeKeys,
				Products:      i.Products,
				IsAvailable:   i.IsAvailableStatus,
				CategoryId:    i.CategoryId,
				DistributorId: i.DistributorId,
				Images:        i.Images,
			}, nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.configurables {
		resp = append(resp, port.GetResponse{
			Id:         i.Id,
			Name:       i.Name,
			Desc:       i.Desc,
			ExternalId: i.ExternalId,
			// Attributes:    i.AttributeKeys,
			Products:      i.Products,
			IsAvailable:   i.IsAvailableStatus,
			CategoryId:    i.CategoryId,
			DistributorId: i.DistributorId,
			Images:        i.Images,
		})

	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) GetByName(ctx context.Context, name string) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.configurables {
		if i.Name == name {
			resp = append(resp, port.GetResponse{
				Id:         i.Id,
				Name:       i.Name,
				Desc:       i.Desc,
				ExternalId: i.ExternalId,
				// Attributes:    i.AttributeKeys,
				Products:      i.Products,
				IsAvailable:   i.IsAvailableStatus,
				CategoryId:    i.CategoryId,
				DistributorId: i.DistributorId,
				Images:        i.Images,
			})
		}
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) GetByExternalId(ctx context.Context, externalId string) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.configurables {
		if i.ExternalId == externalId {
			resp = append(resp, port.GetResponse{
				Id:         i.Id,
				Name:       i.Name,
				Desc:       i.Desc,
				ExternalId: i.ExternalId,
				// Attributes:    i.AttributeKeys,
				Products:      i.Products,
				IsAvailable:   i.IsAvailableStatus,
				CategoryId:    i.CategoryId,
				DistributorId: i.DistributorId,
				Images:        i.Images,
			})
		}
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: resp,
	}, nil
}
