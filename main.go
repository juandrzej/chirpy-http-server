package main

import (
	"log"
	"net/http"
)

func main() {
	serverHandler := http.NewServeMux()
	server := http.Server{
		Handler: serverHandler,
		Addr:    ":8080",
	}

	serverHandler.Handle("/", http.FileServer(http.Dir(".")))

	log.Fatal(server.ListenAndServe())
}
