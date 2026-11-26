package pgx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"image-uploader/internal/apperr"
	"image-uploader/internal/postgres"
)

type pgxRows struct {
	pgx.Rows
}

func (r pgxRows) Scan(dest ...any) error {
	return mapErrors(r.Rows.Scan(dest...))
}

func (r pgxRows) Err() error {
	return mapErrors(r.Rows.Err())
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return mapErrors(err)
	}
	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

const (
	pgxViolatesForeignKeyErrorCode = "23503"
	pgxViolatesUniqueErrorCode     = "23505"
)

func mapErrors(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %w", postgres.ErrNoRows, err)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgxViolatesForeignKeyErrorCode {
			return fmt.Errorf("%w: %w", postgres.ErrViolatesForeignKey, err)
		}
		if pgErr.Code == pgxViolatesUniqueErrorCode {
			return fmt.Errorf("%w: %w", postgres.ErrViolatesUnique, err)
		}
		mapped := fmt.Errorf("%w: %w", postgres.ErrUnknown, err)
		if strings.HasPrefix(pgErr.Code, "08") {
			return apperr.NewUnavailableError(mapped)
		}
		return mapped
	}

	mapped := fmt.Errorf("%w: %w", postgres.ErrUnknown, err)
	var opErr *net.OpError
	var dnsErr *net.DNSError
	if errors.As(err, &opErr) || errors.As(err, &dnsErr) || pgconn.Timeout(err) ||
		errors.Is(err, context.DeadlineExceeded) || errors.Is(err, net.ErrClosed) ||
		errors.Is(err, pgconn.ErrConnClosed) || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNREFUSED) ||
		errors.Is(err, syscall.ECONNABORTED) || errors.Is(err, syscall.EPIPE) {
		return apperr.NewUnavailableError(mapped)
	}
	return mapped
}

type pgxTx struct {
	pgx.Tx
	opTimeout time.Duration
}

func newTx(tx pgx.Tx, opTimeout time.Duration) pgxTx {
	return pgxTx{
		tx,
		opTimeout,
	}
}

func (t pgxTx) Query(ctx context.Context, sql string, args ...any) (postgres.Rows, error) {
	rows, err := t.Tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, mapErrors(err)
	}
	return pgxRows{rows}, nil
}

func (t pgxTx) QueryRow(ctx context.Context, sql string, args ...any) postgres.Row {
	row := t.Tx.QueryRow(ctx, sql, args...)
	return pgxRow{row}
}

func (t pgxTx) Exec(ctx context.Context, sql string, args ...any) (postgres.CommandTag, error) {
	tag, err := t.Tx.Exec(ctx, sql, args...)
	if err != nil {
		return nil, mapErrors(err)
	}
	return pgxCommandTag{tag}, nil
}

func (t pgxTx) OpTimeout() time.Duration {
	return t.opTimeout
}

func (t pgxTx) Commit(ctx context.Context) error {
	return mapErrors(t.Tx.Commit(ctx))
}

func (t pgxTx) Rollback(ctx context.Context) error {
	return mapErrors(t.Tx.Rollback(ctx))
}
