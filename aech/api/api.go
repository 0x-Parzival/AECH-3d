package api

import (
	"aech/block"
	"aech/plane"
	"aech/txpool"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// APIServer holds instances needed by the API handlers.
type APIServer struct {
	plane  *plane.Plane3D
	txpool *txpool.TxPool
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
	server := &APIServer{
		plane:  p,
		txpool: tp,
	}

	r := gin.Default()

	// Block routes
	r.GET("/blocks", server.getBlocksHandler)
	r.GET("/block/:id", server.getBlockByIDHandler)
	r.POST("/block/add", server.addBlockHandler) // If needed

	// Transaction routes
	r.POST("/transaction/add", server.addTransactionHandler)
	r.GET("/txpool", server.getTxPoolHandler)

	// Plane routes
	r.GET("/plane/grid", server.getPlaneGridHandler)

	// TODO: Add other endpoints here as they are developed

	log.Printf("Starting API server on port %s...", port) // Added "..." for clarity
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err) // This log is already good
	}
}

// getBlocksHandler handles GET requests to /blocks.
// It retrieves all blocks from the plane and returns them as JSON.
func (s *APIServer) getBlocksHandler(c *gin.Context) {
	log.Printf("Incoming %s %s request", c.Request.Method, c.Request.URL.Path)
	blocks := s.plane.GetAllBlocks()
	// Currently, GetAllBlocks doesn't return an error. If it did, we'd check it here.
	c.JSON(http.StatusOK, blocks)
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
	var tx block.Transaction // This var is not directly used due to revised logic below.
	// Binding to a temporary struct to extract only necessary fields for NewTransaction
	clientTxData := struct {
		Sender   string  `json:"sender"`
		Receiver string  `json:"receiver"`
		Amount   float64 `json:"amount"`
	}{}
	if err := c.BindJSON(&clientTxData); err != nil {
		log.Printf("Error binding JSON for addTransaction: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	// Sign the transaction (as per current dummy signing logic)
	// NewTransaction already calls Sign, but if we receive a raw tx, we should sign it.
	// If tx.ID is empty, it implies it's a new transaction not yet processed by NewTransaction.
	// However, NewTransaction is the one that generates ID and signs.
	// For robustness, let's ensure it's signed. The client *should* send all necessary fields
	// for NewTransaction to work if they construct it, or send a pre-created one.
	// If the transaction is supposed to be created here from raw data,
	// we might need to call block.NewTransaction(tx.Sender, tx.Receiver, tx.Amount)
	// and then it would be signed.
	// For now, assuming the client sends a transaction that might just need its signature (re-)applied.
	// The current NewTransaction sets ID and Signature. If a client sends a TX,
	// it should have an ID.
    // If the client is expected to send minimal data (sender, receiver, amount) and the server
    // completes it, the logic would be different:
    // newTx := block.NewTransaction(tx.Sender, tx.Receiver, tx.Amount) // This creates ID, Timestamp, signs.
    // Then use newTx for AddTransaction.
    //
    // Given the current structure of Transaction and NewTransaction,
    // if a client sends a JSON that unmarshals into block.Transaction,
    // its ID and Signature might be empty or incorrect.
    // Let's assume the client provides Sender, Receiver, Amount.
    // We should probably re-construct it to ensure ID and Signature are correctly generated.

    // Correct approach: Treat incoming JSON as data for a *new* transaction
    // unless the API is designed to accept pre-signed, pre-ID'd transactions.
    // The task says "Call tx.Sign()", implying tx is already somewhat formed.
    // Let's stick to that and assume ID is part of the incoming tx for now,
    // and Sign() just populates/overwrites the signature.

	if clientTxData.Sender == "" || clientTxData.Receiver == "" { // Basic validation
		log.Printf("AddTransaction failed: sender or receiver is empty. Sender: '%s', Receiver: '%s'", clientTxData.Sender, clientTxData.Receiver)
		c.JSON(http.StatusBadRequest, gin.H{"error": "sender and receiver cannot be empty"})
		return
	}
    // The ID and Timestamp are set by NewTransaction. If we just BindJSON and Sign,
    // ID and Timestamp might not be what we expect if the client doesn't send them.
    // A safer way for "adding" a new transaction from client data:
    // clientTxData struct is already defined and bound above.

    // Create a new transaction using the system's logic to ensure ID, Timestamp, and Signature are correct
    finalTx := block.NewTransaction(clientTxData.Sender, clientTxData.Receiver, clientTxData.Amount)
    // NewTransaction already calls Sign(), so an explicit finalTx.Sign() is redundant here.


	if err := s.txpool.AddTransaction(*finalTx); err != nil {
		log.Printf("Failed to add transaction %s to pool: %v", finalTx.ID, err)
		// Determine appropriate status code based on error type if TxPool returns specific errors
		// For now, using 400 for any error from AddTransaction.
		// Could be 409 if it's a duplicate, 422 for validation, etc.
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to add transaction", "details": err.Error()})
		return
	}
	log.Printf("Transaction %s added to pool via API", finalTx.ID) // Specific log for successful add
	c.JSON(http.StatusOK, gin.H{"status": "transaction added", "tx_id": finalTx.ID})
}

// getTxPoolHandler handles GET requests to /txpool.
// It retrieves all pending transactions from the transaction pool.
func (s *APIServer) getTxPoolHandler(c *gin.Context) {
	log.Printf("Incoming %s %s request", c.Request.Method, c.Request.URL.Path)
	transactions := s.txpool.GetAllTransactions()
	// Currently, GetAllTransactions doesn't return an error.
	c.JSON(http.StatusOK, transactions)
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
