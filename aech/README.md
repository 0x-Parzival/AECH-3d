# AECH: 3D Blockchain Financial Operating System (Core Implementation)

This repository contains the initial core implementation of AECH, a conceptual 3D blockchain-based financial operating system. This early version focuses on the fundamental structures for creating blocks, storing them in a 3D spatial grid, and initializing a Genesis block.

## Current Features

*   **Block3D Structure**: Defines blocks with `X,Y,Z` coordinates, `Timestamp`, `Hash`, `PreviousHash`, and a list of `Transactions`.
*   **Plane System**: A 3D map (`Cube`) for storing and retrieving `Block3D` objects based on their spatial coordinates.
*   **Block Operations**:
    *   `block.NewBlock()`: Creates new blocks, automatically generating a timestamp and hash.
    *   `block.GenerateHash()`: Computes a SHA256 hash for a block.
    *   `block.ValidateBlock()`: A basic stub for block validation.
*   **Plane Operations**:
    *   `plane.InsertBlock()`: Inserts a block into the 3D `Cube`.
    *   `plane.GetBlock()`: Retrieves a block from specified coordinates.
    *   `plane.GetNeighbors()`: Finds adjacent blocks in the 6 cardinal directions.
*   **Genesis Block**: `main.go` initializes the blockchain with a Genesis block at coordinates (0,0,0).
*   **Unit Tests**: Basic tests for core `block` and `plane` functionalities are available in the `/test` directory.

## Directory Structure

- `/block/`: Contains the `Block3D` structure definition and related functions.
- `/plane/`: Contains the `Cube` data structure for 3D spatial storage and related functions.
- `/test/`: Contains unit tests for the project (e.g., `block_plane_test.go`).
- `main.go`: Entry point of the application; currently initializes the Genesis block.
- `go.mod`, `go.sum`: Go module files.

## Block Structure Specification (`block.Block3D`)

- `X, Y, Z (int)`: Coordinates of the block in the 3D space.
- `Timestamp (time.Time)`: Time of block creation.
- `Hash (string)`: SHA256 hash derived from X,Y,Z, Timestamp, PreviousHash, and concatenated transaction hashes.
- `PreviousHash (string)`: Hash of a conceptual preceding block. For the Genesis block, this is "0".
- `Transactions ([]Transaction)`: A list of transactions included in the block.
    - `Transaction struct { Data []byte }`: Placeholder for actual transaction data.

## How to Run

1.  **Prerequisites**:
    *   Go (version 1.21 or later recommended).
2.  **Clone the repository**:
    ```bash
    # git clone <repository-url>
    # cd aech
    ```
3.  **Initialize Go modules** (if you haven't run any `go` commands yet):
    ```bash
    go mod tidy
    ```
4.  **Run the main application** (initializes Genesis block):
    ```bash
    go run main.go
    ```
5.  **Run tests**:
    ```bash
    go test ./...
    ```

## Sample Output (`go run main.go`)

The output will be similar to (timestamps will vary):

```log
2025/05/26 18:47:37 Genesis block created and inserted at coordinates (0, 0, 0)
2025/05/26 18:47:37 Genesis block Hash: 05676806c90e6b384b9bcaf6aaddbe494a959fdf87390c9e72734965be91ef84
2025/05/26 18:47:37 Genesis block Timestamp: 2025-05-26 18:47:37.855062506 +0000 UTC m=+0.000084618
```

## Architecture Diagram

*(Placeholder: A visual diagram, possibly using ASCII art or a linked image (e.g., from draw.io), will be added here in future iterations to better illustrate the 3D cube structure and block relationships.)*
