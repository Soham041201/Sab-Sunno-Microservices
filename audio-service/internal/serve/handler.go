package serve

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Soham041201/Sab-Sunno-Microservices/audio-service/internal/webRTC"
	"github.com/pion/sdp/v3"
	"github.com/pion/webrtc/v3"
)

// Define request and response structs within the session package.
type HttpRequest struct {
	sdp.SessionDescription
}

// Handler struct (if you need to store state)
type Handler struct {
	// Add any necessary state here (e.g., database connection)
}

// NewHandler creates a new Handler instance.
func NewHandler() *Handler {
	return &Handler{} // Initialize any state if needed
}

// SessionHandler is the actual HTTP handler function.
func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Received request\n")

	if r.Method != http.MethodPost {
		fmt.Printf("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received SDP: %s\n", bodyBytes)

	sdpString := string(bodyBytes)
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  sdpString,
	}

	webRTC.SetupWebRTCForConnection(offer, w)
}
