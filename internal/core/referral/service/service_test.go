package referral

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	goodmoney "github.com/the-nucleus-project/good_money"
	"marketplace/internal/core/referral/domain"
	"marketplace/pkg/pagination"

	"github.com/stretchr/testify/require"
)

type mockReferralRepository struct {
	affiliates           map[int64]*domain.Affiliate
	affiliatesByEmail    map[string]int64
	affiliatesByPhone    map[string]int64
	referralCodes        map[int64]*domain.ReferralCode
	referralCodesByCode  map[string]int64
	relationships        map[int64]*domain.ReferralRelationship
	relationshipsByCust  map[int64]int64
	commissions          map[int64]*domain.Commission
	nextID               int64
	createAffiliateFunc  func(context.Context, *domain.Affiliate) error
	findAffiliateByIDFunc func(context.Context, int64) (*domain.Affiliate, error)
}

func newMockReferralRepository() *mockReferralRepository {
	return &mockReferralRepository{
		affiliates:          make(map[int64]*domain.Affiliate),
		affiliatesByEmail:   make(map[string]int64),
		affiliatesByPhone:   make(map[string]int64),
		referralCodes:       make(map[int64]*domain.ReferralCode),
		referralCodesByCode: make(map[string]int64),
		relationships:      make(map[int64]*domain.ReferralRelationship),
		relationshipsByCust: make(map[int64]int64),
		commissions:         make(map[int64]*domain.Commission),
		nextID:             1,
	}
}

func (m *mockReferralRepository) CreateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error {
	if m.createAffiliateFunc != nil {
		return m.createAffiliateFunc(ctx, affiliate)
	}
	if affiliate.ID == 0 {
		affiliate.ID = m.nextID
		m.nextID++
	}
	m.affiliates[affiliate.ID] = affiliate
	if affiliate.Email != "" {
		m.affiliatesByEmail[affiliate.Email] = affiliate.ID
	}
	if affiliate.PhoneNumber != "" {
		m.affiliatesByPhone[affiliate.PhoneNumber] = affiliate.ID
	}
	return nil
}

func (m *mockReferralRepository) UpdateAffiliate(ctx context.Context, affiliate *domain.Affiliate) error {
	if _, ok := m.affiliates[affiliate.ID]; !ok {
		return ErrAffiliateNotFound
	}
	m.affiliates[affiliate.ID] = affiliate
	return nil
}

func (m *mockReferralRepository) FindAffiliateByID(ctx context.Context, id int64) (*domain.Affiliate, error) {
	if m.findAffiliateByIDFunc != nil {
		return m.findAffiliateByIDFunc(ctx, id)
	}
	affiliate, ok := m.affiliates[id]
	if !ok {
		return nil, ErrAffiliateNotFound
	}
	return affiliate, nil
}

func (m *mockReferralRepository) FindAffiliateByEmail(ctx context.Context, email string) (*domain.Affiliate, error) {
	if id, ok := m.affiliatesByEmail[email]; ok {
		return m.affiliates[id], nil
	}
	return nil, ErrAffiliateNotFound
}

func (m *mockReferralRepository) FindAffiliateByPhone(ctx context.Context, phone string) (*domain.Affiliate, error) {
	if id, ok := m.affiliatesByPhone[phone]; ok {
		return m.affiliates[id], nil
	}
	return nil, ErrAffiliateNotFound
}

func (m *mockReferralRepository) FindAffiliateByCustomerID(ctx context.Context, customerID int64) (*domain.Affiliate, error) {
	for _, affiliate := range m.affiliates {
		if affiliate.AffiliateID != nil && *affiliate.AffiliateID == customerID {
			return affiliate, nil
		}
	}
	return nil, ErrAffiliateNotFound
}

func (m *mockReferralRepository) DeleteAffiliate(ctx context.Context, id int64) error {
	if _, ok := m.affiliates[id]; !ok {
		return ErrAffiliateNotFound
	}
	delete(m.affiliates, id)
	return nil
}

func (m *mockReferralRepository) FindAllAffiliates(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Affiliate], error) {
	items := make([]*domain.Affiliate, 0, len(m.affiliates))
	for _, affiliate := range m.affiliates {
		items = append(items, affiliate)
	}
	return pagination.NewPageResult(items, len(items), pageReq), nil
}

func (m *mockReferralRepository) CreateReferralCode(ctx context.Context, code *domain.ReferralCode) error {
	if code.ID == 0 {
		code.ID = m.nextID
		m.nextID++
	}
	m.referralCodes[code.ID] = code
	m.referralCodesByCode[code.Code] = code.ID
	return nil
}

func (m *mockReferralRepository) UpdateReferralCode(ctx context.Context, code *domain.ReferralCode) error {
	if _, ok := m.referralCodes[code.ID]; !ok {
		return ErrReferralCodeNotFound
	}
	m.referralCodes[code.ID] = code
	return nil
}

func (m *mockReferralRepository) FindReferralCodeByID(ctx context.Context, id int64) (*domain.ReferralCode, error) {
	code, ok := m.referralCodes[id]
	if !ok {
		return nil, ErrReferralCodeNotFound
	}
	return code, nil
}

func (m *mockReferralRepository) FindReferralCodeByCode(ctx context.Context, code string) (*domain.ReferralCode, error) {
	if id, ok := m.referralCodesByCode[code]; ok {
		return m.referralCodes[id], nil
	}
	return nil, ErrReferralCodeNotFound
}

func (m *mockReferralRepository) FindReferralCodesByAffiliateID(ctx context.Context, affiliateID int64) ([]*domain.ReferralCode, error) {
	var codes []*domain.ReferralCode
	for _, code := range m.referralCodes {
		if code.AffiliateID == affiliateID {
			codes = append(codes, code)
		}
	}
	return codes, nil
}

func (m *mockReferralRepository) DeleteReferralCode(ctx context.Context, id int64) error {
	if _, ok := m.referralCodes[id]; !ok {
		return ErrReferralCodeNotFound
	}
	delete(m.referralCodes, id)
	return nil
}

func (m *mockReferralRepository) CreateReferralRelationship(ctx context.Context, relationship *domain.ReferralRelationship) error {
	if relationship.ID == 0 {
		relationship.ID = m.nextID
		m.nextID++
	}
	m.relationships[relationship.ID] = relationship
	m.relationshipsByCust[relationship.CustomerID] = relationship.ID
	return nil
}

func (m *mockReferralRepository) UpdateReferralRelationship(ctx context.Context, relationship *domain.ReferralRelationship) error {
	if _, ok := m.relationships[relationship.ID]; !ok {
		return ErrReferralRelationshipNotFound
	}
	m.relationships[relationship.ID] = relationship
	return nil
}

func (m *mockReferralRepository) FindReferralRelationshipByID(ctx context.Context, id int64) (*domain.ReferralRelationship, error) {
	relationship, ok := m.relationships[id]
	if !ok {
		return nil, ErrReferralRelationshipNotFound
	}
	return relationship, nil
}

func (m *mockReferralRepository) FindReferralRelationshipByCustomerID(ctx context.Context, customerID int64) (*domain.ReferralRelationship, error) {
	if id, ok := m.relationshipsByCust[customerID]; ok {
		return m.relationships[id], nil
	}
	return nil, ErrReferralRelationshipNotFound
}

func (m *mockReferralRepository) FindReferralRelationshipsByAffiliateID(ctx context.Context, affiliateID int64, pageReq pagination.PageRequest) (pagination.PageResult[*domain.ReferralRelationship], error) {
	var items []*domain.ReferralRelationship
	for _, rel := range m.relationships {
		if rel.AffiliateID == affiliateID {
			items = append(items, rel)
		}
	}
	return pagination.NewPageResult(items, len(items), pageReq), nil
}

func (m *mockReferralRepository) DeleteReferralRelationship(ctx context.Context, id int64) error {
	if _, ok := m.relationships[id]; !ok {
		return ErrReferralRelationshipNotFound
	}
	delete(m.relationships, id)
	return nil
}

func (m *mockReferralRepository) CreateCommission(ctx context.Context, commission *domain.Commission) error {
	if commission.ID == 0 {
		commission.ID = m.nextID
		m.nextID++
	}
	m.commissions[commission.ID] = commission
	return nil
}

func (m *mockReferralRepository) UpdateCommission(ctx context.Context, commission *domain.Commission) error {
	if _, ok := m.commissions[commission.ID]; !ok {
		return ErrCommissionNotFound
	}
	m.commissions[commission.ID] = commission
	return nil
}

func (m *mockReferralRepository) FindCommissionByID(ctx context.Context, id int64) (*domain.Commission, error) {
	commission, ok := m.commissions[id]
	if !ok {
		return nil, ErrCommissionNotFound
	}
	return commission, nil
}

func (m *mockReferralRepository) FindCommissionsByAffiliateID(ctx context.Context, affiliateID int64, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Commission], error) {
	var items []*domain.Commission
	for _, comm := range m.commissions {
		if comm.AffiliateID == affiliateID {
			items = append(items, comm)
		}
	}
	return pagination.NewPageResult(items, len(items), pageReq), nil
}

func (m *mockReferralRepository) FindCommissionsByOrderID(ctx context.Context, orderID int64) ([]*domain.Commission, error) {
	var items []*domain.Commission
	for _, comm := range m.commissions {
		if comm.OrderID == orderID {
			items = append(items, comm)
		}
	}
	return items, nil
}

func (m *mockReferralRepository) FindCommissionsByReferralRelationshipID(ctx context.Context, referralRelationshipID int64) ([]*domain.Commission, error) {
	var items []*domain.Commission
	for _, comm := range m.commissions {
		if comm.ReferralRelationshipID == referralRelationshipID {
			items = append(items, comm)
		}
	}
	return items, nil
}

func (m *mockReferralRepository) DeleteCommission(ctx context.Context, id int64) error {
	if _, ok := m.commissions[id]; !ok {
		return ErrCommissionNotFound
	}
	delete(m.commissions, id)
	return nil
}

func TestReferralServiceCreateAffiliateSuccess(t *testing.T) {
	repo := newMockReferralRepository()
	service := NewReferralService(repo)

	result, err := service.CreateAffiliate(context.Background(), AffiliateInput{
		FullName:    "John Doe",
		Email:       "john@example.com",
		PhoneNumber: "0912345678",
		Status:      "active",
		Platform:   "Instagram",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotZero(t, result.ID)
	require.Equal(t, domain.AffiliateStatusActive, result.Status)
	require.True(t, result.IsActive)
}

func TestReferralServiceCreateAffiliateEmailConflict(t *testing.T) {
	repo := newMockReferralRepository()
	existing := &domain.Affiliate{
		ID:    1,
		Email: "dup@example.com",
	}
	repo.affiliates[1] = existing
	repo.affiliatesByEmail["dup@example.com"] = 1

	service := NewReferralService(repo)
	_, err := service.CreateAffiliate(context.Background(), AffiliateInput{
		FullName: "New Person",
		Email:    "dup@example.com",
		Status:   "active",
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrAffiliateEmailExists))
}

func TestReferralServiceCreateReferralCodeSuccess(t *testing.T) {
	repo := newMockReferralRepository()
	affiliate := &domain.Affiliate{
		ID:    1,
		Email: "affiliate@example.com",
	}
	repo.affiliates[1] = affiliate

	service := NewReferralService(repo)
	result, err := service.CreateReferralCode(context.Background(), ReferralCodeInput{
		AffiliateID: 1,
		CustomCode:  "TEST123",
		Status:      "active",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "TEST123", result.Code)
	require.True(t, result.IsCustom)
}

func TestReferralServiceCreateReferralCodeAutoGenerate(t *testing.T) {
	repo := newMockReferralRepository()
	affiliate := &domain.Affiliate{
		ID:    1,
		Email: "affiliate@example.com",
	}
	repo.affiliates[1] = affiliate

	service := NewReferralService(repo)
	result, err := service.CreateReferralCode(context.Background(), ReferralCodeInput{
		AffiliateID: 1,
		Status:      "active",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result.Code)
	require.False(t, result.IsCustom)
}

func TestReferralServiceCreateReferralCodeDuplicate(t *testing.T) {
	repo := newMockReferralRepository()
	affiliate := &domain.Affiliate{
		ID:    1,
		Email: "affiliate@example.com",
	}
	repo.affiliates[1] = affiliate
	existingCode := &domain.ReferralCode{
		ID:    1,
		Code:  "EXISTING",
		AffiliateID: 1,
	}
	repo.referralCodes[1] = existingCode
	repo.referralCodesByCode["EXISTING"] = 1

	service := NewReferralService(repo)
	_, err := service.CreateReferralCode(context.Background(), ReferralCodeInput{
		AffiliateID: 1,
		CustomCode:  "EXISTING",
		Status:      "active",
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrReferralCodeExists))
}

func TestReferralServiceCreateReferralRelationshipSuccess(t *testing.T) {
	repo := newMockReferralRepository()
	affiliate := &domain.Affiliate{ID: 1}
	code := &domain.ReferralCode{ID: 1, AffiliateID: 1, Code: "TEST123"}
	repo.affiliates[1] = affiliate
	repo.referralCodes[1] = code

	service := NewReferralService(repo)
	result, err := service.CreateReferralRelationship(context.Background(), ReferralRelationshipInput{
		AffiliateID:    1,
		CustomerID:     2,
		ReferralCodeID: 1,
		Source:        "Instagram",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, int64(1), result.AffiliateID)
	require.Equal(t, int64(2), result.CustomerID)
}

func TestReferralServiceCreateCommissionSuccess(t *testing.T) {
	repo := newMockReferralRepository()
	relationship := &domain.ReferralRelationship{
		ID:         1,
		AffiliateID: 1,
		CustomerID:  2,
	}
	repo.relationships[1] = relationship

	orderTotal, _ := goodmoney.New(10000.0, "ETB") // 100.00 ETB
	orderTotalJSON, _ := json.Marshal(orderTotal)

	service := NewReferralService(repo)
	result, err := service.CreateCommission(context.Background(), CommissionInput{
		AffiliateID:            1,
		ReferralRelationshipID: 1,
		OrderID:                10,
		CustomerID:             2,
		CommissionRate:         10.0,
		OrderTotal:             string(orderTotalJSON),
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, int64(1), result.AffiliateID)
	require.Equal(t, 10.0, result.CommissionRate)
	require.NotNil(t, result.Amount)
}

func TestReferralServiceGetAffiliateNotFound(t *testing.T) {
	repo := newMockReferralRepository()
	service := NewReferralService(repo)

	_, err := service.GetAffiliate(context.Background(), 999)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrAffiliateNotFound))
}

func TestReferralServiceListAffiliates(t *testing.T) {
	repo := newMockReferralRepository()
	repo.affiliates[1] = &domain.Affiliate{ID: 1, FullName: "Affiliate 1"}
	repo.affiliates[2] = &domain.Affiliate{ID: 2, FullName: "Affiliate 2"}

	service := NewReferralService(repo)
	result, err := service.ListAffiliates(context.Background(), pagination.PageRequest{
		Page:  1,
		Limit: 10,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 2)
	require.GreaterOrEqual(t, len(result.Items), 2)
}

