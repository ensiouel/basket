package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type clientKey struct{}

type Tx struct {
	pgx.Tx
}

func WithContext(ctx context.Context, client Client) context.Context {
	return context.WithValue(ctx, clientKey{}, client)
}

func (t *Tx) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, t.Tx, dest, query, args...)
}

func (t *Tx) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, t.Tx, dest, query, args...)
}

func (t *Tx) FromContext(ctx context.Context) Client {
	if tx, ok := ctx.Value(clientKey{}).(Client); ok {
		return tx
	}

	return t
}

func WithinTransaction(ctx context.Context, client Client, txFunc func(ctx context.Context) error) error {
	tx, err := client.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)

			panic(r)
		}
	}()

	err = txFunc(WithContext(ctx, &Tx{tx}))
	if err != nil {
		_ = tx.Rollback(ctx)

		return err
	}

	return tx.Commit(ctx)
}
