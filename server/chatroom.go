package server

import (
	"bufio"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// chatroom is the collection of users which they can receive every message from each other
type chatroom struct {
	name  string
	users []*user
	mu    sync.Mutex
	his   *history
}

// newChatroom returns the pointer to a new chatroom struct
func newChatroom(name string) *chatroom {
	return &chatroom{
		name:  name,
		users: make([]*user, 0),
		his:   newHistory(10),
	}
}

// newUser adds the user to the chatroom and starts a loop for reading from the net.Conn
func (room *chatroom) newUser(user *user, conn net.Conn) {
	room.mu.Lock()
	room.users = append(room.users, user)
	room.mu.Unlock()
	room.broadcast(newSystemMsg(contentHello(user.name)))
	room.writeHistory(conn)
	for {
		reader := bufio.NewReader(conn)
		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			continue
		}
		content := strings.Trim(string(bytes), "\n")
		log.Printf("%s -> %s: %s\n", user.name, room.name, content)
		// if content equals to "exit" then close the connection
		if content == contentExit {
			user.done <- struct{}{}
			room.broadcast(newSystemMsg(contentGoodbye(user.name)))
			_ = conn.Close()
			break
		}
		msg := newUserMsg(user.name, content)
		room.his.push(msg)
		room.broadcast(msg)
	}
}

// broadcast sends the message to every user in the chatroom except the sender
func (room *chatroom) broadcast(msg message) {
	for _, u := range room.users {
		if u.name == msg.from {
			continue
		}
		go func(u *user) {
			select {
			case u.receive <- msg:
				break
			case <-time.After(3 * time.Second):
				break
			}
		}(u)
	}
}

// writeHistory writes the stored history messages to the connection
func (room *chatroom) writeHistory(conn net.Conn) {
	for _, msg := range room.his.get() {
		_, _ = conn.Write(msg.bytes())
	}
}
