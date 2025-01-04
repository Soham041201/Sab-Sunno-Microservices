package gemini

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type GeminiResponse struct {
	ServerContent struct {
		ModelTurn struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"modelTurn"`
	} `json:"serverContent"`
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
