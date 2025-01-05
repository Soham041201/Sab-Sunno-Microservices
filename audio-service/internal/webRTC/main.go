package webRTC

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Soham041201/Sab-Sunno-Microservices/audio-service/internal/gemini"
	"github.com/Soham041201/Sab-Sunno-Microservices/audio-service/utils"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

func SetupWebRTCForConnection(socketEvent utils.SocketEvent, c *websocket.Conn) error {
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
		return fmt.Errorf("creating peer connection: %w", err)
	}
	defer peerConnection.Close() // Close connection on exit

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
		fmt.Println("ICE Candidate:", string(candidateJSON))
		socketEvent := utils.SocketEvent{
			Event: "ice-candidate",
			Data:  candidateJSON,
		}
		c.WriteJSON(socketEvent)
	})

	// Handle data channel messages (same as before)
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open.\n", d.Label(), d.ID())

		})

		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
			gemini.HandleGeminiResponse(string(msg.Data), d)
		})
	})

	// Handle offer and answer
	switch socketEvent.Event {
	case "offer":
		var offer webrtc.SessionDescription
		err := json.Unmarshal(socketEvent.Data, &offer)
		if err != nil {
			return fmt.Errorf("unmarshaling offer: %w", err)
		}

		err = peerConnection.SetRemoteDescription(offer)
		if err != nil {
			return fmt.Errorf("setting remote description: %w", err)
		}

		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			return fmt.Errorf("creating answer: %w", err)
		}

		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			return fmt.Errorf("setting local description: %w", err)
		}

		answerBytes, err := json.Marshal(answer)
		if err != nil {
			return fmt.Errorf("marshaling answer: %w", err)
		}
		socketEvent := utils.SocketEvent{
			Event: "answer",
			Data:  answerBytes,
		}
		c.WriteJSON(socketEvent)

	case "ice-candidate":
		fmt.Println("Received ICE Candidate from React", string(socketEvent.Data))
		var candidate webrtc.ICECandidateInit
		err := json.Unmarshal(socketEvent.Data, &candidate)
		if err != nil {
			return fmt.Errorf("unmarshaling ice candidate: %w", err)
		}
		fmt.Println("Unmarshaled ICE Candidate", candidate)
		err = peerConnection.AddICECandidate(candidate)
		if err != nil {
			return fmt.Errorf("adding ice candidate: %w", err)
		}

	default:
		return fmt.Errorf("invalid event type: %s", socketEvent.Event)
	}

	select {}
}
