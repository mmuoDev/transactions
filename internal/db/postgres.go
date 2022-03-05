package db

import (
	"github.com/mmuoDev/transactions/internal"
	pg "github.com/mmuoDev/transactions/pkg/postgres"
	"github.com/pkg/errors"
)

//InsertTransaction inserts a transaction in a table
func InsertTransaction(dbConnector pg.Connector) InsertTransactionFunc {
	return func(req internal.TransactionDBRequest) (int64, error) {
		query := `INSERT INTO transactions (
			account_id, transaction_id, amount, category
		) VALUES (
			$1, $2, $3, $4
		) RETURNING id`
		res, err := dbConnector.DB.Exec(query, req.AccountID, req.TransactionId, req.Amount, req.Category)
		if err != nil {
			return 0, errors.Wrap(err, "db - unable to insert record")
		}
		id, err := res.LastInsertId()
		if err != nil {
			return 0, errors.Wrap(err, "db - unable to retrieve last inserted id")
		}
		return id, nil
	}
}
