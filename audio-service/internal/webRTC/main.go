package webRTC

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Soham041201/Sab-Sunno-Microservices/audio-service/utils"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

func SetupWebRTCForConnection(clientOffer webrtc.SessionDescription, c *websocket.Conn) {
	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Fatal(err)
	}
	defer peerConnection.Close()

	// Handle ICE candidates
	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		candidateJSON, err := json.Marshal(candidate.ToJSON())

		if err != nil {
			log.Println("error marshaling ice candidate", err)
			return
		}
		fmt.Println("ICE Candidate:", string(candidateJSON)) // Send this to the client
		socketEvent := utils.SocketEvent{
			Event: "ice-candidate",
			Data:  candidateJSON,
		}
		c.WriteJSON(socketEvent)
	})

	// Handle data channel messages
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open.\n", d.Label(), d.ID())
			d.SendText("Hello, Client!")
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})

	// Receive offer from client (this is a placeholder for your signaling mechanism)

	// Set remote description
	err = peerConnection.SetRemoteDescription(clientOffer)
	if err != nil {
		log.Fatal(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Sets the LocalDescription of the local peer
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Answer SDP:\n", answer.SDP) // Send this answer to the client

	c.WriteJSON(answer)

	// Block forever
	// select {}
}
