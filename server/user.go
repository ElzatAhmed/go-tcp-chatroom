package server

import "net"

// user is a single connection to the tcp chat server
type user struct {
	name    string
	receive chan message
	done    chan interface{}
}

// newUser returns the pointer to a new user struct
func newUser(name string) *user {
	return &user{
		name:    name,
		receive: make(chan message),
	}
}

// listen starts a loop to receive from the receive channel and writes to the net.Conn
func (u *user) listen(conn net.Conn) {
	for {
		select {
		case msg := <-u.receive:
			_, _ = conn.Write(msg.bytes())
		case <-u.done:
			break
		}
	}
}
