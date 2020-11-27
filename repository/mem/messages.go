package mem

import (
	"sync"

	"github.com/ayzatziko/chat/repository"
)

var _ repository.Messages = (*LimitMessages)(nil)

type LimitMessages struct {
	N, ptr int

	mu  sync.Mutex
	buf []*repository.Message
}

func (m *LimitMessages) Add(msg *repository.Message) {
	m.mu.Lock()
	if m.buf == nil {
		m.buf = make([]*repository.Message, 0, m.N)
	}

	if len(m.buf) < cap(m.buf) {
		m.buf = append(m.buf, msg)
	} else {
		m.buf[m.ptr] = msg
	}

	m.ptr++
	if m.ptr > m.N-1 {
		m.ptr = 0
	}
	m.mu.Unlock()
}

func (m *LimitMessages) Messages() []*repository.Message {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.buf == nil {
		return nil
	}

	messages := make([]*repository.Message, len(m.buf))
	copy(messages, m.buf[m.ptr:])
	copy(messages[m.ptr:], m.buf[:m.ptr])
	return messages
}
