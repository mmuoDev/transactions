package pkg

//TransactionRequest represents data needed to insert a transaction
type TransactionRequest struct {
	AccountID     int32  `json:"accountId"`
	Amount        int32  `json:"amount"`
	Category      string `json:"category"`
	TransactionId string `json:"transactionId"`
}

//TransactionResponse represents transaction insertion response
type TransactionResponse struct {
	Id int64 `json:"id"`
}
