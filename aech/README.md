# AECH Backend

AECH is a simplified blockchain and 3D ledger system created for educational purposes.
It aims to demonstrate basic concepts of blockchain technology, including:

- Blocks and Transactions
- Cryptographic Hashing
- Transaction Pools
- A 3D "Plane" for spatial block organization (conceptual)
- Basic Consensus (stubbed)
- API for interaction

## Run the Server

To start the AECH backend server:

**1. Configure Port (Optional):**
   Create a `.env` file in the `aech` project root (you can copy `.env.example`):
   ```env
   PORT=8081 
   ```
   If `.env` is not present or `PORT` is not set, the server defaults to port `8080`.

**2. Run Directly:**
   ```bash
   go run main.go
   ```

**3. Build and Run Executable (Recommended for stable execution):**
   ```bash
   # From the aech directory
   go build -o aech_server main.go
   ./aech_server
   ```
The server will log its startup and the port it's using. You'll see further log messages in your console for incoming requests and other operations.

## Running Tests

To run the unit tests for the various packages:

```bash
# From the aech directory
go test ./...
```
This command will execute all test files (`*_test.go`) in the current directory and its subdirectories. Add the `-v` flag for verbose output: `go test -v ./...`.

## API Reference

The API allows you to interact with the AECH system.

### Add a Transaction

**Endpoint:** `POST /transaction/add`

**Purpose:** Submits a new transaction to the transaction pool.

**Request Body (JSON):**
```json
{
  "sender": "Alice",
  "receiver": "Bob",
  "amount": 10.5
}
```
*   `sender` (string): The identifier of the transaction sender.
*   `receiver` (string): The identifier of the transaction receiver.
*   `amount` (float64): The amount to be transferred.

**Response (Success - 200 OK):**
```json
{
  "status": "transaction added",
  "tx_id": "generated_transaction_id_hash"
}
```

### Add a Block (Simulate Mining)

**Endpoint:** `POST /block/add`

**Purpose:** Simulates the mining of a new block by taking pending transactions from the pool and adding them to a new block on the plane.

**Request Body (JSON):**
```json
{
  "x": 1,
  "y": 2,
  "z": 0
}
```
*   `x` (int): The X-coordinate for the new block.
*   `y` (int): The Y-coordinate for the new block.
*   `z` (int): The Z-coordinate for the new block.

**Response (Success - 201 Created):**
A `block.Block3D` object. Example:
```json
{
    "X": 1,
    "Y": 2,
    "Z": 0,
    "Timestamp": "2023-10-27T10:05:00.123456789Z",
    "Hash": "generated_block_hash",
    "PreviousHash": "hash_of_block_at_x_y_z-1_or_0",
    "Transactions": [
        {
            "ID": "tx_id_1",
            "Sender": "SenderName",
            "Receiver": "ReceiverName",
            "Amount": 10.0,
            "Timestamp": 1678886400000000000,
            "Signature": "7369676e65642d74785f69645f31"
        }
    ]
}
```
*Note: `Block3D.Timestamp` is marshaled to RFC3339Nano format. `Transaction.Timestamp` is UnixNano (int64). `Transaction.Signature` is hex-encoded.*

### Get All Blocks

**Endpoint:** `GET /blocks`

**Purpose:** Retrieves a list of all blocks currently in the 3D plane.

**Response (Success - 200 OK):**
An array of `block.Block3D` objects. Example:
```json
[
  {
    "X": 0,
    "Y": 0,
    "Z": 0,
    "Timestamp": "2023-10-27T10:00:00Z",
    "Hash": "abc123genesis",
    "PreviousHash": "0",
    "Transactions": []
  },
  {
    "X": 1,
    "Y": 2,
    "Z": 0,
    "Timestamp": "2023-10-27T10:05:00.123456789Z",
    "Hash": "def456block2",
    "PreviousHash": "abc123genesis",
    "Transactions": [
        {
            "ID": "tx_id_1",
            "Sender": "SenderName",
            "Receiver": "ReceiverName",
            "Amount": 10.0,
            "Timestamp": 1678886400000000000,
            "Signature": "7369676e65642d74785f69645f31"
        }
    ]
  }
]
```

### Get Block by ID

**Endpoint:** `GET /block/:id`

**Purpose:** Retrieves details for a specific block by its hash ID.

**URL Parameter:**
*   `id` (string): The hash of the block to retrieve.

**Response (Success - 200 OK):**
A single `block.Block3D` object, similar to the one shown in the `POST /block/add` success response.

**Response (Failure - 404 Not Found):**
Example:
```json
{
    "details": "block with ID your_block_id not found",
    "error": "block not found"
}
```

### Get Transaction Pool

**Endpoint:** `GET /txpool`

**Purpose:** Retrieves all transactions currently pending in the transaction pool.

**Response (Success - 200 OK):**
An array of `block.Transaction` objects. Example:
```json
[
    {
        "ID": "tx_id_pending_1",
        "Sender": "SenderPending",
        "Receiver": "ReceiverPending",
        "Amount": 5.0,
        "Timestamp": 1678886500000000000,
        "Signature": "7369676e65642d74785f69645f70656e64696e675f31"
    }
]
```
*Note: `Transaction.Timestamp` is UnixNano (int64). `Transaction.Signature` is hex-encoded.*

### Get Plane Grid Layout

**Endpoint:** `GET /plane/grid`

**Purpose:** Returns the current layout of the 3D plane by listing all blocks. This is effectively the same as `GET /blocks`.

**Response (Success - 200 OK):**
An array of `block.Block3D` objects, same as `GET /blocks`.

## Modules

The AECH system is organized into the following packages:

- `block`: Defines `Block3D` and `Transaction` structures, and core operations like hashing and signing.
- `txpool`: Manages the pool of pending transactions before they are included in blocks.
- `plane`: Manages the 3D grid representation where blocks are spatially organized.
- `api`: Provides the HTTP API interface using the Gin framework.
- `main.go`: The entry point for the application. It initializes all components (transaction pool, plane) and starts the API server.

Other conceptual modules (currently stubbed or not fully implemented):
- `consensus`: Intended for consensus algorithm integration.
- `wallet`: Intended for user wallet functionalities.
- `ai`: Conceptual module for AI-driven interactions or analysis.
- `bridge`: Conceptual module for cross-chain communication.
- `crosschain`: Support for cross-chain interactions.
- `explorer`: For a potential block explorer interface.
- `frontend`: For a potential user interface.

## Disclaimer

This project is for learning and demonstration. It is not intended for production use and lacks many security features and optimizations found in real-world blockchain systems.
```
