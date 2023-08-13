package server

import (
	"log"
	"net/http"
)

func StartServer(service Service, port string) {
	http.HandleFunc("/block-number", GetCurrentBlockNumberHandler(service))
	http.HandleFunc("/subscribe", SubscribeHandler(service))
	http.HandleFunc("/transactions/", GetTransactionsHandler(service))

	log.Printf("Server started on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
