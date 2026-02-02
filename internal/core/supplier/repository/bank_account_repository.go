package supplier

import (
	"context"
	"database/sql"
	"errors"

	"marketplace/internal/core/supplier/domain"
)

var (
	ErrBankAccountNotFound = errors.New("bank account not found")
)

// BankAccountRepository persists bank accounts.
type BankAccountRepository struct {
	db *sql.DB
}

func NewBankAccountRepository(db *sql.DB) *BankAccountRepository {
	return &BankAccountRepository{db: db}
}

// Create inserts a new bank account record.
func (r *BankAccountRepository) Create(ctx context.Context, account *domain.BankAccount) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// If this is set as primary, unset other primary accounts for this supplier
	if account.IsPrimary {
		_, err = tx.ExecContext(
			ctx,
			`UPDATE bank_accounts SET is_primary = FALSE, last_modified = NOW() 
			 WHERE supplier_id = $1 AND is_deleted = FALSE AND is_primary = TRUE`,
			account.SupplierID,
		)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO bank_accounts (
			supplier_id,
			bank_name,
			account_number,
			account_holder_name,
			branch_name,
			account_type,
			is_primary,
			is_active,
			created_date,
			last_modified,
			is_deleted
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, FALSE)
		RETURNING id
	`

	err = tx.QueryRowContext(
		ctx,
		query,
		account.SupplierID,
		account.BankName,
		account.AccountNumber,
		account.AccountHolderName,
		nullableString(account.BranchName),
		string(account.AccountType),
		account.IsPrimary,
		account.IsActive,
		account.CreatedAt,
		account.UpdatedAt,
	).Scan(&account.ID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// Update modifies an existing bank account record.
func (r *BankAccountRepository) Update(ctx context.Context, account *domain.BankAccount) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// If this is set as primary, unset other primary accounts for this supplier
	if account.IsPrimary {
		_, err = tx.ExecContext(
			ctx,
			`UPDATE bank_accounts SET is_primary = FALSE, last_modified = NOW() 
			 WHERE supplier_id = $1 AND id != $2 AND is_deleted = FALSE AND is_primary = TRUE`,
			account.SupplierID,
			account.ID,
		)
		if err != nil {
			return err
		}
	}

	query := `
		UPDATE bank_accounts
		SET bank_name = $2,
		    account_number = $3,
		    account_holder_name = $4,
		    branch_name = $5,
		    account_type = $6,
		    is_primary = $7,
		    is_active = $8,
		    last_modified = $9
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		account.ID,
		account.BankName,
		account.AccountNumber,
		account.AccountHolderName,
		nullableString(account.BranchName),
		string(account.AccountType),
		account.IsPrimary,
		account.IsActive,
		account.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrBankAccountNotFound
	}

	return tx.Commit()
}

// FindByID fetches a bank account by identifier.
func (r *BankAccountRepository) FindByID(ctx context.Context, id int64) (*domain.BankAccount, error) {
	query := `
		SELECT id, supplier_id, bank_name, account_number, account_holder_name, 
		       branch_name, account_type, is_primary, is_active, created_date, last_modified
		FROM bank_accounts
		WHERE id = $1 AND is_deleted = FALSE
	`
	account, err := scanBankAccount(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBankAccountNotFound
		}
		return nil, err
	}
	return account, nil
}

// FindBySupplierID fetches all bank accounts for a supplier.
func (r *BankAccountRepository) FindBySupplierID(ctx context.Context, supplierID int64) ([]*domain.BankAccount, error) {
	query := `
		SELECT id, supplier_id, bank_name, account_number, account_holder_name, 
		       branch_name, account_type, is_primary, is_active, created_date, last_modified
		FROM bank_accounts
		WHERE supplier_id = $1 AND is_deleted = FALSE
		ORDER BY is_primary DESC, created_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, supplierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*domain.BankAccount
	for rows.Next() {
		account, err := scanBankAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

// Delete performs a soft delete on the bank account.
func (r *BankAccountRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE bank_accounts
		SET is_deleted = TRUE, last_modified = NOW()
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrBankAccountNotFound
	}

	return nil
}

type bankAccountScanner interface {
	Scan(dest ...interface{}) error
}

func scanBankAccount(scanner bankAccountScanner) (*domain.BankAccount, error) {
	var account domain.BankAccount
	var branchName, accountType sql.NullString
	if err := scanner.Scan(
		&account.ID,
		&account.SupplierID,
		&account.BankName,
		&account.AccountNumber,
		&account.AccountHolderName,
		&branchName,
		&accountType,
		&account.IsPrimary,
		&account.IsActive,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if branchName.Valid {
		account.BranchName = branchName.String
	}
	if accountType.Valid {
		account.AccountType = domain.BankAccountType(accountType.String)
	} else {
		account.AccountType = domain.BankAccountTypeChecking
	}

	return &account, nil
}

