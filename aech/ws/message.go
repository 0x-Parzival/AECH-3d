package ws

// WSMessage defines the structure for messages sent over WebSocket.
type WSMessage struct {
	Type string      `json:"type"` // e.g., "new_block", "new_tx"
	Data interface{} `json:"data"` // The actual data payload (e.g., a block.Block3D or block.Transaction)
}
