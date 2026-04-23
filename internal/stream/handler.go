package stream

import (
	"net/http"
	"time"
)

type Handler struct {
	broadcaster *Broadcaster
}

func NewHandler(broadcaster *Broadcaster) *Handler {
	return &Handler{broadcaster: broadcaster}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set headers for streaming
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create client
	clientID := generateClientID()
	client := &Client{
		ID:     clientID,
		Writer: w,
		Done:   make(chan struct{}),
	}

	// Register client
	h.broadcaster.Register(client)
	defer h.broadcaster.Unregister(clientID)

	// Keep connection open until client disconnects
	<-client.Done
}

func generateClientID() string {
	return time.Now().Format("20060102150405.000000")
}
