package api

import (
	"aech/block"
	"aech/plane"
	"aech/txpool"
	"aech/ws"      // Added for WebSocket integration
	"github.com/gin-gonic/gin"
	"encoding/hex"  // Added for ID calculation
	"encoding/json" // Added for WebSocket message marshaling
	"fmt"           // Added for error formatting
	"log"
	"math"    // Added for pagination (math.Ceil)
	"net/http"
	"regexp"  // Added for name validation
	"strconv" // Added for pagination (Atoi)
	"strings" // Added for Content-Type check
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

// APIServer holds instances needed by the API handlers.
type APIServer struct {
	plane  *plane.Plane3D
	txpool *txpool.TxPool
	wsHub  *ws.Hub // Added WebSocket hub
}

// AddBlockRequest defines the expected request body for adding a new block.
type AddBlockRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
	// PreviousHash is intentionally omitted; server will determine it.
}

// RunServer initializes and starts the Gin HTTP server.
// It sets up routes for block and transaction interactions.
func RunServer(port string, p *plane.Plane3D, tp *txpool.TxPool) {
	hub := ws.NewHub() // Initialize the WebSocket hub

	server := &APIServer{
		plane:  p,
		txpool: tp,
		wsHub:  hub, // Pass the hub instance
	}

	r := gin.Default()

	// Register global error handler first
	r.Use(CORSMiddleware()) // Add CORS middleware
	r.Use(globalErrorHandler())

	// Middleware for Request Body Size Limit
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1024*1024) // 1 MB limit
		c.Next()
	})

	r.Use(RequireJSONContentType()) // Add Content-Type check middleware

	// Block routes
	r.GET("/blocks", server.getBlocksHandler)
	r.GET("/block/:id", server.getBlockByIDHandler)
	r.POST("/block/add", server.addBlockHandler) // If needed

	// Transaction routes
	r.POST("/transaction/add", server.addTransactionHandler)
	r.GET("/txpool", server.getTxPoolHandler)

	// Plane routes
	r.GET("/plane/grid", server.getPlaneGridHandler)

	// WebSocket route
	r.GET("/ws", server.handleWebSocketConnections)

	// TODO: Add other endpoints here as they are developed

	log.Printf("Starting API server on port %s...", port) // Added "..." for clarity
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err) // This log is already good
	}
}

// getBlocksHandler handles GET requests to /blocks.
// It retrieves all blocks from the plane and returns them as a paginated JSON response.
func (s *APIServer) getBlocksHandler(c *gin.Context) {
	log.Printf("Incoming %s %s request (with query: %s)", c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery)

	// Parse page query parameter
	pageStr := c.DefaultQuery("page", strconv.Itoa(DefaultPage))
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		log.Printf("Invalid page parameter '%s', defaulting to %d. Error: %v", pageStr, DefaultPage, err)
		page = DefaultPage
	}

	// Parse limit query parameter
	limitStr := c.DefaultQuery("limit", strconv.Itoa(DefaultLimit))
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		log.Printf("Invalid limit parameter '%s', defaulting to %d. Error: %v", limitStr, DefaultLimit, err)
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		log.Printf("Limit parameter %d exceeds MaxLimit %d, capping at MaxLimit.", limit, MaxLimit)
		limit = MaxLimit
	}

	allBlocks := s.plane.GetAllBlocks()
	total := len(allBlocks)

	// Calculate start and end for slicing
	start := (page - 1) * limit
	end := start + limit

	// Adjust start and end if they are out of bounds
	if start > total {
		start = total // This will result in an empty slice if start was already beyond total
	}
	if end > total {
		end = total
	}
        
    var paginatedBlocks []*block.Block3D
    if start >= total { // If start index is beyond or at the end of the slice
        paginatedBlocks = []*block.Block3D{} // Return empty slice
    } else {
        paginatedBlocks = allBlocks[start:end]
    }


	totalPages := 0
	if total > 0 && limit > 0 { // Avoid division by zero if no items or limit is zero (though limit is validated >0)
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}
    if total == 0 { // Ensure totalPages is 0 if total is 0, even if page is 1
        totalPages = 0
    }
    if totalPages == 0 && total > 0 { // If total > 0 but totalPages calculated to 0 (e.g. total < limit), it should be 1
        totalPages = 1
    }


	response := gin.H{
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
		"data":       paginatedBlocks,
	}
	c.JSON(http.StatusOK, response)
}

// getBlockByIDHandler handles GET requests to /block/:id.
// It retrieves a specific block by its ID from the plane.
func (s *APIServer) getBlockByIDHandler(c *gin.Context) {
	log.Printf("Incoming %s %s request", c.Request.Method, c.Request.URL.Path)
	id := c.Param("id")
	block, err := s.plane.GetBlockByID(id)
	if err != nil {
		log.Printf("Error getting block by ID %s: %v", id, err)
		// Assuming GetBlockByID returns an error type that indicates "not found"
		// For example, if it returns a specific error like plane.ErrBlockNotFound
		// For now, we check if err is non-nil, which means block not found or other error.
		// The current GetBlockByID returns fmt.Errorf, so direct type assertion isn't straightforward.
		// A more robust solution would be to define custom error types in the plane package.
		c.JSON(http.StatusNotFound, gin.H{"error": "block not found", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, block)
}

// addTransactionHandler handles POST requests to /transaction/add.
// It adds a new transaction to the transaction pool.
func (s *APIServer) addTransactionHandler(c *gin.Context) {
	log.Printf("Incoming %s %s request", c.Request.Method, c.Request.URL.Path)

	var receivedTx block.Transaction
	if err := c.BindJSON(&receivedTx); err != nil {
		log.Printf("Error binding JSON for addTransaction: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction format", "details": err.Error()})
		return
	}

	// Basic Input Validation
	const maxNameLength = 256 // Define a reasonable max length
	if len(receivedTx.Sender) > maxNameLength {
		log.Printf("AddTransaction failed: sender name too long (max %d chars). Length: %d", maxNameLength, len(receivedTx.Sender))
		c.JSON(http.StatusBadRequest, gin.H{"error": "sender name too long", "max_length": maxNameLength})
		return
	}
	if len(receivedTx.Receiver) > maxNameLength {
		log.Printf("AddTransaction failed: receiver name too long (max %d chars). Length: %d", maxNameLength, len(receivedTx.Receiver))
		c.JSON(http.StatusBadRequest, gin.H{"error": "receiver name too long", "max_length": maxNameLength})
		return
	}
	if receivedTx.Amount <= 0 {
		log.Printf("AddTransaction failed: amount must be positive. Amount: %f", receivedTx.Amount)
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive"})
		return
	}
	if receivedTx.Sender == "" || receivedTx.Receiver == "" {
		log.Printf("AddTransaction failed: sender or receiver is empty. Sender: '%s', Receiver: '%s'", receivedTx.Sender, receivedTx.Receiver)
		c.JSON(http.StatusBadRequest, gin.H{"error": "sender and receiver cannot be empty"})
		return
	}

	var validNamePattern = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	if !validNamePattern.MatchString(receivedTx.Sender) {
		log.Printf("AddTransaction failed: sender name contains invalid characters. Sender: '%s'", receivedTx.Sender)
		c.JSON(http.StatusBadRequest, gin.H{"error": "sender name contains invalid characters", "details": "Sender must be alphanumeric and can include _ or -."})
		return
	}
	if !validNamePattern.MatchString(receivedTx.Receiver) {
		log.Printf("AddTransaction failed: receiver name contains invalid characters. Receiver: '%s'", receivedTx.Receiver)
		c.JSON(http.StatusBadRequest, gin.H{"error": "receiver name contains invalid characters", "details": "Receiver must be alphanumeric and can include _ or -."})
		return
	}

	// Perform Signature Verification
	if err := block.VerifySignature(receivedTx); err != nil {
		log.Printf("Transaction signature verification failed for Tx ID %s (or provisional ID): %v", receivedTx.ID, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid transaction signature", "details": err.Error()})
		return
	}
	log.Printf("Transaction signature verified successfully for Tx ID %s", receivedTx.ID)

	// Recalculate and Verify/Set Transaction ID
	expectedIDBytes := block.HashTransactionContent(receivedTx)
	expectedID := hex.EncodeToString(expectedIDBytes[:])

	if receivedTx.ID != expectedID {
		log.Printf("Transaction ID mismatch for received Tx. Client ID: %s, Server-calculated ID: %s. Content was verified against client signature, but ID does not match content hash.", receivedTx.ID, expectedID)
		if receivedTx.ID != "" { // Client provided an ID and it's wrong
			c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID mismatch", "details": fmt.Sprintf("Provided ID '%s' does not match calculated content hash ID '%s'", receivedTx.ID, expectedID)})
			return
		}
		// If client ID was empty, or we choose to always use server-calculated ID (current policy is to set if empty):
		log.Printf("Setting transaction ID to server-calculated hash: %s (Client ID was: '%s')", expectedID, receivedTx.ID)
		receivedTx.ID = expectedID
	} else {
		log.Printf("Transaction ID %s matches calculated content hash.", receivedTx.ID)
	}

	// Add to Transaction Pool
	if err := s.txpool.AddTransaction(receivedTx); err != nil { // Pass receivedTx directly
		log.Printf("Failed to add transaction %s to pool: %v", receivedTx.ID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to add transaction to pool", "details": err.Error()})
		return
	}

	// Success Response
	log.Printf("Transaction %s accepted and added to pool", receivedTx.ID)

	// Broadcast WebSocket message for new transaction
	wsMsgTx := ws.WSMessage{Type: "new_tx", Data: receivedTx}
	jsonMessageTx, err := json.Marshal(wsMsgTx)
	if err != nil {
		log.Printf("Error marshaling new_tx WebSocket message: %v", err)
		// Don't fail the HTTP request, just log the WS broadcast error
	} else {
		s.wsHub.Broadcast(jsonMessageTx)
	}

	c.JSON(http.StatusCreated, gin.H{"status": "Transaction accepted", "id": receivedTx.ID}) // Use StatusCreated
}

// getTxPoolHandler handles GET requests to /txpool.
// It retrieves all pending transactions from the transaction pool as a paginated JSON response.
func (s *APIServer) getTxPoolHandler(c *gin.Context) {
	log.Printf("Incoming %s %s request (with query: %s)", c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery)

	// Parse page query parameter
	pageStr := c.DefaultQuery("page", strconv.Itoa(DefaultPage))
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		log.Printf("Invalid page parameter '%s' for txpool, defaulting to %d. Error: %v", pageStr, DefaultPage, err)
		page = DefaultPage
	}

	// Parse limit query parameter
	limitStr := c.DefaultQuery("limit", strconv.Itoa(DefaultLimit))
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		log.Printf("Invalid limit parameter '%s' for txpool, defaulting to %d. Error: %v", limitStr, DefaultLimit, err)
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		log.Printf("Limit parameter %d for txpool exceeds MaxLimit %d, capping at MaxLimit.", limit, MaxLimit)
		limit = MaxLimit
	}

	allTransactions := s.txpool.GetAllTransactions() // Uses GetAllTransactions
	total := len(allTransactions)

	// Calculate start and end for slicing
	start := (page - 1) * limit
	end := start + limit

	// Adjust start and end if they are out of bounds
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
    
    var paginatedTransactions []block.Transaction
    if start >= total {
        paginatedTransactions = []block.Transaction{} // Return empty slice
    } else {
        paginatedTransactions = allTransactions[start:end]
    }

	totalPages := 0
	if total > 0 && limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}
    if total == 0 {
        totalPages = 0
    }
    if totalPages == 0 && total > 0 {
        totalPages = 1
    }

	response := gin.H{
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
		"data":       paginatedTransactions,
	}
	c.JSON(http.StatusOK, response)
}

// TODO: Implement addBlockHandler if direct block addition via API is desired.
// func (s *APIServer) addBlockHandler(c *gin.Context) {
//    // ... logic to bind block data, validate, and add to plane ...
//    // Note: Adding blocks directly might bypass consensus or other logic.
// }

// addBlockHandler handles POST requests to /block/add.
// It simulates mining a new block by taking coordinates, fetching transactions,
// determining a previous hash, creating a block, adding it to the plane,
// and then removing the transactions from the pool.
func (s *APIServer) addBlockHandler(c *gin.Context) {
	log.Printf("Incoming %s %s request", c.Request.Method, c.Request.URL.Path)
	var req AddBlockRequest
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Error binding JSON for addBlock: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	// Input Validation for addBlockHandler coordinates
	if req.X < 0 || req.Y < 0 || req.Z < 0 {
		log.Printf("addBlockHandler failed: coordinates must be non-negative. X: %d, Y: %d, Z: %d", req.X, req.Y, req.Z)
		c.JSON(http.StatusBadRequest, gin.H{"error": "coordinates must be non-negative"})
		return
	}
	// Optional: Upper bound checks can be added here if needed.

	log.Printf("Block creation initiated for coordinates (%d,%d,%d)", req.X, req.Y, req.Z)

	// Retrieve pending transactions
	// The number of transactions can be configurable. Using 10 as an example.
	transactionsForBlock := s.txpool.GetPendingTransactions(10) // Assuming this function exists and is appropriate
	log.Printf("Retrieved %d transactions for new block at (%d,%d,%d)", len(transactionsForBlock), req.X, req.Y, req.Z)

	// Determine previousHash
	var previousHash string
	// Try to get the block at (req.X, req.Y, req.Z - 1)
	// Note: GetBlock(x,y,z) is not defined in plane.go, but GetBlockByID is.
	// We need to iterate all blocks or have a way to get a block by coordinates.
	// For simplicity, let's assume Plane3D has or will have a GetBlockByCoords(x,y,z) method.
	// If not, this logic needs adjustment.
	// The current Plane3D.Grid is map[string]*Block3D, so direct coordinate lookup isn't trivial.
	// Let's adjust this: iterate blocks to find one at Z-1. This is inefficient.
	// A better way would be for Plane3D to offer a method like GetBlockAt(x,y,z).
	// For now, let's use a simplified previous hash determination:
	// If Z > 0, attempt to find a block at (X, Y, Z-1).
	// This is a placeholder for more complex parent selection logic.
	if req.Z > 0 {
		// This is a simplified search. A real system would need an efficient way to find potential parents.
		// We'll iterate through all blocks to find one at (req.X, req.Y, req.Z-1)
		allBlocks := s.plane.GetAllBlocks() // This could be inefficient for large planes
		foundPrevious := false
		for _, b := range allBlocks {
			if b.X == req.X && b.Y == req.Y && b.Z == (req.Z-1) {
				previousHash = b.Hash
				foundPrevious = true
				log.Printf("Found previous block %s for new block at (%d,%d,%d)", b.Hash, req.X, req.Y, req.Z)
				break
			}
		}
		if !foundPrevious {
			previousHash = "0" // No block found directly below, use default
			log.Printf("No previous block found directly below (%d,%d,%d-1), using '0' as PreviousHash", req.X, req.Y, req.Z)
		}
	} else {
		previousHash = "0" // Block is at Z=0, considered a base block in this context
		log.Printf("Block at Z=0, using '0' as PreviousHash for block at (%d,%d,%d)", req.X, req.Y, req.Z)
	}

	// Create a new block
	newBlock := block.NewBlock(req.X, req.Y, req.Z, previousHash, transactionsForBlock)
	if newBlock == nil { // Should not happen with current NewBlock, but good practice
		log.Printf("Failed to create new block instance for (%d,%d,%d)", req.X, req.Y, req.Z)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create new block"})
		return
	}
	log.Printf("New block %s created (not yet added) for coordinates (%d,%d,%d) with %d transactions.", newBlock.Hash, newBlock.X, newBlock.Y, newBlock.Z, len(newBlock.Transactions))

	// Add the new block to the plane
	if err := s.plane.AddBlock(req.X, req.Y, req.Z, newBlock); err != nil {
		log.Printf("Failed to add block %s to plane at (%d,%d,%d): %v", newBlock.Hash, req.X, req.Y, req.Z, err)
		// Check if the error is due to position already occupied (or other specific errors)
		// For now, using 409 Conflict for any AddBlock error.
		c.JSON(http.StatusConflict, gin.H{"error": "failed to add block to plane", "details": err.Error()})
		return
	}
	// Log for AddBlock in plane.go covers the "added to plane" part.
	log.Printf("Block %s successfully processed by addBlockHandler for plane at (%d,%d,%d)", newBlock.Hash, req.X, req.Y, req.Z)


	// If block added successfully, remove transactions from the pool
	if len(transactionsForBlock) > 0 {
		log.Printf("Removing %d transactions from pool for block %s...", len(transactionsForBlock), newBlock.Hash)
		for _, tx := range transactionsForBlock {
			err := s.txpool.RemoveTransactionByID(tx.ID)
			if err != nil {
				// Log error, but don't fail the whole block addition.
				// This transaction might have been removed by another process, or an error occurred.
				log.Printf("Error removing transaction %s from pool: %v (block %s)", tx.ID, err, newBlock.Hash)
			}
			// Log for RemoveTransactionByID in txpool.go covers the "removed from pool" part.
		}
		log.Printf("Finished removing transactions from pool for block %s.", newBlock.Hash)
	}

	// Broadcast WebSocket message for new block
	wsMsgBlock := ws.WSMessage{Type: "new_block", Data: newBlock}
	jsonMessageBlock, err := json.Marshal(wsMsgBlock)
	if err != nil {
		log.Printf("Error marshaling new_block WebSocket message: %v", err)
		// Don't fail the HTTP request, just log the WS broadcast error
	} else {
		s.wsHub.Broadcast(jsonMessageBlock)
	}

	c.JSON(http.StatusCreated, newBlock)
}

// getPlaneGridHandler handles GET requests to /plane/grid.
// It returns all blocks in the plane, effectively giving the "layout".
func (s *APIServer) getPlaneGridHandler(c *gin.Context) {
	log.Printf("Incoming %s %s request", c.Request.Method, c.Request.URL.Path)
	blocks := s.plane.GetAllBlocks()
	// Currently, GetAllBlocks doesn't return an error.
	// If it did, error handling would be needed here.
	c.JSON(http.StatusOK, blocks)
}

// handleWebSocketConnections is a method on APIServer that wraps ws.HandleConnections.
func (s *APIServer) handleWebSocketConnections(c *gin.Context) {
	ws.HandleConnections(s.wsHub)(c) // Calls the handler from the ws package
}

// CORSMiddleware sets up permissive CORS headers for development.
// For production, origins should be restricted.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // Use 204 No Content for OPTIONS preflight
			return
		}
		c.Next()
	}
}

// RequireJSONContentType is a middleware that ensures requests that typically have a body
// (POST, PUT, PATCH) have the Content-Type header set to "application/json".
func RequireJSONContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == "POST" || method == "PUT" || method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			// It's common for Content-Type to include charset, e.g., "application/json; charset=utf-8"
			// So, we check if it starts with "application/json".
			if !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
				log.Printf("Request with method %s rejected due to invalid Content-Type: '%s'", method, contentType)
				c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
					"error":   "Unsupported Media Type",
					"details": "Requests with a body (POST, PUT, PATCH) must have Content-Type set to 'application/json'.",
				})
				return
			}
		}
		c.Next()
	}
}

// globalErrorHandler is a middleware to catch panics and respond with a JSON error.
// It also logs any errors that were added to c.Errors but not handled by other middleware/handlers.
func globalErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Panic recovered in globalErrorHandler: %v", r)
                // Ensure a response hasn't already been sent
                if !c.Writer.Written() {
                    c.JSON(http.StatusInternalServerError, gin.H{
                        "error":   "Internal Server Error",
                        "details": "The server encountered an unrecoverable situation.",
                    })
                }
                c.Abort() // Abort further processing
            }
        }()

        c.Next() // Process request

        // After request, log any errors that were attached to the context but not handled.
        // This is useful for debugging if a handler sets c.Error() but doesn't abort.
        // Note: Most of our handlers already use c.JSON for errors, which aborts.
        if len(c.Errors) > 0 {
            for _, e := range c.Errors {
                log.Printf("Unhandled error in context: %v", e.Err)
            }
            // Optionally, if no response was written, send a generic error.
            // However, this might interfere if a later middleware is supposed to handle it.
            // if !c.Writer.Written() && !c.IsAborted() {
            //    c.JSON(http.StatusInternalServerError, gin.H{"error": "Unhandled error occurred"})
            // }
        }
    }
}
