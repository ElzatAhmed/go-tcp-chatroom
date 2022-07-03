package client

import (
	"bufio"
	"fmt"
	"github.com/elzatahmed/go-tcp-chatroom/comm"
	"log"
	"net"
	"os"
	"sync"
)

// client is the model for a client type
type client struct {
	serverAddr string
	username   string
	roomname   string
}

// New returns the pointer to a new client struct
func New(host string, port int, username string, roomname string) *client {
	return &client{
		serverAddr: fmt.Sprintf("%s:%d", host, port),
		username:   username,
		roomname:   roomname,
	}
}

// Spin connects to the chat server and starts the chatting
func (cli *client) Spin() {
	conn, err := net.Dial("tcp", cli.serverAddr)
	if err != nil {
		log.Fatalf("server connection failed with err: %s\n", err.Error())
	}
	_, err = conn.Write(cli.protocol())
	if err != nil {
		_ = conn.Close()
		log.Fatalf("send protocol failed with err: %s\n", err.Error())
	}
	reader := bufio.NewReader(conn)
	bytes, err := reader.ReadBytes('\n')
	if err != nil {
		_ = conn.Close()
		log.Fatalf("read response from the server failed with err: %s\n", err.Error())
	}

	res := string(bytes)
	if res == string(comm.BytesProtocolErr) || res == string(comm.BytesUsernameExists) {
		_ = conn.Close()
		log.Fatalf(res)
	}
	fmt.Print(res)
	done := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(2)
	go cli.read(reader, done, &wg)
	go cli.write(conn, done, &wg)
	wg.Wait()
	_ = conn.Close()
}

func (cli *client) protocol() []byte {
	return []byte(fmt.Sprintf("%s;%s\n", cli.username, cli.roomname))
}

func (cli *client) read(reader *bufio.Reader, done <-chan interface{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-done:
			wg.Done()
			break
		default:
			bytes, err := reader.ReadBytes('\n')
			if err != nil {
				continue
			}
			fmt.Print(string(bytes))
		}
	}
}

func (cli *client) write(conn net.Conn, done chan<- interface{}, wg *sync.WaitGroup) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		text := scanner.Text()
		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			log.Printf("write text `%s` failed with err: %s\n", text, err.Error())
			continue
		}
		if text == "exit" {
			log.Println("closing the connection with the chat server...")
			done <- struct{}{}
			log.Println("exiting the program...")
			break
		}
	}
	wg.Done()
}
