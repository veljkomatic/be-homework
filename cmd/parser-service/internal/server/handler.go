package server

import (
	"encoding/json"
	"github.com/veljkomatic/be-homework/pkg/blockchain"
	"net/http"
	"strings"
)

type httpHandler func(w http.ResponseWriter, r *http.Request)

type GetCurrentBlockNumberResponse struct {
	BlockNumber int `json:"block_number"`
}

func GetCurrentBlockNumberHandler(service Service) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		currentBlockNumber := service.GetCurrentBlockNumber(r.Context())
		resp := GetCurrentBlockNumberResponse{
			BlockNumber: currentBlockNumber,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
}

type SubscribeResponse struct {
	Subscribed bool `json:"subscribed"`
}

type SubscribeBody struct {
	Address string `json:"address"`
}

func SubscribeHandler(service Service) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		var body SubscribeBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		subscribed := service.Subscribe(r.Context(), body.Address)
		resp := SubscribeResponse{
			Subscribed: subscribed,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
}

// GetTransactionsResponse TODO in future do not return blockchain.Transaction rather some DTO
type GetTransactionsResponse struct {
	Transactions []*blockchain.Transaction `json:"transactions"`
}

func GetTransactionsHandler(service Service) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 || parts[1] != "transactions" {
			http.NotFound(w, r)
			return
		}

		address := parts[2]
		transactions := service.GetTransactions(r.Context(), address)

		resp := GetTransactionsResponse{
			Transactions: transactions,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
}
