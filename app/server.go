package main

import (
	"fmt"
	"net"
	"os"
)


func main() {	
	fmt.Println("Starting server on port 6379")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	for {
		resp := NewResponse(conn)
		val, err := resp.Read()
		if err != nil {
			fmt.Println("Failed to read data:", err.Error())
			os.Exit(1)
		}

		fmt.Println("Received data:", val)
		conn.Write([]byte("+PONG\r\n"))
	}
}

// func handleRequest(conn net.Conn) {
// 	// Create a buffer to hold the incoming data
// 	buf := make([]byte, 1024)
// 	_, err := conn.Read(buf)
// 	if err != nil {
// 		if err == net.ErrClosed {
// 			return
// 		}
// 		fmt.Println("Failed to read data:", err.Error())
// 		os.Exit(1)
// 	}

// 	// Hardcode PONG response
// 	conn.Write([]byte("+OK\r\n"))

// 	fmt.Println("Received data:", string(buf))

// 	conn.Close()
// }