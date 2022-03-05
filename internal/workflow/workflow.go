package workflow

import (
	"context"
	"strings"

	"github.com/mmuoDev/core-proto/gen/wallet"
	"github.com/mmuoDev/transactions/internal"
	"github.com/mmuoDev/transactions/internal/db"
	"github.com/mmuoDev/transactions/pkg"
	"github.com/pkg/errors"
)

//InsertTransactionFunc provides functionality to insert transaction
type InsertTransactionFunc func(req pkg.TransactionRequest) (pkg.TransactionResponse, error)

//InsertTransaction inserts a transaction
func InsertTransaction(addTransaction db.InsertTransactionFunc, walletClient wallet.WalletClient) InsertTransactionFunc {
	return func(req pkg.TransactionRequest) (pkg.TransactionResponse, error) {
		//Ideally, we would create a new wallet separately so we are focused squarely on updating a wallet here.
		//check if wallet exists else create
		retrieve := &wallet.RetrieveWalletRequest{
			AccountId: req.AccountID,
		}
		w, err := walletClient.RetrieveWallet(context.Background(), retrieve)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				createWallet := &wallet.CreateWalletRequest{
					AccountId:       req.AccountID,
					PreviousBalance: 0,
					CurrentBalance:  0,
				}
				if _, err := walletClient.CreateWallet(context.Background(), createWallet); err != nil {
					return pkg.TransactionResponse{}, errors.Wrap(err, "unable to create wallet")
				}
			} else {
				return pkg.TransactionResponse{}, errors.Wrap(err, "unable to retrieve wallet")
			}
		}

		category := getCategory(req.Category)
		data := internal.TransactionDBRequest{
			AccountID:     req.AccountID,
			Amount:        req.Amount,
			Category:      category,
			TransactionId: req.TransactionId,
		}
		id, err := addTransaction(data)
		if err != nil {
			return pkg.TransactionResponse{}, errors.Wrap(err, "workflow - unable to insert transaction")
		}
		curBalance := w.CurrentBalance
		newBalance := curBalance + req.Amount
		if category == 0 {
			newBalance = curBalance - req.Amount
		}
		update := &wallet.UpdateWalletRequest{
			AccountId:       req.AccountID,
			PreviousBalance: curBalance,
			CurrentBalance:  newBalance,
		}
		if _, err := walletClient.UpdateWallet(context.Background(), update); err != nil {
			return pkg.TransactionResponse{}, errors.Wrap(err, "unable to update wallet")
		}
		return pkg.TransactionResponse{Id: id}, nil
	}
}

//getCategory returns category
func getCategory(category string) int {
	if strings.ToLower(category) == "debit" {
		return 0
	}
	return 1
}
