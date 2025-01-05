package utils

import (
	"encoding/json"
)

type SocketEvent struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func IsSocketEvent(message []byte) bool {
	var socketEvent SocketEvent
	err := json.Unmarshal(message, &socketEvent)
	return err == nil
}
