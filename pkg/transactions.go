package pkg

//TransactionRequest represents data needed to record a transaction
type TransactionRequest struct {
	AccountID int32  `json:"accountId"`
	Amount    int32  `json:"amount"`
	Category  string `json:"category"`
}
