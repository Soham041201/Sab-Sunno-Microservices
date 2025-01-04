package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Soham041201/Sab-Sunno-Microservices/audio-service/internal/gemini" // Correct import path
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		log.Fatal("GOOGLE_API_KEY environment variable not set")
	}

	client, err := gemini.NewGeminiClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	err = client.SendSetup()
	if err != nil {
		log.Fatal("setup:", err)
	}

	err = client.SendTextMessage("Hello Gemini from Go!")
	if err != nil {
		log.Fatal("text message:", err)
	}

	messageChan, errorChan := client.ReceiveMessages()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case message := <-messageChan:
			var response gemini.GeminiResponse
			err := json.Unmarshal(message, &response)
			if err != nil {
				log.Fatalf("Error unmarshaling JSON: %v, Raw message: %s", err, message)
			}

			if len(response.ServerContent.ModelTurn.Parts) > 0 {
				fmt.Println(response.ServerContent.ModelTurn.Parts[0].Text)
			} else {
				fmt.Println("No text parts found in the response.")
			}
		case err := <-errorChan:
			log.Println("receive error:", err)
			return
		case <-interrupt:
			fmt.Println("interrupt")
			return
		}
	}
}
