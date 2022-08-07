package db

import (
	"database/sql"
)

type Repo interface {
	Querier
}

type SQLRepo struct {
	*Queries
	db *sql.DB
}

func NewRepo(db *sql.DB) Repo {
	return &SQLRepo{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
// func (repo *SQLRepo) execTx(ctx context.Context, fn func(*Queries) error) error {
// 	tx, err := repo.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}

// 	q := New(tx)
// 	err = fn(q)
// 	if err != nil {
// 		if rbErr := tx.Rollback(); rbErr != nil {
// 			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
// 		}

// 	}
// 	return tx.Commit()
// }
