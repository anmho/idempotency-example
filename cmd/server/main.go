package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

var idempotencyKeyStore = map[string]struct{}{}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /chargeCustomer", func(w http.ResponseWriter, r *http.Request) {
		idempotencyKey := r.Header.Get("x-idempotency-key")
		if idempotencyKey != "" {
			// check if we already stored the idempotencyKey
			// if yes then abort the request early as a success
			if _, ok := idempotencyKeyStore[idempotencyKey]; ok {
				log.Println("already handled or am handling this payment request")
				w.WriteHeader(http.StatusOK)
				return
			}

			// set the idempotency key if it is set
			idempotencyKeyStore[idempotencyKey] = struct{}{}
		} else {
			log.Println("no idempotency key supplied")
		}

		log.Println("handling this payment")

	})

	srv := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	fmt.Println("starting server")
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("error running server %s", err)
		}
	}
}
