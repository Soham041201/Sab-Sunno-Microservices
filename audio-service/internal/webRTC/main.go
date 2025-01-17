package webRTC

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Soham041201/Sab-Sunno-Microservices/audio-service/utils"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type WebRtcPeerConnection struct {
	c              *websocket.Conn
	done           chan struct{}
	peerConnection *webrtc.PeerConnection
}

func NewWebRtcSocket(c *websocket.Conn, done chan struct{}, pc *webrtc.PeerConnection) *WebRtcPeerConnection {
	return &WebRtcPeerConnection{
		c:              c,
		done:           done,
		peerConnection: pc,
	}
}

func SetupWebRTCForConnection(socketEvent utils.SocketEvent, c *websocket.Conn) *webrtc.PeerConnection {

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)

	if err != nil {
		fmt.Print("setting remote description: %w", err.Error())

	}

	return peerConnection

}

func (w *WebRtcPeerConnection) HandleIceCandidate(candidate *webrtc.ICECandidate) {
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

func (w *WebRtcPeerConnection) HandleTrack(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
	fmt.Println("Track received")
	HandleTrack(track, w.peerConnection)
}

func (w *WebRtcPeerConnection) HandleDataChannel(d *webrtc.DataChannel) {
	fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

	d.OnOpen(func() {
		fmt.Printf("Data channel '%s'-'%d' open.\n", d.Label(), d.ID())

	})

	d.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		// gemini.HandleGeminiResponse(string(msg.Data), d)
	})
}

func (w *WebRtcPeerConnection) HandleConnectioChange(state webrtc.PeerConnectionState) {
	fmt.Println("Connection State has changed", state.String())
	if state == webrtc.PeerConnectionStateDisconnected {
		fmt.Println("Peer connection is broken. Exiting.", state.String())
		w.peerConnection.Close()

	}

	if state == webrtc.PeerConnectionStateConnected {
		fmt.Println("Peer connection is now connected")
	}
}

func (w *WebRtcPeerConnection) HandlePeerConnectionOffer(offer webrtc.SessionDescription) {

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

func (w *WebRtcPeerConnection) HandleIceCandidateSocketEvent(data []byte) {
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
