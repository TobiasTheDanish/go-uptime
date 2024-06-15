package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		f := rand.Float32()

		if f >= 0.9 {
			fmt.Println("Server error")
			w.WriteHeader(500)
			w.Write([]byte("Server error"))
		} else {
			fmt.Println("Server succes")
			w.Write([]byte("Server running"))
		}
	})

	server := &http.Server{
		Addr:    ":4000",
		Handler: mux,
	}

	fmt.Println("Server listening on port 4000")
	server.ListenAndServe()
}
