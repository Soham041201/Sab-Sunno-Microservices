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
	"github.com/pion/webrtc/v3"
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
					fmt.Println("Error unmarshaling SocketEvent:", err)
					return // Or take other appropriate action
				}

				peerConnection := webRTC.SetupWebRTCForConnection(data, c)
				pc := webRTC.NewWebRtcSocket(c, done, peerConnection)
				defer peerConnection.Close()

				switch data.Event {
				case "offer":
					var offer webrtc.SessionDescription
					err := json.Unmarshal(data.Data, &offer)
					if err != nil {
						fmt.Print("setting remote description: %w", err.Error())
					}
					pc.HandlePeerConnectionOffer(offer)
				case "ice-candidate":
					pc.HandleIceCandidateSocketEvent(data.Data)
				}

				peerConnection.OnICECandidate(pc.HandleIceCandidate)
				peerConnection.OnTrack(pc.HandleTrack)
				peerConnection.OnDataChannel(pc.HandleDataChannel)
				peerConnection.OnConnectionStateChange(pc.HandleConnectioChange)
			} else {
				fmt.Println("Invalid SocketEvent received", string(message))
			}
		}
	}()

	select {
	case <-done:
	case <-interrupt:
	}

}
