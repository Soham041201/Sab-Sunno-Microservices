package utils

import (
	"encoding/json"
)

type SocketEvent struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

func IsSocketEvent(message []byte) bool {
	var socketEvent SocketEvent
	err := json.Unmarshal(message, &socketEvent)
	if err != nil {
		return false // Not a valid SocketEvent (JSON unmarshal failed)
	}
	return true
}
