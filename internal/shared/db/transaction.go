package db

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/shared/ports/outbound"
)

type TransactionManager struct {
	db *DB
}

func NewTransactionManager(database *DB) *TransactionManager {
	return &TransactionManager{
		db: database,
	}
}

func (m *TransactionManager) WithinTransaction(
	ctx context.Context,
	fn func(context.Context) error,
) error {

	tx, err := m.db.conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	ctx = WithExecutor(ctx, tx)

	if err := fn(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

var _ outbound.TransactionManager = (*TransactionManager)(nil)