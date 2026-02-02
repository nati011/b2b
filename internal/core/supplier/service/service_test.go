package supplier

import (
	"context"
	"errors"
	"strings"
	"testing"

	"marketplace/internal/core/supplier/domain"
	"marketplace/pkg/pagination"

	"github.com/stretchr/testify/require"
)

type mockSupplierRepository struct {
	suppliers    map[int64]*domain.Supplier
	byEmail      map[string]int64
	byPhone      map[string]int64
	nextID       int64
	createFunc   func(context.Context, *domain.Supplier) error
	updateFunc   func(context.Context, *domain.Supplier) error
	findByIDFunc func(context.Context, int64) (*domain.Supplier, error)
}

type mockBankAccountRepository struct{}

func (m *mockBankAccountRepository) Create(ctx context.Context, account *domain.BankAccount) error {
	return nil
}

func (m *mockBankAccountRepository) Update(ctx context.Context, account *domain.BankAccount) error {
	return nil
}

func (m *mockBankAccountRepository) FindByID(ctx context.Context, id int64) (*domain.BankAccount, error) {
	return nil, ErrBankAccountNotFound
}

func (m *mockBankAccountRepository) FindBySupplierID(ctx context.Context, supplierID int64) ([]*domain.BankAccount, error) {
	return []*domain.BankAccount{}, nil
}

func (m *mockBankAccountRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func newMockSupplierRepository() *mockSupplierRepository {
	return &mockSupplierRepository{
		suppliers: make(map[int64]*domain.Supplier),
		byEmail:   make(map[string]int64),
		byPhone:   make(map[string]int64),
		nextID:    1,
	}
}

func (m *mockSupplierRepository) Create(ctx context.Context, supplier *domain.Supplier) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, supplier)
	}
	if supplier.ID == 0 {
		supplier.ID = m.nextID
		m.nextID++
	}
	m.suppliers[supplier.ID] = supplier
	if supplier.SupportEmail != "" {
		m.byEmail[supplier.SupportEmail] = supplier.ID
	}
	if supplier.SupportPhone != "" {
		m.byPhone[supplier.SupportPhone] = supplier.ID
	}
	return nil
}

func (m *mockSupplierRepository) Update(ctx context.Context, supplier *domain.Supplier) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, supplier)
	}
	if _, ok := m.suppliers[supplier.ID]; !ok {
		return ErrSupplierNotFound
	}
	m.suppliers[supplier.ID] = supplier
	if supplier.SupportEmail != "" {
		m.byEmail[supplier.SupportEmail] = supplier.ID
	}
	if supplier.SupportPhone != "" {
		m.byPhone[supplier.SupportPhone] = supplier.ID
	}
	return nil
}

func (m *mockSupplierRepository) FindByID(ctx context.Context, id int64) (*domain.Supplier, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	supplier, ok := m.suppliers[id]
	if !ok {
		return nil, ErrSupplierNotFound
	}
	return supplier, nil
}

func (m *mockSupplierRepository) FindByEmail(ctx context.Context, email string) (*domain.Supplier, error) {
	if id, ok := m.byEmail[email]; ok {
		return m.suppliers[id], nil
	}
	return nil, ErrSupplierNotFound
}

func (m *mockSupplierRepository) FindByPhone(ctx context.Context, phone string) (*domain.Supplier, error) {
	if id, ok := m.byPhone[phone]; ok {
		return m.suppliers[id], nil
	}
	return nil, ErrSupplierNotFound
}

func (m *mockSupplierRepository) Delete(ctx context.Context, id int64) error {
	if _, ok := m.suppliers[id]; !ok {
		return ErrSupplierNotFound
	}
	delete(m.suppliers, id)
	return nil
}

func (m *mockSupplierRepository) FindAll(ctx context.Context, pageReq pagination.PageRequest, search string) (pagination.PageResult[*domain.Supplier], error) {
	items := make([]*domain.Supplier, 0, len(m.suppliers))
	for _, supplier := range m.suppliers {
		// Apply search filter if provided
		if search != "" {
			searchLower := strings.ToLower(search)
			matches := strings.Contains(strings.ToLower(supplier.BusinessName), searchLower) ||
				strings.Contains(strings.ToLower(supplier.SupportEmail), searchLower) ||
				strings.Contains(strings.ToLower(supplier.SupportPhone), searchLower)
			if !matches {
				continue
			}
		}
		items = append(items, supplier)
	}
	return pagination.NewPageResult(items, len(items), pageReq), nil
}

func TestSupplierServiceCreateSuccess(t *testing.T) {
	repo := newMockSupplierRepository()
	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)

	result, err := service.Create(context.Background(), SupplierInput{
		BusinessName: "Acme Corporation",
		Status:       "active",
		SupportEmail: "support@acme.com",
		SupportPhone: "0912345678",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotZero(t, result.ID)
	require.Equal(t, domain.SupplierStatusActive, result.Status)
	require.True(t, result.IsActive)
	require.Equal(t, "Acme Corporation", result.BusinessName)
	require.Equal(t, "support@acme.com", result.SupportEmail)
	require.Equal(t, "0912345678", result.SupportPhone)
}

func TestSupplierServiceCreateWithEmailOnly(t *testing.T) {
	repo := newMockSupplierRepository()
	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)

	result, err := service.Create(context.Background(), SupplierInput{
		BusinessName: "Tech Solutions Inc",
		Status:       "active",
		SupportEmail: "info@techsolutions.com",
		SupportPhone: "",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "info@techsolutions.com", result.SupportEmail)
	require.Empty(t, result.SupportPhone)
}

func TestSupplierServiceCreateWithPhoneOnly(t *testing.T) {
	repo := newMockSupplierRepository()
	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)

	result, err := service.Create(context.Background(), SupplierInput{
		BusinessName: "Phone Only Corp",
		Status:       "inactive",
		SupportEmail: "",
		SupportPhone: "0911111111",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Empty(t, result.SupportEmail)
	require.Equal(t, "0911111111", result.SupportPhone)
	require.Equal(t, domain.SupplierStatusInactive, result.Status)
	require.False(t, result.IsActive)
}

func TestSupplierServiceCreateEmailConflict(t *testing.T) {
	repo := newMockSupplierRepository()
	existingSupplier, _ := domain.NewSupplier("Existing Corp", "active", "dup@example.com", "0910000000")
	existingSupplier.ID = 10
	repo.Create(context.Background(), existingSupplier)
	repo.byEmail["dup@example.com"] = 10

	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)
	_, err := service.Create(context.Background(), SupplierInput{
		BusinessName: "New Corporation",
		Status:       "active",
		SupportEmail: "dup@example.com",
		SupportPhone: "0922222222",
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrSupplierEmailExists))
}

func TestSupplierServiceCreatePhoneConflict(t *testing.T) {
	repo := newMockSupplierRepository()
	existingSupplier, _ := domain.NewSupplier("Existing Corp", "active", "existing@example.com", "0911111111")
	existingSupplier.ID = 10
	repo.Create(context.Background(), existingSupplier)
	repo.byPhone["0911111111"] = 10

	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)
	_, err := service.Create(context.Background(), SupplierInput{
		BusinessName: "New Corporation",
		Status:       "active",
		SupportEmail: "new@example.com",
		SupportPhone: "0911111111",
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrSupplierPhoneExists))
}

func TestSupplierServiceGetSuccess(t *testing.T) {
	repo := newMockSupplierRepository()
	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)

	supplier, _ := domain.NewSupplier("Test Corp", "active", "test@example.com", "0912345678")
	supplier.ID = 1
	repo.suppliers[1] = supplier

	result, err := service.Get(context.Background(), 1)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, int64(1), result.ID)
	require.Equal(t, "Test Corp", result.BusinessName)
	require.Equal(t, domain.SupplierStatusActive, result.Status)
}

func TestSupplierServiceGetNotFound(t *testing.T) {
	repo := newMockSupplierRepository()
	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)

	_, err := service.Get(context.Background(), 999)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrSupplierNotFound))
}

func TestSupplierServiceUpdateSuccess(t *testing.T) {
	repo := newMockSupplierRepository()
	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)

	existingSupplier, _ := domain.NewSupplier("Old Corp", "active", "old@example.com", "0910000000")
	existingSupplier.ID = 1
	repo.suppliers[1] = existingSupplier
	repo.byEmail["old@example.com"] = 1
	repo.byPhone["0910000000"] = 1

	result, err := service.Update(context.Background(), 1, SupplierInput{
		BusinessName: "Updated Corp",
		Status:       "active",
		SupportEmail: "new@example.com",
		SupportPhone: "0922222222",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "Updated Corp", result.BusinessName)
	require.Equal(t, "new@example.com", result.SupportEmail)
	require.Equal(t, "0922222222", result.SupportPhone)
}

func TestSupplierServiceUpdateNotFound(t *testing.T) {
	repo := newMockSupplierRepository()
	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)

	_, err := service.Update(context.Background(), 999, SupplierInput{
		BusinessName: "New Corp",
		Status:       "active",
		SupportEmail: "new@example.com",
		SupportPhone: "0912345678",
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrSupplierNotFound))
}

func TestSupplierServiceUpdateEmailConflict(t *testing.T) {
	repo := newMockSupplierRepository()
	repo.suppliers[1] = &domain.Supplier{
		ID:           1,
		BusinessName: "Primary",
		Status:       domain.SupplierStatusActive,
		SupportPhone: "0911000000",
		SupportEmail: "primary@example.com",
	}
	repo.suppliers[2] = &domain.Supplier{
		ID:           2,
		BusinessName: "Other",
		Status:       domain.SupplierStatusActive,
		SupportPhone: "0911222333",
		SupportEmail: "other@example.com",
	}
	repo.byEmail["primary@example.com"] = 1
	repo.byEmail["other@example.com"] = 2
	repo.byPhone["0911000000"] = 1
	repo.byPhone["0911222333"] = 2

	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)
	_, err := service.Update(context.Background(), 1, SupplierInput{
		BusinessName: "Primary",
		Status:       "active",
		SupportPhone: "0911000000",
		SupportEmail: "other@example.com",
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrSupplierEmailExists))
}

func TestSupplierServiceUpdatePhoneConflict(t *testing.T) {
	repo := newMockSupplierRepository()
	repo.suppliers[1] = &domain.Supplier{
		ID:           1,
		BusinessName: "Primary",
		Status:       domain.SupplierStatusActive,
		SupportPhone: "0911000000",
		SupportEmail: "primary@example.com",
	}
	repo.suppliers[2] = &domain.Supplier{
		ID:           2,
		BusinessName: "Other",
		Status:       domain.SupplierStatusActive,
		SupportPhone: "0911222333",
		SupportEmail: "other@example.com",
	}
	repo.byEmail["primary@example.com"] = 1
	repo.byEmail["other@example.com"] = 2
	repo.byPhone["0911000000"] = 1
	repo.byPhone["0911222333"] = 2

	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)
	_, err := service.Update(context.Background(), 1, SupplierInput{
		BusinessName: "Primary",
		Status:       "active",
		SupportPhone: "0911222333",
		SupportEmail: "primary@example.com",
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrSupplierPhoneExists))
}

func TestSupplierServiceUpdateSameEmailAndPhone(t *testing.T) {
	repo := newMockSupplierRepository()
	existingSupplier, _ := domain.NewSupplier("Test Corp", "active", "test@example.com", "0912345678")
	existingSupplier.ID = 1
	repo.suppliers[1] = existingSupplier
	repo.byEmail["test@example.com"] = 1
	repo.byPhone["0912345678"] = 1

	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)
	result, err := service.Update(context.Background(), 1, SupplierInput{
		BusinessName: "Updated Corp",
		Status:       "active",
		SupportEmail: "test@example.com",
		SupportPhone: "0912345678",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "Updated Corp", result.BusinessName)
	require.Equal(t, "test@example.com", result.SupportEmail)
	require.Equal(t, "0912345678", result.SupportPhone)
}

func TestSupplierServiceDeleteSuccess(t *testing.T) {
	repo := newMockSupplierRepository()
	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)

	supplier, _ := domain.NewSupplier("To Delete", "active", "delete@example.com", "0912345678")
	supplier.ID = 1
	repo.suppliers[1] = supplier

	err := service.Delete(context.Background(), 1)
	require.NoError(t, err)
	_, exists := repo.suppliers[1]
	require.False(t, exists)
}

func TestSupplierServiceDeleteNotFound(t *testing.T) {
	repo := newMockSupplierRepository()
	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)

	err := service.Delete(context.Background(), 999)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrSupplierNotFound))
}

func TestSupplierServiceListSuccess(t *testing.T) {
	repo := newMockSupplierRepository()
	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)

	supplier1, _ := domain.NewSupplier("Corp 1", "active", "corp1@example.com", "0911111111")
	supplier1.ID = 1
	repo.suppliers[1] = supplier1

	supplier2, _ := domain.NewSupplier("Corp 2", "inactive", "corp2@example.com", "0922222222")
	supplier2.ID = 2
	repo.suppliers[2] = supplier2

	result, err := service.List(context.Background(), pagination.PageRequest{Page: 1, Limit: 10}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 2, result.Total)
	require.Len(t, result.Items, 2)
}

func TestSupplierServiceCreateInvalidInput(t *testing.T) {
	repo := newMockSupplierRepository()
	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)

	// Missing business name
	_, err := service.Create(context.Background(), SupplierInput{
		BusinessName: "",
		Status:       "active",
		SupportEmail: "test@example.com",
		SupportPhone: "0912345678",
	})
	require.Error(t, err)

	// Missing both email and phone
	_, err = service.Create(context.Background(), SupplierInput{
		BusinessName: "Test Corp",
		Status:       "active",
		SupportEmail: "",
		SupportPhone: "",
	})
	require.Error(t, err)

	// Invalid email format
	_, err = service.Create(context.Background(), SupplierInput{
		BusinessName: "Test Corp",
		Status:       "active",
		SupportEmail: "invalid-email",
		SupportPhone: "",
	})
	require.Error(t, err)
}

func TestSupplierServiceUpdateStatusChange(t *testing.T) {
	repo := newMockSupplierRepository()
	bankAccountRepo := &mockBankAccountRepository{}
	service := NewSupplierService(repo, bankAccountRepo)

	existingSupplier, _ := domain.NewSupplier("Test Corp", "inactive", "test@example.com", "0912345678")
	existingSupplier.ID = 1
	existingSupplier.IsActive = false
	repo.suppliers[1] = existingSupplier
	repo.byEmail["test@example.com"] = 1
	repo.byPhone["0912345678"] = 1

	result, err := service.Update(context.Background(), 1, SupplierInput{
		BusinessName: "Test Corp",
		Status:       "active",
		SupportEmail: "test@example.com",
		SupportPhone: "0912345678",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, domain.SupplierStatusActive, result.Status)
	require.True(t, result.IsActive)
}
