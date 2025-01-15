package gemini

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/pion/webrtc/v3"
)

func SetupGeminiServiceToGenerateResponsee() (*GeminiClient, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		log.Fatal("GOOGLE_API_KEY environment variable not set")
	}
	client, err := NewGeminiClient(apiKey)

	if err != nil {
		log.Fatal(err)
	}

	err = client.SendSetup()

	if err != nil {
		log.Fatal("setup:", err)
	}

	return client, nil

}

func HandleGeminiResponse(data []byte, sampleRate int, d *webrtc.PeerConnection) error {

	client, err := SetupGeminiServiceToGenerateResponsee() // Correct function name
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	err = client.SendAudioMessage(data, sampleRate) // Correct function name
	if err != nil {
		log.Fatal("text message:", err)
	}

	messageChan, errorChan := client.ReceiveMessages()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case message := <-messageChan:
			var response GeminiResponse
			err := json.Unmarshal(message, &response)
			if err != nil {
				log.Fatalf("Error unmarshaling JSON: %v, Raw message: %s", err, message)
			}
			fmt.Print("response", response)
			if len(response.ServerContent.ModelTurn.Parts) > 0 {
				fmt.Println(response.ServerContent)
				// d.SendText(response.ServerContent.ModelTurn.Parts[0].Text)
			} else {
				fmt.Println("No text parts found in the response.")
			}
		case err := <-errorChan:
			log.Println("receive error:", err)
			return fmt.Errorf("error from gemini %d", err)
		case <-interrupt:
			fmt.Println("interrupt")
			return fmt.Errorf("error from gemini %d", err)
		}
	}
}
