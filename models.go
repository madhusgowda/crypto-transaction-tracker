package main

type TransactionResponse struct {
	Transactions []struct {
		TransactionID string  `json:"transactionID"`
		From          string  `json:"from"`
		To            string  `json:"to"`
		Coin          string  `json:"coin"`
		Value         float64 `json:"value"`
		Fee           float64 `json:"fee"`
		Timestamp     string  `json:"timestamp"`
	} `json:"transactions"`
}

type BalanceResponse struct {
	InitialState []struct {
		Address string `json:"address"`
		Coin    string `json:"coin"`
		Balance string `json:"balance"`
	} `json:"initial_state"`
}

type CommonResponse struct {
	Message        string           `json:"message"`
	Balance        float64          `json:"balance"`
	MaxTransaction *TransactionData `json:"maxTransaction,omitempty"`
}

type TransactionData struct {
	TransactionID string  `json:"transactionID"`
	From          string  `json:"from"`
	To            string  `json:"to"`
	Coin          string  `json:"coin"`
	Value         float64 `json:"value"`
	Timestamp     string  `json:"timestamp"`
}
