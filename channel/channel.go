package channel

import (
	"sync"
)

type Channel struct {
	buffer     []interface{}
	bufferSize int
	Mutex      sync.Mutex
	Cond       *sync.Cond
	closed     bool
}

// NewChannel A factory function to create new channel
func NewChannel(bufferSize int) *Channel {
	ch := &Channel{
		buffer:     make([]interface{}, 0, bufferSize),
		bufferSize: bufferSize,
		closed:     false,
	}

	ch.Cond = sync.NewCond(&ch.Mutex)
	return ch
}

func (ch *Channel) Send(value interface{}) {
	ch.Mutex.Lock()
	defer ch.Mutex.Unlock()

	for len(ch.buffer) == ch.bufferSize && !ch.closed {
		ch.Cond.Wait()
	}

	if ch.closed {
		panic("Sending on closed channel")
	}

	ch.buffer = append(ch.buffer, value)
	ch.Cond.Signal()

}

func (ch *Channel) Receive() (interface{}, bool) {
	ch.Mutex.Lock()
	defer ch.Mutex.Unlock()

	for len(ch.buffer) == 0 && !ch.closed {
		ch.Cond.Wait()
	}

	if len(ch.buffer) == 0 && ch.closed {
		return nil, false
	}

	value := ch.buffer[0]
	ch.buffer = ch.buffer[1:]
	ch.Cond.Signal()
	return value, true
}

func (ch *Channel) Close() {
	ch.Mutex.Lock()
	defer ch.Mutex.Unlock()

	ch.closed = true
	ch.Cond.Broadcast()
}
