package domain

import "time"

// BankAccountType represents the type of bank account.
type BankAccountType string

const (
	BankAccountTypeChecking BankAccountType = "checking"
	BankAccountTypeSavings  BankAccountType = "savings"
	BankAccountTypeCurrent  BankAccountType = "current"
)

// BankAccount represents a supplier's bank account information.
type BankAccount struct {
	ID              int64
	SupplierID      int64
	BankName        string
	AccountNumber   string
	AccountHolderName string
	BranchName      string
	AccountType     BankAccountType
	IsPrimary       bool
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewBankAccount creates a new bank account with validated inputs.
func NewBankAccount(supplierID int64, bankName, accountNumber, accountHolderName, branchName, accountType string, isPrimary bool) (*BankAccount, error) {
	if supplierID <= 0 {
		return nil, ErrInvalidSupplierID
	}
	if bankName == "" {
		return nil, ErrInvalidBankName
	}
	if accountNumber == "" {
		return nil, ErrInvalidAccountNumber
	}
	if accountHolderName == "" {
		return nil, ErrInvalidAccountHolderName
	}

	var parsedType BankAccountType
	switch accountType {
	case "checking", "savings", "current":
		parsedType = BankAccountType(accountType)
	default:
		parsedType = BankAccountTypeChecking // default
	}

	now := time.Now()
	return &BankAccount{
		SupplierID:        supplierID,
		BankName:          bankName,
		AccountNumber:     accountNumber,
		AccountHolderName: accountHolderName,
		BranchName:        branchName,
		AccountType:       parsedType,
		IsPrimary:         isPrimary,
		IsActive:          true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

// Update modifies bank account fields in-place.
func (b *BankAccount) Update(bankName, accountNumber, accountHolderName, branchName, accountType string, isPrimary, isActive bool) error {
	if bankName == "" {
		return ErrInvalidBankName
	}
	if accountNumber == "" {
		return ErrInvalidAccountNumber
	}
	if accountHolderName == "" {
		return ErrInvalidAccountHolderName
	}

	var parsedType BankAccountType
	switch accountType {
	case "checking", "savings", "current":
		parsedType = BankAccountType(accountType)
	default:
		parsedType = BankAccountTypeChecking
	}

	b.BankName = bankName
	b.AccountNumber = accountNumber
	b.AccountHolderName = accountHolderName
	b.BranchName = branchName
	b.AccountType = parsedType
	b.IsPrimary = isPrimary
	b.IsActive = isActive
	b.UpdatedAt = time.Now()
	return nil
}

var (
	ErrInvalidSupplierID      = &ValidationError{Field: "supplier_id", Message: "supplier_id is required"}
	ErrInvalidBankName        = &ValidationError{Field: "bank_name", Message: "bank_name is required"}
	ErrInvalidAccountNumber   = &ValidationError{Field: "account_number", Message: "account_number is required"}
	ErrInvalidAccountHolderName = &ValidationError{Field: "account_holder_name", Message: "account_holder_name is required"}
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

