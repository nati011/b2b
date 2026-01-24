package supplier

import (
	"context"
	"errors"

	"marketplace/internal/core/supplier/domain"
	"marketplace/pkg/logger"
	"marketplace/pkg/pagination"
)

var (
	// ErrSupplierNotFound indicates a supplier lookup failure.
	ErrSupplierNotFound = errors.New("supplier not found")
	// ErrSupplierEmailExists indicates a supplier with the email already exists.
	ErrSupplierEmailExists = errors.New("supplier email already exists")
	// ErrSupplierPhoneExists indicates a supplier with the phone number already exists.
	ErrSupplierPhoneExists = errors.New("supplier phone number already exists")
)

// Repository defines the supplier persistence contract.
type Repository interface {
	Create(ctx context.Context, supplier *domain.Supplier) error
	Update(ctx context.Context, supplier *domain.Supplier) error
	FindByID(ctx context.Context, id int64) (*domain.Supplier, error)
	FindByEmail(ctx context.Context, email string) (*domain.Supplier, error)
	FindByPhone(ctx context.Context, phone string) (*domain.Supplier, error)
	Delete(ctx context.Context, id int64) error
	FindAll(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Supplier], error)
}

// SupplierInput captures incoming supplier profile fields.
type SupplierInput struct {
	BusinessName string
	Status       string
	SupportEmail string
	SupportPhone string
}

// SupplierService coordinates supplier business logic.
type SupplierService struct {
	repository Repository
}

func NewSupplierService(repository Repository) *SupplierService {
	return &SupplierService{repository: repository}
}

// Create registers a new supplier.
func (s *SupplierService) Create(ctx context.Context, input SupplierInput) (*domain.Supplier, error) {
	if err := s.ensureEmailAvailable(ctx, input.SupportEmail, nil); err != nil {
		return nil, err
	}
	if err := s.ensurePhoneAvailable(ctx, input.SupportPhone, nil); err != nil {
		return nil, err
	}

	supplier, err := domain.NewSupplier(
		input.BusinessName,
		input.Status,
		input.SupportEmail,
		input.SupportPhone,
	)
	if err != nil {
		logger.Warn("Supplier creation failed: validation error", "business_name", input.BusinessName, "error", err)
		return nil, err
	}

	if err := s.repository.Create(ctx, supplier); err != nil {
		logger.Error("Supplier creation failed: repository error", "business_name", input.BusinessName, "error", err)
		return nil, err
	}

	logger.Info("Supplier created successfully", "supplier_id", supplier.ID)
	return supplier, nil
}

// Get fetches a supplier by identifier.
func (s *SupplierService) Get(ctx context.Context, id int64) (*domain.Supplier, error) {
	supplier, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSupplierNotFound) {
			logger.Debug("Supplier retrieval failed: supplier not found", "supplier_id", id)
		} else {
			logger.Error("Supplier retrieval failed: repository error", "supplier_id", id, "error", err)
		}
		return nil, err
	}
	logger.Debug("Supplier retrieved successfully", "supplier_id", id)
	return supplier, nil
}

// Update modifies an existing supplier.
func (s *SupplierService) Update(ctx context.Context, id int64, input SupplierInput) (*domain.Supplier, error) {
	supplier, err := s.repository.FindByID(ctx, id)
	if err != nil {
		logger.Debug("Supplier update failed: supplier not found", "supplier_id", id)
		return nil, err
	}

	if input.SupportEmail != "" && input.SupportEmail != supplier.SupportEmail {
		if err := s.ensureEmailAvailable(ctx, input.SupportEmail, &id); err != nil {
			return nil, err
		}
	}
	if input.SupportPhone != "" && input.SupportPhone != supplier.SupportPhone {
		if err := s.ensurePhoneAvailable(ctx, input.SupportPhone, &id); err != nil {
			return nil, err
		}
	}

	if err := supplier.Update(
		input.BusinessName,
		input.Status,
		input.SupportEmail,
		input.SupportPhone,
	); err != nil {
		logger.Warn("Supplier update failed: validation error", "supplier_id", id, "error", err)
		return nil, err
	}

	if err := s.repository.Update(ctx, supplier); err != nil {
		logger.Error("Supplier update failed: repository error", "supplier_id", id, "error", err)
		return nil, err
	}

	logger.Info("Supplier updated successfully", "supplier_id", id)
	return supplier, nil
}

// Delete removes a supplier (soft delete).
func (s *SupplierService) Delete(ctx context.Context, id int64) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrSupplierNotFound) {
			logger.Debug("Supplier deletion failed: supplier not found", "supplier_id", id)
		} else {
			logger.Error("Supplier deletion failed: repository error", "supplier_id", id, "error", err)
		}
		return err
	}

	logger.Info("Supplier deleted successfully", "supplier_id", id)
	return nil
}

// List retrieves a paginated list of suppliers.
func (s *SupplierService) List(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Supplier], error) {
	result, err := s.repository.FindAll(ctx, pageReq)
	if err != nil {
		logger.Error("Supplier pagination failed: repository error", "page", pageReq.Page, "limit", pageReq.Limit, "error", err)
		return pagination.PageResult[*domain.Supplier]{}, err
	}
	logger.Debug("Supplier pagination completed", "page", result.Page, "total", result.Total, "items", len(result.Items))
	return result, nil
}

func (s *SupplierService) ensureEmailAvailable(ctx context.Context, email string, excludeID *int64) error {
	if email == "" {
		return nil
	}
	existing, err := s.repository.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		if excludeID == nil || existing.ID != *excludeID {
			logger.Warn("Supplier email conflict detected", "email", email)
			return ErrSupplierEmailExists
		}
	}
	if err != nil && !errors.Is(err, ErrSupplierNotFound) {
		logger.Error("Supplier email lookup failed", "email", email, "error", err)
		return err
	}
	return nil
}

func (s *SupplierService) ensurePhoneAvailable(ctx context.Context, phoneNumber string, excludeID *int64) error {
	if phoneNumber == "" {
		return nil
	}
	existing, err := s.repository.FindByPhone(ctx, phoneNumber)
	if err == nil && existing != nil {
		if excludeID == nil || existing.ID != *excludeID {
			logger.Warn("Supplier phone conflict detected", "phone_number", phoneNumber)
			return ErrSupplierPhoneExists
		}
	}
	if err != nil && !errors.Is(err, ErrSupplierNotFound) {
		logger.Error("Supplier phone lookup failed", "phone_number", phoneNumber, "error", err)
		return err
	}
	return nil
}
