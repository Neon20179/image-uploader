package postgres

import (
	"context"
	"fmt"
	"time"
)

type TxManager struct {
	pool Pool
}

type rollbackError struct {
	primary error
	cleanup error
}

func (e *rollbackError) Error() string {
	return fmt.Sprintf("%v; rollback tx: %v", e.primary, e.cleanup)
}
func (e *rollbackError) Unwrap() error {
	return e.primary
}

func rollback(ctx context.Context, tx Tx) error {
	timeout := tx.OpTimeout()
	if timeout <= 0 {
		timeout = time.Second
	}
	cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), timeout)
	defer cancel()
	return tx.Rollback(cleanup)
}

func NewTxManager(pool Pool) *TxManager {
	return &TxManager{
		pool: pool,
	}
}

func (m *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	txCtx := ContextWithQuerier(ctx, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = rollback(ctx, tx)
			panic(p)
		}
	}()

	finish := func(primary error) error {
		if rbErr := rollback(ctx, tx); rbErr != nil {
			return &rollbackError{primary: primary, cleanup: rbErr}
		}
		return primary
	}
	if err := fn(txCtx); err != nil {
		return finish(err)
	}

	return tx.Commit(ctx)
}
