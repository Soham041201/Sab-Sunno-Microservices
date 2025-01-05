package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	// Import from your internal package
	"github.com/Soham041201/Sab-Sunno-Microservices/audio-service/internal/serve" // Correct import path
	// Correct import path
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	handler := serve.NewHandler() // Initialize your handler from the internal package

	http.HandleFunc("/", handler.HandleRequest)

	fmt.Printf("Server listening on port %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))

}
