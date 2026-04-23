package stream

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	ID     string
	Writer http.ResponseWriter
	Done   chan struct{}
}

type Broadcaster struct {
	clients      map[string]*Client
	mu           sync.RWMutex
	chunkChan    <-chan []byte
	streamTimeout time.Duration
}

func New(chunkChan <-chan []byte, timeout time.Duration) *Broadcaster {
	return &Broadcaster{
		clients:      make(map[string]*Client),
		chunkChan:    chunkChan,
		streamTimeout: timeout,
	}
}

func (b *Broadcaster) Register(client *Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[client.ID] = client
	log.Printf("Client registered: %s (total: %d)", client.ID, len(b.clients))
}

func (b *Broadcaster) Unregister(clientID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if client, ok := b.clients[clientID]; ok {
		close(client.Done)
		delete(b.clients, clientID)
		log.Printf("Client unregistered: %s (total: %d)", clientID, len(b.clients))
	}
}

func (b *Broadcaster) Start() {
	go b.broadcastLoop()
}

func (b *Broadcaster) broadcastLoop() {
	for chunk := range b.chunkChan {
		b.mu.RLock()
		clients := make([]*Client, 0, len(b.clients))
		for _, client := range b.clients {
			clients = append(clients, client)
		}
		b.mu.RUnlock()

		for _, client := range clients {
			select {
			case <-client.Done:
				b.Unregister(client.ID)
				continue
			default:
				if !b.writeToClient(client, chunk) {
					b.Unregister(client.ID)
				}
			}
		}
	}
}

func (b *Broadcaster) writeToClient(client *Client, chunk []byte) bool {
	ctx, cancel := context.WithTimeout(context.Background(), b.streamTimeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		_, err := client.Writer.Write(chunk)
		if f, ok := client.Writer.(http.Flusher); ok {
			f.Flush()
		}
		done <- err
	}()

	select {
	case <-ctx.Done():
		log.Printf("Client %s timed out", client.ID)
		return false
	case err := <-done:
		if err != nil {
			log.Printf("Error writing to client %s: %v", client.ID, err)
			return false
		}
		return true
	}
}

func (b *Broadcaster) ClientCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.clients)
}
