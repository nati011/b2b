package transaction_control

import "context"

type Transaction interface {
	Rollback() error
	Commit() error
}

type TransactionManager struct {
	rollbackRoutines []func()
}

func NewTransactionManager() TransactionManager {
	return TransactionManager{}
}

func (t *TransactionManager) WithinTransaction(ctx context.Context, txFunc func(context.Context) (interface{}, error)) (interface{}, error) {
	t.Init(ctx)
	defer func() {
		if r := recover(); r != nil {
			t.Rollback(ctx)
			panic(r)
		}
	}()

	result, err := txFunc(ctx)
	if err != nil {
		t.Rollback(ctx)
		return nil, err
	}

	if err := t.Commit(ctx); err != nil {
		return nil, err
	}

	return result, nil
}

func (t *TransactionManager) Init(ctx context.Context) error {
	return nil
}

func (t *TransactionManager) Commit(ctx context.Context) error {
	return nil
}

func (t *TransactionManager) Rollback(ctx context.Context) error {
	for _, i := range t.rollbackRoutines {
		i()
	}
	return nil
}

func (t *TransactionManager) Register(routine func()) {
	t.rollbackRoutines = append(t.rollbackRoutines, routine)
}
