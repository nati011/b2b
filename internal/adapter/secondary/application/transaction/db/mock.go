package transaction

import (
	"context"
	"time"

	port "b2b.nati011.github.com/internal/port/application/transaction/db"
)

type GetResponse struct {
	Id         int
	User_Id    int
	Date       time.Time
	Amount     int32
	Partner_Id int
}

type GetAllResponse struct {
	List []GetAllResponse
}

type CreateRequest struct {
	User_Id    int
	Amount     int64
	Partner_Id int
}

type GetByParamRequest struct {
	Date       time.Time
	Partner_Id int
	User_Id    int
}

type MockTransaction struct {
	Id         int
	User_Id    int
	Date       time.Time
	Amount     int64
	Partner_Id int
}

type Mock struct {
	resources []MockTransaction
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	for _, i := range m.resources {
		if i.Id == id {
			return port.GetResponse{
				Id:         i.Id,
				User_Id:    i.User_Id,
				Date:       i.Date,
				Amount:     i.Amount,
				Partner_Id: i.Partner_Id,
			}, nil
		}
	}
	return port.GetResponse{}, port.ErrSysNoRows
}

func (m *Mock) GetAll(ctx context.Context, pagination *port.Pagination) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	count := 0
	for _, i := range m.resources {
		count++
		response = append(response, port.GetResponse{
			Id:         i.Id,
			User_Id:    i.User_Id,
			Date:       i.Date,
			Amount:     i.Amount,
			Partner_Id: i.Partner_Id,
		})
		if count > pagination.Limit && count == pagination.Limit {
			break
		}
	}
	if len(response) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: response,
	}, nil
}

func (m *Mock) GetByDate(ctx context.Context, date time.Time, pagination *port.Pagination) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	for _, i := range m.resources {
		if i.Date == date {
			response = append(response, port.GetResponse{
				Id:         i.Id,
				User_Id:    i.User_Id,
				Date:       i.Date,
				Amount:     i.Amount,
				Partner_Id: i.Partner_Id,
			})
		}
	}
	if len(response) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: response,
	}, nil
}

func (m *Mock) GetByUserId(ctx context.Context, userId int, pagination *port.Pagination) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	for _, i := range m.resources {
		if i.User_Id == userId {
			response = append(response, port.GetResponse{
				Id:         i.Id,
				User_Id:    i.User_Id,
				Date:       i.Date,
				Amount:     i.Amount,
				Partner_Id: i.Partner_Id,
			})
		}
	}
	if len(response) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: response,
	}, nil
}

func (m *Mock) GetByPartnerId(ctx context.Context, partnerId int, pagination *port.Pagination) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	for _, i := range m.resources {
		if i.Partner_Id == partnerId {
			response = append(response, port.GetResponse{
				Id:         i.Id,
				User_Id:    i.User_Id,
				Date:       i.Date,
				Amount:     i.Amount,
				Partner_Id: i.Partner_Id,
			})
		}
	}
	if len(response) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: response,
	}, nil
}

func (m *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	newResourceId := len(m.resources) + 1
	m.resources = append(m.resources, MockTransaction{
		Id:         newResourceId,
		User_Id:    req.User_Id,
		Amount:     req.Amount,
		Partner_Id: req.Partner_Id,
		Date:       time.Now(),
	})
	return newResourceId, nil
}

func (m *Mock) Delete(ctx context.Context, id int) error {
	updatedResources := []MockTransaction{}
	for _, i := range m.resources {
		if i.Id != id {
			updatedResources = append(updatedResources, MockTransaction{
				Id:         i.Id,
				User_Id:    i.User_Id,
				Date:       i.Date,
				Amount:     i.Amount,
				Partner_Id: i.Partner_Id,
			})
		}

	}
	m.resources = updatedResources
	return nil
}
