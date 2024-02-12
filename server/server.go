package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

type Message struct {
	Content string
	Type    string
}

var connections int = 0

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4242")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer listener.Close()

	fmt.Println("Server is listening on port 4242")

	for {
		conn, err := listener.Accept()
		if connections >= 1 {
			fmt.Println("Recieved duplicate connection")
			continue
		}

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer handleDisconnect(conn)
	fmt.Println("Client connected")
	connections += 1

	decoder := json.NewDecoder(conn)
	for {
		var m Message
		err := decoder.Decode(&m)
		if err == io.EOF {
			break
		} else if err != nil {
			break
		}

		if m.Type == "stdout" {
			fmt.Printf("[OUT] %s", m.Content)
		} else if m.Type == "stdin" {
			fmt.Printf("[IN] %s", m.Content)
		} else {
			break
		}
	}
}

func handleDisconnect(conn net.Conn) {
	fmt.Println("Client disconnected")
	connections -= 1
}
