package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
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

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:4242")
	if err != nil {
		fmt.Println("Couldn't connect to server!")
		return
	}

	cmd := exec.Command("/bin/sh")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("failed to make StdoutPipe!")
		return
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("failed to make StdinPipe!")
		return
	}

	cmd.Stderr = cmd.Stdout

	logs := []string{}
	render_flag := true

	// pipe os stdin
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			stdin.Write([]byte(text + "\n"))
			logs = append(logs, "[IN] "+text)
			render_flag = true
			message := ClientMessage{Content: text, Type: "stdin"}
			data, err := json.Marshal(message)
			if err != nil {
				fmt.Println("failed to encode message")
			}
			conn.Write(data)
		}
	}()

	// pipe connection stdin
	go func() {
		decoder := json.NewDecoder(conn)
		for {
			var m ServerMessage
			err := decoder.Decode(&m)
			if err == io.EOF {
				break
			} else if err != nil {
				break
			}

			stdin.Write([]byte(m.Content + "\n"))
			logs = append(logs, "[IN] "+m.Content)
			render_flag = true
		}
	}()

	// pipe command stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			text := scanner.Text()
			logs = append(logs, "[OUT] "+text)
			render_flag = true
			message := ClientMessage{Content: text, Type: "stdout"}
			data, err := json.Marshal(message)
			if err != nil {
				fmt.Println("failed to encode message")
			}
			conn.Write(data)
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Println("failed to start command!")
		return
	}

	go func() {
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
	}()

	cmd.Wait()
}
