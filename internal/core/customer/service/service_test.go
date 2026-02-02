package customer

import (
	"context"
	"errors"
	"testing"

	"marketplace/internal/core/customer/domain"
	"marketplace/pkg/pagination"

	"github.com/stretchr/testify/require"
)

const (
	testEmailDup     = "dup@example.com"
	testEmailPrimary = "primary@example.com"
	testEmailOther   = "other@example.com"
)

type mockCustomerRepository struct {
	customers    map[int64]*domain.Customer
	byEmail      map[string]int64
	byPhone      map[string]int64
	nextID       int64
	createFunc   func(context.Context, *domain.Customer) error
	updateFunc   func(context.Context, *domain.Customer) error
	findByIDFunc func(context.Context, int64) (*domain.Customer, error)
}

func newMockCustomerRepository() *mockCustomerRepository {
	return &mockCustomerRepository{
		customers: make(map[int64]*domain.Customer),
		byEmail:   make(map[string]int64),
		byPhone:   make(map[string]int64),
		nextID:    1,
	}
}

func (m *mockCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, customer)
	}
	if customer.ID == 0 {
		customer.ID = m.nextID
		m.nextID++
	}
	m.customers[customer.ID] = customer
	if customer.Email != "" {
		m.byEmail[customer.Email] = customer.ID
	}
	if customer.PhoneNumber != "" {
		m.byPhone[customer.PhoneNumber] = customer.ID
	}
	return nil
}

func (m *mockCustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, customer)
	}
	if _, ok := m.customers[customer.ID]; !ok {
		return ErrCustomerNotFound
	}
	m.customers[customer.ID] = customer
	if customer.Email != "" {
		m.byEmail[customer.Email] = customer.ID
	}
	if customer.PhoneNumber != "" {
		m.byPhone[customer.PhoneNumber] = customer.ID
	}
	return nil
}

func (m *mockCustomerRepository) FindByID(ctx context.Context, id int64) (*domain.Customer, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	customer, ok := m.customers[id]
	if !ok {
		return nil, ErrCustomerNotFound
	}
	return customer, nil
}

func (m *mockCustomerRepository) FindByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	if id, ok := m.byEmail[email]; ok {
		return m.customers[id], nil
	}
	return nil, ErrCustomerNotFound
}

func (m *mockCustomerRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Customer, error) {
	if id, ok := m.byPhone[phoneNumber]; ok {
		return m.customers[id], nil
	}
	return nil, ErrCustomerNotFound
}

func (m *mockCustomerRepository) Delete(ctx context.Context, id int64) error {
	if _, ok := m.customers[id]; !ok {
		return ErrCustomerNotFound
	}
	delete(m.customers, id)
	return nil
}

func (m *mockCustomerRepository) FindAll(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Customer], error) {
	return pagination.NewPageResult([]*domain.Customer{}, 0, pageReq), nil
}

func (m *mockCustomerRepository) FindAllWithSupplierFilter(ctx context.Context, pageReq pagination.PageRequest, supplierID int64, search string) (pagination.PageResult[*domain.Customer], error) {
	// Mock implementation - return empty list for now
	return pagination.NewPageResult([]*domain.Customer{}, 0, pageReq), nil
}

func TestCustomerServiceCreateSuccess(t *testing.T) {
	repo := newMockCustomerRepository()
	service := NewCustomerService(repo)

	result, err := service.Create(context.Background(), CustomerInput{
		FullName:    "Ada Lovelace",
		Status:      "active",
		City:        "Addis Ababa",
		Region:      "Addis Ababa",
		Woreda:      "01",
		PhoneNumber: "0912345678",
		Email:       "ada@example.com",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotZero(t, result.ID)
	require.Equal(t, domain.CustomerStatusActive, result.Status)
	require.True(t, result.IsActive)
}

func TestCustomerServiceCreateEmailConflict(t *testing.T) {
	repo := newMockCustomerRepository()
	repo.Create(context.Background(), &domain.Customer{
		ID:       10,
		FullName: "Existing",
		Email:    testEmailDup,
	})
	repo.byEmail[testEmailDup] = 10

	service := NewCustomerService(repo)
	_, err := service.Create(context.Background(), CustomerInput{
		FullName: "New Person",
		Status:   "active",
		Email:    testEmailDup,
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrCustomerEmailExists))
}

func TestCustomerServiceUpdatePhoneConflict(t *testing.T) {
	repo := newMockCustomerRepository()
	repo.customers[1] = &domain.Customer{
		ID:          1,
		FullName:    "Primary",
		Status:      domain.CustomerStatusActive,
		PhoneNumber: "0911000000",
		Email:       testEmailPrimary,
	}
	repo.customers[2] = &domain.Customer{
		ID:          2,
		FullName:    "Other",
		Status:      domain.CustomerStatusActive,
		PhoneNumber: "0911222333",
		Email:       testEmailOther,
	}
	repo.byEmail[testEmailPrimary] = 1
	repo.byEmail[testEmailOther] = 2
	repo.byPhone["0911000000"] = 1
	repo.byPhone["0911222333"] = 2

	service := NewCustomerService(repo)
	_, err := service.Update(context.Background(), 1, CustomerInput{
		FullName:    "Primary",
		Status:      "active",
		PhoneNumber: "0911222333",
		Email:       testEmailPrimary,
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrCustomerPhoneExists))
}
