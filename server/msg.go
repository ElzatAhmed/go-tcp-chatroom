package server

import (
	"fmt"
	"time"
)

type msgType uint8

// message is the model for every message flows through the chatroom
type message struct {
	typ     msgType
	from    string
	content string
	when    time.Time
}

func newUserMsg(from string, content string) message {
	return message{
		typ:     msgTypeUser,
		from:    from,
		content: content,
		when:    time.Now(),
	}
}

func newSystemMsg(content string) message {
	return message{
		typ:     msgTypeSystem,
		from:    "",
		content: content,
		when:    time.Now(),
	}
}

func (msg message) string() string {
	switch msg.typ {
	case msgTypeUser:
		return fmt.Sprintf("%s %s: %s\n",
			msg.when.Format("2006-01-02 15:04:05"), msg.from, msg.content)
	case msgTypeSystem:
		return fmt.Sprintf("%s server: %s\n",
			msg.when.Format("2006-01-02 15:04:05"), msg.content)
	}
	return ""
}

func (msg message) bytes() []byte {
	return []byte(msg.string())
}
