package serve

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Define request and response structs within the session package.
type SessionRequest struct {
	Username string `json:"username"`
	Data     string `json:"data"`
}

type SessionResponse struct {
	Message   string `json:"message"`
	SessionID string `json:"session_id"`
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
	if r.Method != http.MethodPost {
		fmt.Printf("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SessionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if req.Data == "" {
		http.Error(w, "Data is required", http.StatusBadRequest)
		return
	}

	sessionID := fmt.Sprintf("session-%s-%s", req.Username, req.Data)

	resp := SessionResponse{
		Message:   fmt.Sprintf("Session created for user %s", req.Username),
		SessionID: sessionID,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
}
