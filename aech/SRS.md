# Software Requirements Specification (SRS) for AECH

**Version 1.1**

**Prepared by:** AECH Development Team

**Date:** [Current Date]

## Table of Contents

1.  [Introduction](#1-introduction)
    1.  [Purpose](#11-purpose)
    2.  [Scope](#12-scope)
    3.  [Definitions, Acronyms, and Abbreviations](#13-definitions-acronyms-and-abbreviations)
    4.  [References](#14-references)
    5.  [Overview](#15-overview)
2.  [Overall Description](#2-overall-description)
    1.  [Product Perspective](#21-product-perspective)
    2.  [Product Functions](#22-product-functions)
    3.  [User Characteristics](#23-user-characteristics)
    4.  [Constraints](#24-constraints)
    5.  [Assumptions and Dependencies](#25-assumptions-and-dependencies)
3.  [Specific Requirements](#3-specific-requirements)
    1.  [Functional Requirements](#31-functional-requirements)
        1.  [FR1: Transaction Creation](#311-fr1-transaction-creation)
        2.  [FR2: Transaction Validation](#312-fr2-transaction-validation)
        3.  [FR3: Transaction Pool Management](#313-fr3-transaction-pool-management)
        4.  [FR4: Block Creation (Mining Simulation)](#314-fr4-block-creation-mining-simulation)
        5.  [FR5: Block Validation](#315-fr5-block-validation)
        6.  [FR6: Plane (3D Ledger) Management](#316-fr6-plane-3d-ledger-management)
        7.  [FR7: API Endpoints](#317-fr7-api-endpoints)
        8.  [FR8: WebSocket Real-time Updates](#318-fr8-websocket-real-time-updates)
        9.  [FR9: Cross-Chain Communication (Mock)](#319-fr9-cross-chain-communication-mock)
    2.  [Non-Functional Requirements](#32-non-functional-requirements)
        1.  [NFR1: Performance](#321-nfr1-performance)
        2.  [NFR2: Scalability](#322-nfr2-scalability)
        3.  [NFR3: Security](#323-nfr3-security)
        4.  [NFR4: Usability (API)](#324-nfr4-usability-api)
        5.  [NFR5: Maintainability](#325-nfr5-maintainability)
        6.  [NFR6: Configurability](#326-nfr6-configurability)
    3.  [Interface Requirements](#33-interface-requirements)
        1.  [API Interface](#331-api-interface)
        2.  [WebSocket Interface](#332-websocket-interface)
4.  [Other Requirements](#4-other-requirements)
    1.  [Logging](#41-logging)
    2.  [Error Handling](#42-error-handling)
5.  [IEEE 830 Compliance Verification](#5-ieee-830-compliance-verification)

---

## 1. Introduction

### 1.1 Purpose

This Software Requirements Specification (SRS) document describes the functional and non-functional requirements for the AECH (Another Educational CHain) system. AECH is a backend system designed to simulate basic blockchain and 3D ledger concepts for educational and demonstrative purposes. This document is intended for developers, testers, and project stakeholders.

### 1.2 Scope

The scope of AECH includes:
-   Core data structures: Transactions and Blocks (with 3D spatial coordinates).
-   Core operations: Transaction creation (client-side signed), signature verification, block creation (simulated mining), block validation.
-   Management of pending transactions in a transaction pool.
-   Management of blocks in a conceptual 3D spatial grid ("Plane").
-   A RESTful HTTP API for interacting with the system.
-   Real-time updates via WebSockets for new transactions and blocks.
-   Mocked cross-chain communication capabilities.
-   Basic security considerations for the API.
-   Configurable parameters such as server port and CORS origins via environment variables.

The system is **not** intended for production use and will explicitly exclude:
-   Complex consensus algorithms (e.g., Proof-of-Work, Proof-of-Stake). Consensus is stubbed.
-   Advanced cryptographic techniques beyond basic ECDSA for signatures.
-   Peer-to-peer networking for distributed ledger synchronization.
-   User account management and authentication beyond public key identifiers.
-   A persistent storage layer (data is in-memory).
-   A graphical user interface (GUI).

### 1.3 Definitions, Acronyms, and Abbreviations

-   **AECH**: Another Educational CHain, the name of the system.
-   **API**: Application Programming Interface.
-   **Block**: A data structure containing a set of transactions, a hash of the previous block, a timestamp, and spatial coordinates.
-   **Transaction**: A record of transfer or interaction, digitally signed by the sender.
-   **TxPool**: Transaction Pool, a collection of pending transactions.
-   **Plane**: The conceptual 3D grid where blocks are spatially organized.
-   **SRS**: Software Requirements Specification.
-   **JSON**: JavaScript Object Notation.
-   **HTTP**: Hypertext Transfer Protocol.
-   **WebSocket**: A communication protocol providing full-duplex communication channels over a single TCP connection.
-   **CORS**: Cross-Origin Resource Sharing.
-   **ECDSA**: Elliptic Curve Digital Signature Algorithm.
-   **P-256**: A specific elliptic curve (also known as secp256r1 or prime256v1).

### 1.4 References

-   Go Programming Language Specification ([https://golang.org/ref/spec](https://golang.org/ref/spec))
-   Gin Web Framework ([https://gin-gonic.com/docs/](https://gin-gonic.com/docs/))
-   Gorilla WebSocket ([https://github.com/gorilla/websocket](https://github.com/gorilla/websocket))
-   IEEE Std 830-1998: Recommended Practice for Software Requirements Specifications (conceptual model for structure).

### 1.5 Overview

This document is organized into several sections:
-   **Section 1 (Introduction)**: Provides context for the SRS.
-   **Section 2 (Overall Description)**: Describes the general factors affecting the product and its requirements.
-   **Section 3 (Specific Requirements)**: Details the functional, non-functional, and interface requirements.
-   **Section 4 (Other Requirements)**: Covers aspects like logging and error handling.
-   **Section 5 (IEEE 830 Compliance Verification)**: Maps sections of this document to the IEEE 830 standard.

---

## 2. Overall Description

### 2.1 Product Perspective

AECH is a self-contained backend application. It provides an API for external clients (e.g., command-line tools, web frontends, testing scripts) to interact with its simulated blockchain environment. It is designed as a modular Go application, with distinct packages for core logic (blocks, transactions), spatial organization (plane), transaction pooling, and API handling.

### 2.2 Product Functions

The primary functions of AECH are:
1.  **Transaction Management**: Accepting client-signed transactions, validating their signatures and basic structure, and managing them in a transaction pool.
2.  **Block Management**: Creating new blocks from pending transactions (simulating mining), validating blocks, and organizing them in a conceptual 3D spatial grid.
3.  **API Provision**: Offering HTTP API endpoints for creating transactions, adding blocks, and querying system state (blocks, transaction pool).
4.  **Real-time Notifications**: Broadcasting information about new transactions and blocks to connected WebSocket clients.
5.  **Cross-Chain Simulation**: Mocking the sending and receiving of messages to/from other (simulated) blockchain planes.

### 2.3 User Characteristics

The anticipated users of the AECH system are:
-   **Developers**: Using the API to build applications or test integrations with blockchain-like systems.
-   **Testers**: Verifying the functionality and behavior of the system.
-   **Students/Educators**: Using AECH as a tool to learn or teach basic blockchain concepts.

Users are expected to have technical knowledge regarding APIs, JSON, and fundamental blockchain principles.

### 2.4 Constraints

-   The system must be implemented in Go.
-   The API must use the Gin Web Framework.
-   WebSocket communication must use the Gorilla WebSocket library.
-   The system will be an in-memory simulation; no persistent database is required.
-   Cryptographic operations will be limited to SHA-256 for hashing and ECDSA (P-256 curve) for transaction signatures.
-   Development will follow standard Go project structure and practices.
-   The system must be configurable via environment variables for port and CORS settings.

### 2.5 Assumptions and Dependencies

-   Users (clients of the API) are responsible for creating and signing transactions before submission. The server will verify signatures.
-   The "3D Plane" is a conceptual organization and does not imply complex 3D graphical rendering or physics.
-   Cross-chain communication is simulated and does not involve actual external network interactions with other live blockchains.
-   The system will run as a single server instance. Distributed operation is outside the scope.
-   The Go toolchain (compatible with the version specified in `go.mod`, currently Go 1.23.x) is available in the development and execution environment.

---

## 3. Specific Requirements

### 3.1 Functional Requirements

#### 3.1.1 FR1: Transaction Creation
-   **FR1.1**: The system shall define a `Transaction` structure including: `ID` (string, server-set hash of content), `Timestamp` (int64, client-set UnixNano), `Sender` (string, client-set), `Receiver` (string, client-set), `Amount` (float64, client-set), `Signature` (string, client-set hex-encoded ECDSA signature), `PubKey` (string, client-set hex-encoded compressed ECDSA public key).
-   **FR1.2**: The system shall provide a helper function `block.NewTransaction(senderName, receiverName, amount)` for server-side generation of new, signed transactions (e.g., for testing or coinbase simulation). This function will generate a new ECDSA key pair, set the `PubKey`, hash the core content (Sender, Receiver, Amount, Timestamp) to create the `ID`, and sign this `ID` hash to create the `Signature`.
-   **FR1.3**: The system shall provide a helper function `block.HashTransactionContent(tx Transaction)` that computes a SHA-256 hash of the transaction's core content (Sender, Receiver, Amount, Timestamp) used for signing and ID generation.

#### 3.1.2 FR2: Transaction Validation
-   **FR2.1**: The system shall provide a function `block.VerifySignature(tx Transaction)` to verify the `tx.Signature` against the `tx.PubKey` and the hash of the transaction's core content (derived using `HashTransactionContent`). This function must support hex-encoded compressed public keys (P-256) and hex-encoded ECDSA signatures (R and S components concatenated).
-   **FR2.2**: The API endpoint for adding transactions (`POST /transaction/add`) shall validate incoming transactions:
    -   **FR2.2.1**: Bind the JSON payload to the `block.Transaction` struct.
    -   **FR2.2.2**: Perform basic field validations: `Sender` and `Receiver` names (length, allowed characters: alphanumeric, underscore, hyphen), `Amount` (must be positive).
    -   **FR2.2.3**: Verify the transaction's signature using `block.VerifySignature`.
    -   **FR2.2.4**: Calculate the expected transaction `ID` using `block.HashTransactionContent`. If the client provided an `ID`, it must match this calculated ID. If the client-provided `ID` is empty, the server shall set it to the calculated ID.

#### 3.1.3 FR3: Transaction Pool Management
-   **FR3.1**: The system shall maintain a transaction pool (`TxPool`) for pending transactions.
-   **FR3.2**: `TxPool` shall provide a thread-safe `AddTransaction(tx Transaction)` method. It shall return an error if the transaction already exists in the pool (checked by ID) or if `tx.Validate()` (currently a stub) fails.
-   **FR3.3**: `TxPool` shall provide a `GetPendingTransactions(count int)` method to retrieve a specified number of transactions.
-   **FR3.4**: `TxPool` shall provide a `RemoveTransactionByID(txID string)` method.
-   **FR3.5**: `TxPool` shall provide a `GetAllTransactions()` method to retrieve all transactions currently in the pool (used for API and pagination).

#### 3.1.4 FR4: Block Creation (Mining Simulation)
-   **FR4.1**: The system shall define a `Block3D` structure including: `X, Y, Z` (int spatial coordinates), `Timestamp` (time.Time), `Hash` (string, hash of block content), `PreviousHash` (string), and `Transactions` ([]Transaction).
-   **FR4.2**: The system shall provide a `block.NewBlock(x, y, z, previousHash, transactions)` function that creates a new `Block3D` instance, sets its timestamp, and generates its hash.
-   **FR4.3**: The `Block3D.GenerateHash()` method shall compute a SHA-256 hash based on `X, Y, Z, Timestamp, PreviousHash`, and a string representation of all included transactions (including their IDs, Sender, Receiver, Amount, Timestamps, and Signatures).
-   **FR4.4**: The API endpoint for adding blocks (`POST /block/add`) shall simulate block creation:
    -   **FR4.4.1**: Accept `X, Y, Z` coordinates.
    -   **FR4.4.2**: Retrieve pending transactions from `TxPool` (e.g., up to a configurable limit).
    -   **FR4.4.3**: Determine `PreviousHash` (e.g., by looking for a block at `X, Y, Z-1`, or "0" if none).
    -   **FR4.4.4**: Create a new block using `block.NewBlock`.
    -   **FR4.4.5**: Add the block to the `Plane` using `plane.AddBlock`.
    -   **FR4.4.6**: If block addition to the plane is successful, remove the included transactions from `TxPool`.

#### 3.1.5 FR5: Block Validation
-   **FR5.1**: The system shall provide a `block.ValidateBlock(b Block3D)` function. Initially, this checks if `Hash` is non-empty and coordinates are non-negative. (Future extension: re-calculate and compare hash).

#### 3.1.6 FR6: Plane (3D Ledger) Management
-   **FR6.1**: The system shall maintain a `Plane3D` structure to store blocks in a 3D grid using coordinate-based keys.
-   **FR6.2**: `Plane3D` shall provide `AddBlock(x, y, z, block)`: Adds a block if the position is valid (non-negative coordinates) and not occupied.
-   **FR6.3**: `Plane3D` shall provide `GetBlock(x, y, z)`: Retrieves a block by coordinates.
-   **FR6.4**: `Plane3D` shall provide `GetBlockByID(id string)`: Retrieves a block by its hash ID.
-   **FR6.5**: `Plane3D` shall provide `GetAllBlocks()`: Retrieves all blocks.
-   **FR6.6**: `Plane3D` shall provide `ListNeighbors(x, y, z)`: Lists adjacent blocks.
-   **FR6.7**: `Plane3D` shall provide `IsValidPosition(x, y, z)`: Checks if coordinates are non-negative (initial implementation).

#### 3.1.7 FR7: API Endpoints
The system shall provide the following HTTP API endpoints:
-   **FR7.1**: `POST /transaction/add`: Accepts a new transaction (see FR2.2). Responds with status and transaction ID.
-   **FR7.2**: `POST /block/add`: Simulates mining a new block (see FR4.4). Responds with the created block.
-   **FR7.3**: `GET /blocks`: Retrieves all blocks, paginated (see FR6.5, FR7.7).
-   **FR7.4**: `GET /block/:id`: Retrieves a specific block by its hash ID (see FR6.4).
-   **FR7.5**: `GET /txpool`: Retrieves all pending transactions, paginated (see FR3.5, FR7.7).
-   **FR7.6**: `GET /plane/grid`: Retrieves all blocks (same as `GET /blocks`), paginated (see FR6.5, FR7.7).
-   **FR7.7**: Pagination: `GET /blocks`, `GET /txpool`, and `GET /plane/grid` endpoints shall support `page` and `limit` query parameters. Defaults: page=1, limit=10. Max limit: 100. Response includes `page`, `limit`, `total`, `totalPages`, and `data`.
-   **FR7.8**: `GET /health`: A health check endpoint returning `{"status": "ok", "time": "timestamp"}`.

#### 3.1.8 FR8: WebSocket Real-time Updates
-   **FR8.1**: The system shall provide a WebSocket endpoint `GET /ws`.
-   **FR8.2**: A `Hub` shall manage active WebSocket clients (`Client` channels).
    -   `Hub.Register(client)`: Adds a client.
    -   `Hub.Unregister(client)`: Removes a client and closes its channel.
    -   `Hub.Broadcast(message)`: Sends a message to all registered clients. Handles blocked/slow clients by unregistering them.
-   **FR8.3**: The WebSocket handler (`ws.HandleConnections`) shall upgrade HTTP connections, register clients with the Hub, and manage read/write pumps for each connection.
    -   `writePump`: Sends messages from the client's channel to the WebSocket connection.
    -   `readPump`: Detects client disconnections and unregisters the client. Client-sent messages are logged but ignored.
-   **FR8.4**: The system shall define a `WSMessage` struct (`{Type string, Data interface{}}`) for WebSocket messages.
-   **FR8.5**: After a new transaction is successfully added via `POST /transaction/add`, a `WSMessage` of type `"new_tx"` with the transaction object as `Data` shall be broadcast to all WebSocket clients.
-   **FR8.6**: After a new block is successfully added via `POST /block/add`, a `WSMessage` of type `"new_block"` with the block object as `Data` shall be broadcast to all WebSocket clients.

#### 3.1.9 FR9: Cross-Chain Communication (Mock)
-   **FR9.1**: The system shall define a `CrossChainMessage` struct (`{FromPlaneID, ToPlaneID, BlockHash, Payload, Timestamp}`).
-   **FR9.2**: `crosschain.SendMessage(msg)`: Simulates sending a message by logging it and calling `ReceiveMessage` (local loopback).
-   **FR9.3**: `crosschain.ReceiveMessage(msg)`: Simulates receiving by logging the message.
-   **FR9.4**: When a block is successfully added to the plane in `plane.AddBlock`, a `CrossChainMessage` shall be constructed and "sent" using `crosschain.SendMessage`.

### 3.2 Non-Functional Requirements

#### 3.2.1 NFR1: Performance
-   **NFR1.1**: API responses for read operations (e.g., `GET /blocks`, `GET /block/:id`) should generally complete within 500ms under normal load (simulated environment with up to 1,000 blocks and 10,000 transactions in pool).
-   **NFR1.2**: Transaction submission (`POST /transaction/add`) should complete within 200ms (excluding network latency).
-   **NFR1.3**: WebSocket broadcasts should be initiated within 100ms of the triggering event (new transaction/block).

#### 3.2.2 NFR2: Scalability
-   **NFR2.1**: The system should handle up to 100 concurrent API requests.
-   **NFR2.2**: The system should manage up to 50 concurrent WebSocket clients.
-   **NFR2.3**: The in-memory data structures (plane, txpool) should efficiently manage up to 10,000 blocks and 100,000 transactions in the pool without significant degradation of core functions (add, retrieve).

#### 3.2.3 NFR3: Security
-   **NFR3.1**: API endpoints accepting JSON payloads (e.g., `POST /transaction/add`, `POST /block/add`) must validate the `Content-Type` header to be `application/json`.
-   **NFR3.2**: API endpoints must implement a request body size limit (e.g., 1MB) to prevent excessively large payloads.
-   **NFR3.3**: Input validation must be performed on all client-provided data for API endpoints (e.g., transaction fields, block coordinates) to prevent common injection or processing errors.
-   **NFR3.4**: CORS must be configurable via environment variables (`CORS_ALLOWED_ORIGINS`). Default to `"*"` if not set.
-   **NFR3.5**: Sensitive information (if any were present beyond this simulation) should not be logged excessively. (Current logging is for operational/debug purposes).
-   **NFR3.6**: The system must implement global error handling to catch unhandled panics and return a generic 500 error, logging the panic details.

#### 3.2.4 NFR4: Usability (API)
-   **NFR4.1**: The API should be well-documented in `README.md`, including endpoint descriptions, request/response formats, and examples.
-   **NFR4.2**: API error responses should be in JSON format and provide meaningful error messages.
-   **NFR4.3**: API endpoints requiring data input should use standard HTTP methods (POST, PUT) and expect data in JSON format.
-   **NFR4.4**: API endpoints for data retrieval should use the GET method.

#### 3.2.5 NFR5: Maintainability
-   **NFR5.1**: The codebase shall be organized into logical packages (`block`, `txpool`, `plane`, `api`, `ws`, `crosschain`).
-   **NFR5.2**: Code shall be commented, particularly for public functions, structs, and complex logic.
-   **NFR5.3**: Unit tests shall be provided for core functionalities in each package.
-   **NFR5.4**: The system shall use Go modules for dependency management.
-   **NFR5.5**: A `.golangci.yml` linting configuration should be provided and linters should pass.

#### 3.2.6 NFR6: Configurability
-   **NFR6.1**: The API server port shall be configurable via a `PORT` environment variable (e.g., from a `.env` file). Default to `8080`.
-   **NFR6.2**: Allowed CORS origins shall be configurable via a `CORS_ALLOWED_ORIGINS` environment variable. Default to `"*"` (allow all).

### 3.3 Interface Requirements

#### 3.3.1 API Interface
-   The system shall expose a RESTful HTTP API as detailed in Functional Requirement FR7.
-   Data exchange format for request and response bodies shall be JSON.

#### 3.3.2 WebSocket Interface
-   The system shall expose a WebSocket interface at `GET /ws` as detailed in Functional Requirement FR8.
-   Messages exchanged over WebSocket shall be in JSON format, following the `WSMessage` structure.

---

## 4. Other Requirements

### 4.1 Logging
-   **L1**: All API requests should be logged with method, path, and processing status/outcome.
-   **L2**: Significant events within modules (e.g., transaction added/removed from pool, block added to plane, client registered/unregistered from WebSocket hub) should be logged.
-   **L3**: Errors encountered during processing (e.g., signature verification failure, failed data validation, issues during broadcast) should be logged with relevant details.
-   **L4**: The system startup, including the port being used and environment variable loading status, should be logged.
-   **L5**: Cross-chain message simulation (send/receive) should be logged.

### 4.2 Error Handling
-   **EH1**: API errors should return appropriate HTTP status codes (e.g., 400 for bad request, 401 for unauthorized, 404 for not found, 415 for unsupported media type, 500 for internal server error).
-   **EH2**: API error responses should be in JSON format, e.g., `{"error": "message", "details": "optional_details"}`.
-   **EH3**: The global error handling middleware shall catch unhandled panics, log them, and return a generic 500 error response to the client if no response has been sent.

---

## 5. IEEE 830 Compliance Verification

| IEEE 830 Section                    | Corresponding SRS Section(s) | Notes                                                                 |
| :---------------------------------- | :--------------------------- | :-------------------------------------------------------------------- |
| **1. Introduction**                 |                              |                                                                       |
| 1.1 Purpose                         | 1.1 Purpose                  | AECH system purpose defined.                                          |
| 1.2 Scope                           | 1.2 Scope                    | Scope of AECH backend, features included/excluded.                    |
| 1.3 Definitions, Acronyms, Abbrev.  | 1.3 Definitions...           | Key terms and acronyms used in the document.                          |
| 1.4 References                      | 1.4 References               | Relevant documents and standards.                                     |
| 1.5 Overview                        | 1.5 Overview                 | Document structure outlined.                                          |
| **2. Overall Description**          |                              |                                                                       |
| 2.1 Product Perspective             | 2.1 Product Perspective      | AECH as a self-contained backend, its modularity.                     |
| 2.2 Product Functions               | 2.2 Product Functions        | Summary of key system functionalities.                                |
| 2.3 User Characteristics            | 2.3 User Characteristics     | Target users (developers, testers, students).                        |
| 2.4 Constraints                     | 2.4 Constraints              | Technical and operational constraints (Go, Gin, in-memory, etc.).     |
| 2.5 Assumptions and Dependencies    | 2.5 Assumptions...           | Key assumptions (client-side signing, single instance).               |
| **3. Specific Requirements**        |                              |                                                                       |
| 3.1 External Interface Requirements | 3.3 Interface Requirements   | Details API (HTTP/JSON) and WebSocket interfaces.                     |
| 3.2 Functional Requirements         | 3.1 Functional Requirements  | Detailed FRs for transactions, blocks, pool, plane, API, WS, etc.     |
| 3.3 Performance Requirements        | 3.2.1 NFR1: Performance      | Response times, broadcast initiation targets.                         |
| 3.4 Design Constraints              | 2.4 Constraints              | (Partially covered, also specific tech choices in FRs imply constraints) |
| 3.5 Attributes (Security, Maint.)   | 3.2.3 NFR3: Security         | Security NFRs (Content-Type, size limits, CORS, etc.).                |
|                                     | 3.2.5 NFR5: Maintainability  | Maintainability NFRs (modularity, comments, tests).                   |
| 3.6 Other Requirements              | 4. Other Requirements        | Logging (L1-L5) and Error Handling (EH1-EH3).                         |
| **4. Supporting Information**       |                              | (Appendices, Index - Not included in this SRS version)                |

*Note: This SRS aims for conceptual alignment with IEEE 830 structure. Some IEEE 830 sections like "Specific Requirements" are further broken down into Functional, Non-Functional, and Interface requirements as per common practice.*
```
