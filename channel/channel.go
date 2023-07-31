package channel

import "sync"

type Channel struct {
	buffer        []interface{}
	bufferSize    int
	senderMutex   sync.Mutex
	receiverMutex sync.Mutex
	senderCond    *sync.Cond
	receiverCond  *sync.Cond
	closed        bool
}

// NewChannel A factory function to create new channel
func NewChannel(bufferSize int) *Channel {
	ch := &Channel{
		buffer:     make([]interface{}, 0, bufferSize),
		bufferSize: bufferSize,
		closed:     false,
	}

	ch.senderCond = sync.NewCond(&ch.senderMutex)
	ch.receiverCond = sync.NewCond(&ch.receiverMutex)
	return ch
}

func (ch *Channel) Send(value interface{}) {
	ch.senderMutex.Lock()
	defer ch.senderMutex.Unlock()

	for len(ch.buffer) == ch.bufferSize && !ch.closed {
		ch.senderCond.Wait()
	}

	if ch.closed {
		panic("Sending on closed channel")
	}

	ch.buffer = append(ch.buffer, value)
	ch.receiverCond.Signal()

}

func (ch *Channel) Receive() (interface{}, bool) {
	ch.receiverMutex.Lock()
	defer ch.receiverMutex.Unlock()

	for len(ch.buffer) == 0 && !ch.closed {
		ch.receiverCond.Wait()
	}

	// cannot receive if channel is closed
	if len(ch.buffer) == 0 && ch.closed {
		return nil, false
	}

	value := ch.buffer[0]
	ch.buffer = ch.buffer[1:]
	ch.receiverCond.Signal()
	return value, true
}

func (ch *Channel) Close() {
	ch.senderMutex.Lock()
	defer ch.senderMutex.Unlock()

	ch.closed = true
	ch.senderCond.Broadcast()
	ch.receiverCond.Broadcast()
}
