package webRTC

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Soham041201/Sab-Sunno-Microservices/audio-service/utils"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type WebRtcSocket struct {
	c              *websocket.Conn
	done           chan struct{}
	peerConnection *webrtc.PeerConnection
}

func NewWebRtcSocket(c *websocket.Conn, done chan struct{}, pc *webrtc.PeerConnection) *WebRtcSocket {
	return &WebRtcSocket{
		c:              c,
		done:           done,
		peerConnection: pc,
	}
}

func SetupWebRTCForConnection(socketEvent utils.SocketEvent, c *websocket.Conn) {
	// Prepare the configuration
	done := make(chan struct{})

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	rtcSocketInstace := NewWebRtcSocket(c, done, peerConnection)

	if err != nil {
		fmt.Print("setting remote description: %w", err.Error())

	}
	defer peerConnection.Close() // Close connection on exit

	// Handle ICE candidates
	peerConnection.OnICECandidate(rtcSocketInstace.HandleIceCandidate)
	peerConnection.OnTrack(rtcSocketInstace.HandleTrack)
	peerConnection.OnDataChannel(rtcSocketInstace.HandleDataChannel)
	peerConnection.OnConnectionStateChange(rtcSocketInstace.HandleConnectioChange)

	// Handle offer and answer
	switch socketEvent.Event {
	case "offer":
		var offer webrtc.SessionDescription
		err := json.Unmarshal(socketEvent.Data, &offer)
		if err != nil {
			fmt.Print("setting remote description: %w", err.Error())

		}
		rtcSocketInstace.HandlePeerConnectionOffer(offer)
	case "ice-candidate":
		rtcSocketInstace.HandleIceCandidateSocketEvent(socketEvent.Data)
	default:
		fmt.Print("setting remote description: %w", err.Error())
	}

	<-done // Wait for connection to close
}

func (w *WebRtcSocket) HandleIceCandidate(candidate *webrtc.ICECandidate) {
	if candidate == nil {
		return
	}

	candidateJSON, err := json.Marshal(candidate.ToJSON())
	if err != nil {
		log.Println("error marshaling ice candidate", err)
	}
	fmt.Println("ICE Candidate:", string(candidateJSON))
	socketEvent := utils.SocketEvent{
		Event: "ice-candidate",
		Data:  candidateJSON,
	}
	w.c.WriteJSON(socketEvent)
}

func (w *WebRtcSocket) HandleTrack(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
	fmt.Println("Track received")
}

func (w *WebRtcSocket) HandleDataChannel(d *webrtc.DataChannel) {
	fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

	d.OnOpen(func() {
		fmt.Printf("Data channel '%s'-'%d' open.\n", d.Label(), d.ID())

	})

	d.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		// gemini.HandleGeminiResponse(string(msg.Data), d)
	})
}

func (w *WebRtcSocket) HandleConnectioChange(state webrtc.PeerConnectionState) {
	fmt.Println("Connection State has changed", state.String())
	if state == webrtc.PeerConnectionStateDisconnected {
		fmt.Println("Peer connection is broken. Exiting.", state.String())
		close(w.done) // Signal completion

	}

	if state == webrtc.PeerConnectionStateConnected {
		fmt.Println("Peer connection is now connected")
	}
}

func (w *WebRtcSocket) HandlePeerConnectionOffer(offer webrtc.SessionDescription) {

	err := w.peerConnection.SetRemoteDescription(offer)
	if err != nil {
		fmt.Print("setting remote description: %w", err.Error())
	}

	answer, err := w.peerConnection.CreateAnswer(nil)
	if err != nil {
		fmt.Print("setting remote description: %w", err.Error())
	}

	err = w.peerConnection.SetLocalDescription(answer)
	if err != nil {
		fmt.Print("setting remote description: %w", err.Error())

	}
	answerBytes, err := json.Marshal(answer)
	if err != nil {
		fmt.Print("setting remote description: %w", err.Error())
	}
	socketEvent := utils.SocketEvent{
		Event: "answer",
		Data:  answerBytes,
	}
	w.c.WriteJSON(socketEvent)
}

func (w *WebRtcSocket) HandleIceCandidateSocketEvent(data []byte) {
	fmt.Println("Received ICE Candidate from React", string(data))
	var candidate webrtc.ICECandidateInit
	err := json.Unmarshal(data, &candidate)
	if err != nil {
		fmt.Print("setting remote description: %w", err.Error())

	}
	fmt.Println("Unmarshaled ICE Candidate", candidate)
	err = w.peerConnection.AddICECandidate(candidate)
	if err != nil {
		fmt.Print("setting remote description: %w", err.Error())
	}
}
