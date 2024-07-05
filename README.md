# Redis Client in Go

The goal of this project is to create a simple key-value store that implements using a subset of the Redis commands and data structures. The project is written in Go and uses the RESP for communication between the client and server.

## Objectives

-   [x] Implement a basic server that listens for connections on a port.
-   [x] Implement a client that can connect to the server.
-   [x] Implement a deserializer that can parse the Redis protocol.
-   [ ] Implement a serializer that can serialize data to the Redis protocol.
-   [ ] Implement a basic in-memory key-value store.
    -   [ ] Implement the `SET` command.
    -   [ ] Implement the `GET` command.
    -   [ ] Implement the `DEL` command.
    -   [ ] Implement the `EXISTS` command.
    -   [ ] Implement the `KEYS` command.
    -   [ ] Implement the `FLUSHDB` command.

## Usage

Make sure you have Go installed on your machine. You can download it from the [official website](https://golang.org/).

You should also have Redis installed, which you can download from the [official website](https://redis.io/).

To start the Redis server which connects to port 6379, run: `./start_redis_server.sh`. To use the Redis CLI, run `redis-cli` in a separate terminal.
