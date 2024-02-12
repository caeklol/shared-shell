package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	tm "github.com/buger/goterm"
)

type ServerMessage struct {
	Content string
}

type ClientMessage struct {
	Content string
	Type    string
}

var logs []string = []string{}
var render_flag bool = true
var connections int = 0
var conn net.Conn = nil

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4242")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer listener.Close()

	go render()
	go pipeOsStdin()

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

func pipeOsStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if conn == nil && connections == 0 {
			logs = append(logs, "[LOG] No client connected")
		} else {
			// -- possible race condition -- [A]
			logs = append(logs, "[IN] "+text)
			message := ServerMessage{Content: text}
			data, err := json.Marshal(message)
			if err != nil {
				fmt.Println("failed to encode message")
			}
			conn.Write(data)
		}

		render_flag = true
	}
}

func handleClient(c net.Conn) {
	defer handleDisconnect(conn)
	logs = append(logs, "[LOG] Client connected")
	render_flag = true
	connections += 1
	conn = c

	decoder := json.NewDecoder(c)
	for {
		var m ClientMessage
		err := decoder.Decode(&m)
		if err == io.EOF {
			break
		} else if err != nil {
			break
		}
		text := m.Content

		if m.Type == "stdout" {
			logs = append(logs, "[OUT] "+text)
		} else if m.Type == "stdin" {
			logs = append(logs, "[IN] "+text)
		} else {
			break
		}
	}
}

func handleDisconnect(conn net.Conn) {
	logs = append(logs, "[LOG] Client disconnected")
	render_flag = true
	connections -= 1
}

func render() {

	for {
		if !render_flag {
			time.Sleep(time.Millisecond)
			continue
		}

		distance_to_limit := len(logs) - tm.Height() + 1
		if distance_to_limit > 0 {
			for i := 0; i < distance_to_limit; i++ {
				logs = logs[1:]
			}
		}

		tm.Clear()
		tm.MoveCursor(0, 1)
		for _, s := range logs {
			tm.Println(s)
		}
		tm.MoveCursor(0, tm.Height())
		tm.Flush()
		render_flag = false
	}
}
