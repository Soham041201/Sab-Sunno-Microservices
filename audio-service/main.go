package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type GeminiClient struct {
	conn *websocket.Conn
}

func (gc *GeminiClient) Close() error {
	return gc.conn.Close()
}

func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	u := url.URL{Scheme: "wss", Host: "generativelanguage.googleapis.com", Path: "/ws/google.ai.generativelanguage.v1alpha.GenerativeService.BidiGenerateContent", RawQuery: "key=" + apiKey}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	return &GeminiClient{conn: c}, nil
}

type GeminiResponse struct {
	ServerContent struct {
		ModelTurn struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"modelTurn"`
	} `json:"serverContent"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file") // Handle the error properly
	}
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		log.Fatal("GOOGLE_API_KEY environment variable not set")
	}

	client, err := NewGeminiClient(apiKey)
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
			// fmt.Printf("Received: %s\n", message)
			var response GeminiResponse
			err := json.Unmarshal(message, &response)
			if err != nil {
				log.Fatalf("Error unmarshaling JSON: %v", err)
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

func (gc *GeminiClient) SendSetup() error {
	setupMsg := map[string]interface{}{
		"setup": map[string]interface{}{
			"model": "models/gemini-2.0-flash-exp",
			"generation_config": map[string]interface{}{
				"response_modalities": []string{"TEXT"},
			},
		},
	}
	return gc.sendMessage(setupMsg)
}

func (gc *GeminiClient) SendTextMessage(text string) error {
	msg := map[string]interface{}{
		"client_content": map[string]interface{}{
			"turn_complete": true,
			"turns": []map[string]interface{}{
				{
					"role": "user",
					"parts": []map[string]interface{}{
						{"text": text},
					},
				},
			},
		},
	}

	return gc.sendMessage(msg)
}

func (gc *GeminiClient) sendMessage(msg map[string]interface{}) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return gc.conn.WriteMessage(websocket.TextMessage, msgJSON)
}

func (gc *GeminiClient) ReceiveMessages() (<-chan []byte, <-chan error) {
	messageChan := make(chan []byte)
	errorChan := make(chan error)

	go func() {
		defer close(messageChan)
		defer close(errorChan)
		for {
			_, message, err := gc.conn.ReadMessage()
			if err != nil {
				errorChan <- fmt.Errorf("read: %w", err)
				return
			}
			messageChan <- message
		}
	}()

	return messageChan, errorChan
}
