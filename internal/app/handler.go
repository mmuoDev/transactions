package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mmuoDev/transactions/internal/db"
	"github.com/mmuoDev/transactions/internal/workflow"
	"github.com/mmuoDev/transactions/pkg"
)

//InsertTransactionHandler returns a http request to add a transaction
func InsertTransactionHandler(addTransaction db.InsertTransactionFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pkg.TransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		category := strings.ToLower(req.Category)
		if category != "credit" && category != "debit" {
			w.Write([]byte(fmt.Sprintf("Invalid transaction category=%s", category)))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		insert := workflow.InsertTransaction(addTransaction)
		if err := insert(req); err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
