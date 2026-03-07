package database

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type TxProvider interface {
	RunInTx(ctx context.Context, fn func(Executor) error) error
}

type txProvider struct {
	db *sqlx.DB
}

func NewTxProvider(db *sqlx.DB) TxProvider {
	return txProvider{db: db}
}

func (p txProvider) RunInTx(ctx context.Context, fn func(Executor) error) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}
