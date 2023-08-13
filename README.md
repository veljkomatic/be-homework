# Backend Homework - Tx Parser

## Goal

Implement Ethereum blockchain parser that will allow to query transactions for subscribed addresses.

## Problem

Users not able to receive push notifications for incoming/outgoing transactions. By Implementing `Parser` interface we would be able to hook this up to notifications service to notify about any incoming/outgoing transactions.

Expose public interface for external usage either via code or command line or rest api that will include supported list of operations defined in the Parser interface

    type Parser interface {
        // last parsed block
        GetCurrentBlock() int
	    // add address to observer
	    Subscribe(address string) bool
	    // list of inbound or outbound transactions for an address
	    GetTransactions(address string) []Transaction
    }

External usage exposed via REST API

    curl -X GET http://localhost:8080/block-number // get last parsed block
    curl -X POST -d '{"address": "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5"}' http://localhost:8080/subscribe // subscribe to address
    curl -X GET http://localhost:8080/transactions/:address // get transactions for address

# Code structure
## cmd directory
The cmd directory is commonly used in Go projects to represent the entry points of the application,
This is particularly useful in the case of larger systems like microservices where you might have multiple services within the same git repository.
Currently, we only parse-service in the cmd directory, but we can add more services in the future.

### parse-service
In main.go, init application and start processing new blocks from the blockchain and start the rest server.
In internal directory, we have the main logic of the application, including:
- block_parser: parse new blocks from the blockchain and send them to the channel, here we start processing from last block number. When starting default block number is 0.
- transaction_filter: filter transactions from the block for observed addresses and store them in storage(in memory). Trade off here we filter all transactions of block synchronously, but we can do it in parallel in the future.
    - Note: here after finding transaction for observed address, we could send it to notification service via message broker. Notification service is not implemented, because it is not in the scope of the task.
- server : rest server to expose the API for the client.

## common directory
In the common directory, we have the common logic of the application
- json-rpc: json-rpc request and response models

## pkg directory
The pkg directory is used to hold libraries and code that's intended to be used by other services.
- blockchain:
    - block: block model represents the block in the blockchain with transactions
    - types: block number and conversion functions
- parser: parser interface and implementation, this is given interface from the task. Note, I added context as first argument to the methods, its golang good practice to provide context to the methods.
- provider: rpc provider interface and implementation, rpc url is cloudflare-eth endpoint, but we can add more providers in the future.
- storage:
    - block: block storage and repository, here we store the last block number
    - transaction: transaction storage and repository, here we store transactions for observed addresses
- subscriber:
    - filter: filter transactions from the block for observed addresses
    - subscriber: subscribe to addresses and store them in storage(in memory)

## TODOs in the future
- Add validations
  - for the request body (provided address)
- Cover code with unit and integration tests
  - tests processing new blocks in parallel for given rage, check that all blocks in the rage are processed
  - test for error handling, ensure that block that failed to process will be sent to the channel and will be processed again
  - test transaction filtering for observed addresses, set observed addresses and check that only transactions for these addresses are stored in storage
  - test in memory storage for transactions and blocks, store and then get transactions and blocks from storage to ensure that they are stored correctly
  - test rest server, send requests to the server and check that responses are correct
- Use different provider in the future
  - add more providers in the pkg/provider directory
  - use for example bloxroute provider instead of cloudflare-eth, because it gives you ability to subscribe to new blocks via websocket or grpc
  - now we only support eth mainnet, but we can add more networks in the future
  - handle reorgs, now we only process new blocks, but we can add reorgs handling in the future

### Architecture
- Switch in-memory storage to the database
- Add more service for example backfilling-service, this service will be responsible for backfilling blocks from the blockchain and store them in the database. For this service, we could use clickhouse to store transactions from block history.
- Then we can add spin up some workers to process old blocks from the clickhouse for newly subscribed addresses.
- Separate rest server from parser-service, and redesign it to be api gateway for the application.
- block_parser, transaction_filter will be separate services, and they will be responsible for processing new blocks and filtering transactions for observed addresses.
- Instead of using channels, use message broker for example kafka, and redesign the application to be event-driven.
