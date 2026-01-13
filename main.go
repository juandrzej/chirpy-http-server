package main

import "net/http"

func main() {
	serverHandler := http.NewServeMux()
	server := http.Server{
		Handler: serverHandler,
		Addr:    ":8080",
	}

	_ = server.ListenAndServe()
}
