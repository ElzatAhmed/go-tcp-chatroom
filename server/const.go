package server

import (
	"errors"
	"fmt"
)

var (
	errProtocolError = errors.New("connection protocol error")
)

const (
	contentExit = "exit"
)

func contentHello(name string) string {
	return fmt.Sprintf("user %s entered the chatroom and says hello!", name)
}

func contentGoodbye(name string) string {
	return fmt.Sprintf("user %s said goodbyte and exited the chatroom.", name)
}

const (
	msgTypeUser msgType = iota
	msgTypeSystem
)
