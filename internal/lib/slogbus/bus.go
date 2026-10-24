package slogbus

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
)

type Bus struct {
	nextRecordID     atomic.Uint64
	nextSubscriberID atomic.Uint64

	mu          sync.RWMutex
	subscribers map[uint64]chan Record
}

func New() *Bus {
	return &Bus{
		subscribers: make(map[uint64]chan Record),
	}
}

func (b *Bus) Handler(next slog.Handler) slog.Handler {
	return &handler{
		Handler: next,
		bus:     b,
	}
}

func (b *Bus) Watch(ctx context.Context) <-chan Record {
	id := b.nextSubscriberID.Add(1)
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

// Close closes the channels of all watchers, so their streams end and the API server can shut down.
func (b *Bus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for id, ch := range b.subscribers {
		close(ch)
		delete(b.subscribers, id)
	}
}

func (b *Bus) publish(record Record) {
	record.ID = b.nextRecordID.Add(1)

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
