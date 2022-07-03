package main

import (
	"flag"
	"github.com/elzatahmed/go-tcp-chatroom/server"
)

func main() {
	host := flag.String("h", "127.0.0.1", "the host name of the server")
	port := flag.Int("p", 8888, "the port number of the server")
	flag.Parse()
	chatServer := server.New(*host, *port)
	chatServer.Spin()
}
