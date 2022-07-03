package server

import "sync"

// history stores history messages up to maxCount
type history struct {
	messages []message
	maxCount int
	curCount int
	mu       sync.Mutex
}

// newHistory returns the pointer to a new history struct
func newHistory(maxCount int) *history {
	return &history{
		messages: make([]message, maxCount),
		maxCount: maxCount,
		curCount: 0,
	}
}

// push appends a new message into the history,
// if message count is greater than maxCount then the first in message will be obsolete
func (h *history) push(msg message) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.curCount == h.maxCount {
		h.messages = h.messages[1:]
		h.messages[9] = msg
		return
	}
	h.messages[h.curCount] = msg
	h.curCount++
}

func (h *history) get() []message {
	return h.messages[:h.curCount]
}
