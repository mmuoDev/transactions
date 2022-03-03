package workflow

import (
	"strings"

	"github.com/mmuoDev/transactions/internal/db"
	"github.com/mmuoDev/transactions/pkg"
	"github.com/pkg/errors"
	"github.com/mmuoDev/wallet/gen/wallet"
)

//InsertTransactionFunc provides functionality to insert transaction
type InsertTransactionFunc func(req pkg.TransactionRequest) error

//InsertTransaction inserts a transaction
func InsertTransaction(addTransaction db.InsertTransactionFunc, walletClient wallet.WalletClient) InsertTransactionFunc {
	return func(req pkg.TransactionRequest) error {
		data := make(map[string]interface{})
		data["account_id"] = req.AccountID
		data["amount"] = req.Amount
		data["category"] = getCategory(req.Category)
		_, err := addTransaction(data)
		if err != nil {
			return errors.Wrap(err, "workflow - unable to insert transaction")
		}
		//update wallet 
		
		return nil
	}
}

//getCategory returns category
func getCategory(category string) int {
	if strings.ToLower(category) == "debit" {
		return 0
	}
	return 1
}
