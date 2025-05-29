package ws

import (
	"log"
	"sync"
)

// Client represents a single WebSocket client by its message channel.
type Client chan []byte // Messages are byte slices.

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	mu      sync.Mutex
	clients map[Client]struct{} // Using struct{} as value for set-like behavior
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[Client]struct{}),
	}
}

// Register adds a new client to the hub.
func (h *Hub) Register(c Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = struct{}{}
	log.Printf("Hub: Client registered. Total clients: %d", len(h.clients))
}

// Unregister removes a client from the hub and closes its channel.
func (h *Hub) Unregister(c Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		// It's crucial to close the channel to signal the client's write goroutine to stop.
		// The write goroutine (in handler.go) should range over this channel.
		close(c) 
		log.Printf("Hub: Client unregistered. Total clients: %d", len(h.clients))
	}
}

// Broadcast sends a message to all registered clients.
// If a client's channel is blocked (e.g., slow consumer), that client is unregistered.
func (h *Hub) Broadcast(msg []byte) {
	h.mu.Lock()
	// Create a slice of clients to iterate over to avoid issues if Unregister is called
	// by a client goroutine during this broadcast (which could modify h.clients).
	// However, with the current lock on the whole Broadcast, direct iteration is fine.
	// For more complex scenarios, consider copying keys or using channels for broadcasting.
	currentClients := make([]Client, 0, len(h.clients))
    for c := range h.clients {
        currentClients = append(currentClients, c)
    }
	h.mu.Unlock() // Unlock while sending to avoid blocking other hub operations for too long.

	log.Printf("Hub: Broadcasting message to %d client(s): %s", len(currentClients), string(msg))
	
	unregisterList := []Client{}

	for _, c := range currentClients {
		select {
		case c <- msg: // Attempt to send the message
		default:
			// Channel is blocked (buffer full or receiver goroutine is stuck/gone).
			// Mark for unregistration.
			log.Printf("Hub: Client channel blocked or full. Marking for unregistration.")
			unregisterList = append(unregisterList, c)
		}
	}

    // Unregister clients that were marked
    if len(unregisterList) > 0 {
        h.mu.Lock()
        for _, c := range unregisterList {
            if _, ok := h.clients[c]; ok {
                delete(h.clients, c)
                close(c) // Close channel after removing from map
                log.Printf("Hub: Client auto-unregistered due to blocked channel. Total clients: %d", len(h.clients))
            }
        }
        h.mu.Unlock()
    }
}
