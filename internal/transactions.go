package internal

//TransactionDBRequest represents a transaction insert DB request
type TransactionDBRequest struct {
	AccountID     int32  `json:"accountId"`
	Amount        int32  `json:"amount"`
	Category      int    `json:"category"`
	TransactionId string `json:"transactionId"`
}
