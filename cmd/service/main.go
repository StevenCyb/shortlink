package main

import (
	"log"
	"shortlink/pkg/api"
	"shortlink/pkg/store"
	"time"
)

func main() {
	// this attributes could be provided via arg or env (e.g. launch json)
	listen := "localhost:8080"
	expiration := 1 * time.Hour
	gcInterval := 1 * time.Minute

	store := store.NewInMemory(expiration, gcInterval)
	server := api.NewServer(store)

	err := server.ListenAndServe(listen)
	if err != nil {
		log.Fatal(err)
	}
}
