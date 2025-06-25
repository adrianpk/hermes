package sqlite

import (
	"context"
	"errors"
	"fmt"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	key = am.Key
)

type HermesRepo struct {
	*am.BaseRepo
	db *sqlx.DB
}

func NewHermesRepo(qm *am.QueryManager, opts ...am.Option) *HermesRepo {
	return &HermesRepo{
		BaseRepo: am.NewRepo("sqlite-auth-repo", qm, opts...),
	}
}

// Setup the database connection.
func (repo *HermesRepo) Setup(ctx context.Context) error {
	dsn, ok := repo.Cfg().StrVal(key.DBSQLiteDSN)
	if !ok {
		return errors.New("database DSN not found in configuration")
	}

	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return err
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return fmt.Errorf("failed to set WAL mode: %w", err)
	}
	repo.db = db
	return nil
}

// Stop closes the database connection.
func (repo *HermesRepo) Stop(ctx context.Context) error {
	if repo.db != nil {
		return repo.db.Close()
	}
	return nil
}

// DB returns the underlying *sqlx.DB for transaction management.
func (repo *HermesRepo) DB() *sqlx.DB {
	return repo.db
}

// getExec returns the correct Execer (Tx or DB) from context.
func (repo *HermesRepo) getExec(ctx context.Context) sqlx.ExtContext {
	tx, ok := am.TxFromContext(ctx)
	if ok {
		if sqlxTx, ok := tx.(*sqlx.Tx); ok {
			return sqlxTx
		}
	}
	return repo.db
}

func (r *HermesRepo) BeginTx(ctx context.Context) (context.Context, am.Tx, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return ctx, nil, err
	}
	ctxWithTx := am.WithTx(ctx, tx)
	return ctxWithTx, tx, nil
}
