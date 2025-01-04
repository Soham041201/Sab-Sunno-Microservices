package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Soham041201/Sab-Sunno-Microservices/audio-service/internal/gemini" // Correct import path
)

func main() {

	client, err := gemini.SetupGeminiServiceToGenerateResponsee() // Correct function name
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	err = client.SendTextMessage("How are you?") // Correct function name
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
