package db

import (
	"context"
	"time"

	port "b2b.nati011.github.com/internal/port/user"
)

type MockUser struct {
	Id         int
	FullName   string
	Email      string
	Phone      string
	Username   string
	DOB        time.Time
	IsActive   bool
	ExternalId string
}

type MockRole struct {
	Id int
}

type MockResource struct {
	Id int
}

type Mock struct {
	users     []MockUser
	roles     []MockRole
	resources []MockResource
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.users {

		resp = append(resp, port.GetResponse{
			Id:         i.Id,
			FullName:   i.FullName,
			Email:      i.Email,
			Phone:      i.Phone,
			Username:   i.Username,
			DOB:        i.DOB,
			IsActive:   i.IsActive,
			ExternalId: i.ExternalId,
		})
	}
	if len(resp) == 0 {
		return port.GetResponse{}, port.ErrSysNoRows
	}

	return resp[0], nil
}

func (m *Mock) GetByEmail(ctx context.Context, email string) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.users {
		if i.Email == email {
			resp = append(resp, port.GetResponse{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		}
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) GetByPhone(ctx context.Context, phone string) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.users {
		if i.Phone == phone {
			resp = append(resp, port.GetResponse{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		}
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) GetByUsername(ctx context.Context, username string) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.users {
		if i.Username == username {
			resp = append(resp, port.GetResponse{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		}
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) GetByActiveStatus(ctx context.Context, status bool) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.users {
		if i.IsActive == status {
			resp = append(resp, port.GetResponse{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		}
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) GetByExternalId(ctx context.Context, extId string) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.users {
		if i.ExternalId == extId {
			resp = append(resp, port.GetResponse{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		}
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) GetAll(context.Context) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.users {

		resp = append(resp, port.GetResponse{
			Id:         i.Id,
			FullName:   i.FullName,
			Email:      i.Email,
			Phone:      i.Phone,
			Username:   i.Username,
			DOB:        i.DOB,
			IsActive:   i.IsActive,
			ExternalId: i.ExternalId,
		})
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	newUserId := len(m.users) + 1
	m.users = append(m.users, MockUser{
		Id:         newUserId,
		FullName:   req.FullName,
		Email:      req.Email,
		Phone:      req.Phone,
		Username:   req.Username,
		DOB:        req.DOB,
		IsActive:   req.IsActive,
		ExternalId: req.ExternalId,
	})
	return newUserId, nil
}

func (m *Mock) CreateAndActivate(ctx context.Context, req *port.CreateRequest) (int, error) {
	newUserId := len(m.users) + 1
	m.users = append(m.users, MockUser{
		Id:         newUserId,
		FullName:   req.FullName,
		Email:      req.Email,
		Phone:      req.Phone,
		Username:   req.Username,
		DOB:        req.DOB,
		IsActive:   true,
		ExternalId: req.ExternalId,
	})
	return newUserId, nil
}

func (m *Mock) Delete(ctx context.Context, id int) error {
	updatedResources := []MockUser{}
	for _, i := range m.users {
		if i.Id != id {
			updatedResources = append(updatedResources, MockUser{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		}
	}
	m.users = updatedResources
	return nil
}

func (m *Mock) UpdateFullName(ctx context.Context, req *port.UpdateFullNameRequest) (int, error) {
	updatedResourceId := req.Id
	updatedResources := []MockUser{}
	for _, i := range m.users {
		if i.Id == req.Id {
			updatedResources = append(updatedResources, MockUser{
				Id:         i.Id,
				FullName:   req.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		} else {
			updatedResources = append(updatedResources, MockUser{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		}
	}
	m.users = updatedResources
	return updatedResourceId, nil
}

func (m *Mock) UpdateEmail(ctx context.Context, req *port.UpdateEmailRequest) (int, error) {
	updatedResourceId := req.Id
	updatedResources := []MockUser{}
	for _, i := range m.users {
		if i.Id == req.Id {
			updatedResources = append(updatedResources, MockUser{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      req.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		} else {
			updatedResources = append(updatedResources, MockUser{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		}
	}
	m.users = updatedResources
	return updatedResourceId, nil
}

func (m *Mock) UpdatePhone(ctx context.Context, req *port.UpdatePhoneRequest) (int, error) {
	updatedResourceId := req.Id
	updatedResources := []MockUser{}
	for _, i := range m.users {
		if i.Id == req.Id {
			updatedResources = append(updatedResources, MockUser{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      req.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		} else {
			updatedResources = append(updatedResources, MockUser{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		}
	}
	m.users = updatedResources
	return updatedResourceId, nil
}

func (m *Mock) UpdateUsername(ctx context.Context, req *port.UpdateUsernameRequest) (int, error) {
	updatedResourceId := req.Id
	updatedResources := []MockUser{}
	for _, i := range m.users {
		if i.Id == req.Id {
			updatedResources = append(updatedResources, MockUser{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   req.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		} else {
			updatedResources = append(updatedResources, MockUser{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		}
	}
	m.users = updatedResources
	return updatedResourceId, nil
}

func (m *Mock) UpdateDOB(ctx context.Context, req *port.UpdateDOBRequest) (int, error) {
	updatedResourceId := req.Id
	updatedResources := []MockUser{}
	for _, i := range m.users {
		if i.Id == req.Id {
			updatedResources = append(updatedResources, MockUser{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        req.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		} else {
			updatedResources = append(updatedResources, MockUser{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		}
	}
	m.users = updatedResources
	return updatedResourceId, nil
}

func (m *Mock) UpdateIsActiveStatus(ctx context.Context, req *port.UpdateIsActiveRequest) (int, error) {
	updatedResourceId := req.Id
	updatedResources := []MockUser{}
	for _, i := range m.users {
		if i.Id == req.Id {
			updatedResources = append(updatedResources, MockUser{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   req.IsActive,
				ExternalId: i.ExternalId,
			})
		} else {
			updatedResources = append(updatedResources, MockUser{
				Id:         i.Id,
				FullName:   i.FullName,
				Email:      i.Email,
				Phone:      i.Phone,
				Username:   i.Username,
				DOB:        i.DOB,
				IsActive:   i.IsActive,
				ExternalId: i.ExternalId,
			})
		}
	}
	m.users = updatedResources
	return updatedResourceId, nil
}

func (m *Mock) AssignRole(ctx context.Context, id int, roleId int) error {
	m.roles = append(m.roles, MockRole{
		Id: roleId,
	})
	return nil
}

func (m *Mock) RemoveAssignedRole(ctx context.Context, id int, roleId int) error {
	for _, i := range m.roles {
		if i.Id != i.Id {
			m.roles = append(m.roles, MockRole{
				Id: i.Id,
			})
		}
	}
	return nil
}

func (m *Mock) GetAllAssignedRole(ctx context.Context, id int) (port.GetAllAssignedRoleResponse, error) {
	var resp port.GetAllAssignedRoleResponse
	for _, i := range m.roles {
		resp.List = append(resp.List, port.GetAssignedRoleResponse{
			Id: i.Id,
		})
	}
	if len(resp.List) == 0 {
		return port.GetAllAssignedRoleResponse{}, port.ErrSysNoRows
	}
	return resp, nil
}

func (m *Mock) HasAccessToResource(ctx context.Context, id int, roleId int) (bool, error) {
	for _, i := range m.resources {
		if i.Id == id {
			return true, nil
		}
	}
	return false, nil
}
