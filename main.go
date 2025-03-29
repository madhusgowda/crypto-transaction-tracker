package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
	Build a system that tracks cryptocurrency transactions for a given addresses across multiple currencies.
	Each transaction contains details like transactionID, from, to, value, fee, and timestamp.
	Initial Balances (current_state.json): This dataset provides the current balance for several addresses across Bitcoin, Dogecoin, and TON.
	Transactions (transactions.json): This dataset contains 1,000 transactions for different addresses in the system.
	1. getCurrentBalance(address,coin) => output: summation balance(if multiple occurrence)
	2. getMaxTransaction(coin) => output: transactionID, from, to, value, fee, and timestamp

	https://www.jsonkeeper.com/b/K4PX
	https://www.jsonkeeper.com/b/ILH9
*/

func main() {
	balanceUrl := "https://www.jsonkeeper.com/b/ILH9"
	transactionUrl := "https://www.jsonkeeper.com/b/K4PX"

	var response CommonResponse

	http.HandleFunc("/balances", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		coin := strings.ToLower(r.URL.Query().Get("coin"))

		if address == "" || coin == "" {
			http.Error(w, "address or coin is required", http.StatusBadRequest)
			return
		}

		var remainingBalance float64

		var balanceResp BalanceResponse
		if err := fetchJSONData(balanceUrl, &balanceResp); err != nil {
			http.Error(w, "Failed to fetch balance data", http.StatusInternalServerError)
			return
		}

		for _, balanceData := range balanceResp.InitialState {
			if balanceData.Address == address && strings.ToLower(balanceData.Coin) == coin {
				stringTrimmedBalance := strings.Fields(balanceData.Balance)
				floatBalance, _ := strconv.ParseFloat(stringTrimmedBalance[0], 64)
				remainingBalance += floatBalance
			}
		}

		var transactionResoponse TransactionResponse
		if err := fetchJSONData(transactionUrl, &transactionResoponse); err != nil {
			http.Error(w, "Failed to fetch transaction data", http.StatusInternalServerError)
			return
		}

		var maxTransaction TransactionData
		maxTransactionFound := false
		maxTxValue := 0.0

		for _, transaction := range transactionResoponse.Transactions {
			if strings.ToLower(transaction.Coin) == coin {
				if transaction.Value > maxTxValue {
					maxTxValue = transaction.Value
					maxTransactionFound = true
					maxTransaction.TransactionID = transaction.TransactionID
					maxTransaction.From = transaction.From
					maxTransaction.To = transaction.To
					maxTransaction.Coin = strings.ToLower(transaction.Coin)
					maxTransaction.Value = transaction.Value
					maxTransaction.Timestamp = transaction.Timestamp
				}
			}
		}

		if maxTransactionFound {
			response.MaxTransaction = &maxTransaction
		}

		response = CommonResponse{
			Message:        "Data fetched successfully",
			Balance:        remainingBalance,
			MaxTransaction: &maxTransaction,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	port := "8080"
	log.Printf("Server running on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func fetchJSONData(url string, responseStruct interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching data from the url: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	return json.Unmarshal(body, responseStruct)
}
