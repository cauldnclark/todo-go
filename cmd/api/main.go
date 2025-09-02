package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message": "Hello, World! Server is running on port 3344"}`)
	})

	fmt.Println("Server starting on port 3344...")
	fmt.Println("Listening on http://localhost:3344")
	fmt.Println("Press Ctrl+C to stop the server")

	// Start the server and keep it running
	if err := http.ListenAndServe(":3344", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
