package customer

import (
	"context"
	"errors"

	"marketplace/internal/core/customer/domain"
	"marketplace/pkg/logger"
	"marketplace/pkg/pagination"
)

var (
	// ErrCustomerNotFound indicates a customer lookup failure.
	ErrCustomerNotFound = errors.New("customer not found")
	// ErrCustomerEmailExists indicates a customer with the email already exists.
	ErrCustomerEmailExists = errors.New("customer email already exists")
	// ErrCustomerPhoneExists indicates a customer with the phone number already exists.
	ErrCustomerPhoneExists = errors.New("customer phone number already exists")
)

// Repository defines the customer persistence contract.
type Repository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	Update(ctx context.Context, customer *domain.Customer) error
	FindByID(ctx context.Context, id int64) (*domain.Customer, error)
	FindByEmail(ctx context.Context, email string) (*domain.Customer, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Customer, error)
	Delete(ctx context.Context, id int64) error
	FindAll(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Customer], error)
	FindAllWithSupplierFilter(ctx context.Context, pageReq pagination.PageRequest, supplierID int64, search string) (pagination.PageResult[*domain.Customer], error)
}

// CustomerInput captures incoming customer profile fields.
type CustomerInput struct {
	FullName    string
	Status      string
	City        string
	Region      string
	Woreda      string
	PhoneNumber string
	Email       string
}

// CustomerService coordinates customer business logic.
type CustomerService struct {
	repository Repository
}

func NewCustomerService(repository Repository) *CustomerService {
	return &CustomerService{repository: repository}
}

// Create registers a new customer.
func (s *CustomerService) Create(ctx context.Context, input CustomerInput) (*domain.Customer, error) {
	if err := s.ensureEmailAvailable(ctx, input.Email, nil); err != nil {
		return nil, err
	}
	if err := s.ensurePhoneAvailable(ctx, input.PhoneNumber, nil); err != nil {
		return nil, err
	}

	customer, err := domain.NewCustomer(
		input.FullName,
		input.Status,
		input.City,
		input.Region,
		input.Woreda,
		input.PhoneNumber,
		input.Email,
	)
	if err != nil {
		logger.Warn("Customer creation failed: validation error", "full_name", input.FullName, "error", err)
		return nil, err
	}

	if err := s.repository.Create(ctx, customer); err != nil {
		logger.Error("Customer creation failed: repository error", "full_name", input.FullName, "error", err)
		return nil, err
	}

	logger.Info("Customer created successfully", "customer_id", customer.ID)
	return customer, nil
}

// Get fetches a customer by identifier.
func (s *CustomerService) Get(ctx context.Context, id int64) (*domain.Customer, error) {
	customer, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			logger.Debug("Customer retrieval failed: customer not found", "customer_id", id)
		} else {
			logger.Error("Customer retrieval failed: repository error", "customer_id", id, "error", err)
		}
		return nil, err
	}
	logger.Debug("Customer retrieved successfully", "customer_id", id)
	return customer, nil
}

// GetByEmail fetches a customer by email address.
func (s *CustomerService) GetByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	customer, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			logger.Debug("Customer retrieval failed: customer not found", "email", email)
		} else {
			logger.Error("Customer retrieval failed: repository error", "email", email, "error", err)
		}
		return nil, err
	}
	logger.Debug("Customer retrieved successfully", "email", email, "customer_id", customer.ID)
	return customer, nil
}

// Update modifies an existing customer.
func (s *CustomerService) Update(ctx context.Context, id int64, input CustomerInput) (*domain.Customer, error) {
	customer, err := s.repository.FindByID(ctx, id)
	if err != nil {
		logger.Debug("Customer update failed: customer not found", "customer_id", id)
		return nil, err
	}

	if input.Email != "" && input.Email != customer.Email {
		if err := s.ensureEmailAvailable(ctx, input.Email, &id); err != nil {
			return nil, err
		}
	}
	if input.PhoneNumber != "" && input.PhoneNumber != customer.PhoneNumber {
		if err := s.ensurePhoneAvailable(ctx, input.PhoneNumber, &id); err != nil {
			return nil, err
		}
	}

	if err := customer.Update(
		input.FullName,
		input.Status,
		input.City,
		input.Region,
		input.Woreda,
		input.PhoneNumber,
		input.Email,
	); err != nil {
		logger.Warn("Customer update failed: validation error", "customer_id", id, "error", err)
		return nil, err
	}

	if err := s.repository.Update(ctx, customer); err != nil {
		logger.Error("Customer update failed: repository error", "customer_id", id, "error", err)
		return nil, err
	}

	logger.Info("Customer updated successfully", "customer_id", id)
	return customer, nil
}

// Delete removes a customer (soft delete).
func (s *CustomerService) Delete(ctx context.Context, id int64) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			logger.Debug("Customer deletion failed: customer not found", "customer_id", id)
		} else {
			logger.Error("Customer deletion failed: repository error", "customer_id", id, "error", err)
		}
		return err
	}

	logger.Info("Customer deleted successfully", "customer_id", id)
	return nil
}

// List retrieves a paginated list of customers.
func (s *CustomerService) List(ctx context.Context, pageReq pagination.PageRequest, search string) (pagination.PageResult[*domain.Customer], error) {
	return s.ListWithSupplierFilter(ctx, pageReq, 0, search)
}

// ListWithSupplierFilter retrieves a paginated list of customers, optionally filtered by supplier.
func (s *CustomerService) ListWithSupplierFilter(ctx context.Context, pageReq pagination.PageRequest, supplierID int64, search string) (pagination.PageResult[*domain.Customer], error) {
	result, err := s.repository.FindAllWithSupplierFilter(ctx, pageReq, supplierID, search)
	if err != nil {
		logger.Error("Customer pagination failed: repository error", "page", pageReq.Page, "limit", pageReq.Limit, "supplier_id", supplierID, "search", search, "error", err)
		return pagination.PageResult[*domain.Customer]{}, err
	}
	logger.Debug("Customer pagination completed", "page", result.Page, "total", result.Total, "items", len(result.Items), "supplier_id", supplierID, "search", search)
	return result, nil
}

func (s *CustomerService) ensureEmailAvailable(ctx context.Context, email string, excludeID *int64) error {
	if email == "" {
		return nil
	}
	existing, err := s.repository.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		if excludeID == nil || existing.ID != *excludeID {
			logger.Warn("Customer email conflict detected", "email", email)
			return ErrCustomerEmailExists
		}
	}
	if err != nil && !errors.Is(err, ErrCustomerNotFound) {
		logger.Error("Customer email lookup failed", "email", email, "error", err)
		return err
	}
	return nil
}

func (s *CustomerService) ensurePhoneAvailable(ctx context.Context, phoneNumber string, excludeID *int64) error {
	if phoneNumber == "" {
		return nil
	}
	existing, err := s.repository.FindByPhoneNumber(ctx, phoneNumber)
	if err == nil && existing != nil {
		if excludeID == nil || existing.ID != *excludeID {
			logger.Warn("Customer phone conflict detected", "phone_number", phoneNumber)
			return ErrCustomerPhoneExists
		}
	}
	if err != nil && !errors.Is(err, ErrCustomerNotFound) {
		logger.Error("Customer phone lookup failed", "phone_number", phoneNumber, "error", err)
		return err
	}
	return nil
}
