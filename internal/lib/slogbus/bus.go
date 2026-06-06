package slogbus

import (
	"context"
	"sync"
	"sync/atomic"
)

type Bus struct {
	nextRecordID     uint64
	nextSubscriberID uint64

	mu          sync.RWMutex
	subscribers map[uint64]chan Record
}

func New() *Bus {
	return &Bus{
		subscribers: make(map[uint64]chan Record),
	}
}

func (b *Bus) Watch(ctx context.Context) <-chan Record {
	id := atomic.AddUint64(&b.nextSubscriberID, 1)
	ch := make(chan Record, 256)

	b.mu.Lock()
	b.subscribers[id] = ch
	b.mu.Unlock()

	go func() {
		<-ctx.Done()
		b.mu.Lock()
		delete(b.subscribers, id)
		b.mu.Unlock()
	}()

	return ch
}

func (b *Bus) publish(record Record) {
	record.ID = atomic.AddUint64(&b.nextRecordID, 1)

	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.subscribers {
		select {
		case ch <- record:
		default:
			// Keep logging non-blocking if a subscriber falls behind.
		}
	}
}
