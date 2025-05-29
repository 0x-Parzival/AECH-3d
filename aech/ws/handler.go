package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// upgrader is a Gorilla WebSocket upgrader instance with basic configuration.
// CheckOrigin allows all origins for development purposes.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by returning true.
		// For production, you should implement proper origin checking.
		// Example: return r.Header.Get("Origin") == "http://localhost:yourfrontendport"
		return true
	},
}

// HandleConnections returns a Gin handler function that manages WebSocket connections.
// It takes a Hub instance to register clients and handle message broadcasting.
func HandleConnections(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Upgrade HTTP server connection to a WebSocket connection
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade WebSocket connection: %v", err)
			// Note: upgrader.Upgrade writes a response on error, so we don't need c.JSON here.
			return
		}
		// Ensure connection is closed when the handler exits
		// defer conn.Close() // The client's writePump and readPump should handle closing.

		// Create a new client channel for this connection
		// Buffer size of 1 for the client channel means the hub won't block
		// if the client's writePump is slightly delayed, but if it's consistently slow,
		// the hub's Broadcast will detect it via the non-blocking send.
		client := make(Client, 1) // Buffered channel for the client
		hub.Register(client)
		log.Printf("WebSocket client connected: %s", conn.RemoteAddr().String())

		// Goroutine to write messages from the hub to the WebSocket client
		go writePump(conn, client, hub)
		// Goroutine to read messages from the WebSocket client (and keep connection alive)
		go readPump(conn, hub, client) // Pass client to unregister it if read fails
	}
}

// writePump pumps messages from the client's channel to the WebSocket connection.
// A Vgoroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this Vgoroutine.
func writePump(conn *websocket.Conn, client Client, hub *Hub) {
	defer func() {
		// When this goroutine exits (e.g., client channel closed by hub, or write error),
		// ensure the WebSocket connection is closed.
		log.Printf("Closing WebSocket connection (writePump exiting): %s", conn.RemoteAddr().String())
		conn.Close()
		// Unregister client is handled by hub or readPump for robustness,
		// but can be done here too if client channel closure always means unregister.
		// hub.Unregister(client) // Usually handled by hub's broadcast or readPump on error
	}()

	for message := range client { // Loop until client channel is closed
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error writing message to WebSocket client %s: %v", conn.RemoteAddr().String(), err)
			// If WriteMessage fails, the connection is likely broken.
			// The client channel might still be open if the hub hasn't detected it yet.
			// Closing the connection here will cause readPump to exit, which can trigger unregistration.
			return // Exit goroutine, defer will close conn.
		}
		log.Printf("Sent message to WebSocket client %s: %s", conn.RemoteAddr().String(), string(message))
	}
	// If client channel is closed by the hub (e.g. during broadcast to a slow client), this loop terminates.
	log.Printf("Client channel closed, writePump for %s exiting.", conn.RemoteAddr().String())
}

// readPump pumps messages from the WebSocket connection to the hub (if needed for bi-directional).
// Currently, it's mainly used to detect client disconnection.
// A Vgoroutine running readPump is started for each connection.
func readPump(conn *websocket.Conn, hub *Hub, client Client) {
	defer func() {
		log.Printf("Unregistering client and closing WebSocket connection (readPump exiting): %s", conn.RemoteAddr().String())
		hub.Unregister(client) // Unregister the client from the hub
		conn.Close()           // Close the WebSocket connection
	}()

	// Configure pong handler, max message size, etc. (optional)
	// conn.SetReadLimit(maxMessageSize)
	// conn.SetReadDeadline(time.Now().Add(pongWait))
	// conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Read messages from the WebSocket connection.
	// This loop is primarily to detect if the client has disconnected.
	for {
		// ReadMessage is a blocking call.
		// If the client disconnects, it will return an error, breaking the loop.
		if _, message, err := conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket client %s disconnected (unexpected error): %v", conn.RemoteAddr().String(), err)
			} else {
				log.Printf("WebSocket client %s disconnected (error reading message): %v", conn.RemoteAddr().String(), err)
			}
			break // Exit loop on any error (typically client disconnection)
		}
		// Messages read from clients are currently ignored.
		// If you need to process messages from clients, add logic here.
		 log.Printf("Received message from client %s (currently ignored): %s", conn.RemoteAddr().String(), string(message))
	}
}
