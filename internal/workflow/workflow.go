package workflow

import (
	"context"
	"strings"
	"time"

	"github.com/mmuoDev/transactions/internal/db"
	"github.com/mmuoDev/transactions/pkg"
	"github.com/mmuoDev/wallet/gen/wallet"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		//create wallet
		createWallet := &wallet.CreateWalletRequest{
			AccountId:       req.AccountID,
			PreviousBalance: 0,
			CurrentBalance:  0,
			CreatedAt:       timestamppb.New(time.Now()),
			UpdatedAt:       timestamppb.New(time.Now()),
		}
		if _, err := walletClient.CreateWallet(context.Background(), createWallet); err != nil {
			return errors.Wrap(err, "unable to create wallet")
		}
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
