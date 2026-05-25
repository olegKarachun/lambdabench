package models

type Transaction struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}

type LambdaPayload struct {
	RequestID    string        `json:"requestId"`
	Transactions []Transaction `json:"transactions"`
}
