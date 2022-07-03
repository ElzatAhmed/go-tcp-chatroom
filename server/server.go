package server

import (
	"bufio"
	"fmt"
	"github.com/elzatahmed/go-tcp-chatroom/comm"
	"log"
	"net"
	"sync"
)

// chatServer is the listening and dispatching server for tcp chatroom,
// it stores information about all the rooms and all the users ever created
type chatServer struct {
	addr   string
	rooms  map[string]*chatroom
	users  map[string]*user
	muRoom sync.Mutex
	muUser sync.Mutex
}

// New returns the pointer to a new chatServer struct
func New(host string, port int) *chatServer {
	return &chatServer{
		addr:  fmt.Sprintf("%s:%d", host, port),
		rooms: make(map[string]*chatroom),
		users: make(map[string]*user),
	}
}

// Spin starts the chatServer at given address
func (server *chatServer) Spin() {
	listener, err := net.Listen("tcp", server.addr)
	if err != nil {
		log.Fatalf("failed to start the server at %s, err: %s\n", server.addr, err.Error())
	}
	log.Printf("server started at address %s...\n", server.addr)
	for {
		conn, err := listener.Accept()
		log.Printf("server accepted a new connection from %s\n", conn.RemoteAddr())
		if err != nil {
			continue
		}
		go server.spin(conn)
	}
}

// spin do the protocol procedure and starts the connection goroutines
func (server *chatServer) spin(conn net.Conn) {
	reader := bufio.NewReader(conn)
	bytes, err := reader.ReadBytes('\n')
	if err != nil {
		log.Printf("connection failed with client %s with err: %s\n",
			conn.RemoteAddr(), err.Error())
		return
	}
	username, roomname, err := parseProtocol(bytes)
	if err != nil {
		_, _ = conn.Write(comm.BytesProtocolErr)
		return
	}
	if _, ok := server.users[username]; ok {
		_, _ = conn.Write(comm.BytesUsernameExists)
		_ = conn.Close()
		log.Printf("connection from %s closed by server\n", conn.RemoteAddr())
		return
	}
	log.Printf("connecting user %s to chatroom %s...\n", username, roomname)
	u := server.newUser(username)
	room, ok := server.rooms[roomname]

	if !ok {
		room = server.newRoom(roomname)
	}
	go room.newUser(u, conn)
	go u.listen(conn)
	log.Printf("user %s is connected to chatroom %s\n", username, roomname)
}

// newRoom constructs a new chatroom, adds it to the rooms map and returns it
func (server *chatServer) newRoom(name string) *chatroom {
	server.muRoom.Lock()
	defer server.muRoom.Unlock()
	room := newChatroom(name)
	server.rooms[name] = room
	return room
}

// newUser constructs a new user, adds it to the users map and returns it
func (server *chatServer) newUser(name string) *user {
	server.muUser.Lock()
	defer server.muUser.Unlock()
	u := newUser(name)
	server.users[name] = u
	return u
}
