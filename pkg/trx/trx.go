package trx

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Block func(tx *sqlx.Tx, err chan error)

func WithTx(db *sqlx.DB, block Block) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	errChan := make(chan error)
	go block(tx, errChan)

	if err := <-errChan; err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			err = fmt.Errorf("%q: %w", errTx, err)
		}
		return err
	}

	return tx.Commit()
}
