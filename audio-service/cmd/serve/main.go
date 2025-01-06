package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/Soham041201/Sab-Sunno-Microservices/audio-service/internal/webRTC"
	"github.com/Soham041201/Sab-Sunno-Microservices/audio-service/utils"
	"github.com/gorilla/websocket"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	addr := flag.String("addr", "localhost:8001", "http service address")
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()

			if err != nil {
				log.Println("read:", err)
				return
			}

			if utils.IsSocketEvent(message) {
				var data utils.SocketEvent
				err := json.Unmarshal(message, &data)
				if err != nil {
					// Handle error gracefully (e.g., log the error, send an error message to the client)
					fmt.Println("Error unmarshaling SocketEvent:", err)
					return // Or take other appropriate action
				}
				// webRTC.SetupWebRTCForConnection(data, c)
				if string(data.Event) == "offer" || string(data.Event) == "ice-candidate" {
					webRTC.SetupWebRTCForConnection(data, c)
					
				} else {
					fmt.Println("Unhandled SocketEvent:", data.Event)
				}
			} else {
				// Handle messages that are not valid SocketEvents
				fmt.Println("Invalid SocketEvent received", string(message))
			}
		}
	}()

	select {
	case <-done:
	case <-interrupt:
	}

}
