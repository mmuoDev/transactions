package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mmuoDev/core-proto/gen/wallet"
	"github.com/mmuoDev/transactions/internal/db"
	"github.com/mmuoDev/transactions/internal/workflow"
	"github.com/mmuoDev/transactions/pkg"
)

//ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

//InsertTransactionHandler returns a http request to add a transaction
func InsertTransactionHandler(addTransaction db.InsertTransactionFunc, walletClient wallet.WalletClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pkg.TransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			res := ErrorResponse{Error: err.Error()}
			ServeJSON(res, w)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		category := strings.ToLower(req.Category)
		if category != "credit" && category != "debit" {
			res := ErrorResponse{Error: fmt.Sprintf("Invalid transaction category=%s", category)}
			ServeJSON(res, w)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		insert := workflow.InsertTransaction(addTransaction, walletClient)
		res, err := insert(req)
		if err != nil {
			res := ErrorResponse{Error: err.Error()}
			ServeJSON(res, w)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ServeJSON(res, w)
		w.WriteHeader(http.StatusOK)
	}
}

func ServeJSON(res interface{}, w http.ResponseWriter) {
	bb, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(bb)
}
