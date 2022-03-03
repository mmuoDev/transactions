package db

import (
	pg "github.com/mmuoDev/transactions/pkg/postgres"
	"github.com/pkg/errors"
)

const (
	transactionTable = "transactions"
)

//InsertTransactionFunc provides a functionality to insert data into table
type InsertTransactionFunc func(data map[string]interface{}) (int64, error)

//AddWalletFunc provides a functionality to create a wallet
type AddWalletFunc func(data map[string]interface{}) error

//InsertTransaction inserts a transaction in a table
func InsertTransaction(dbConnector pg.Connector) InsertTransactionFunc {
	return func(data map[string]interface{}) (int64, error) {
		id, err := dbConnector.Insert(transactionTable, data)
		if err != nil {
			return 0, errors.Wrap(err, "db - unable to insert record")
		}
		return id, nil
	}
}

