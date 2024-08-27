package event

import (
	"io"
	"sync"

	"github.com/gofrs/uuid"
)

type Subscription uuid.UUID

type Emitter[Topic, Msg any] interface {
	io.Closer

	Ready() bool
	Publish(topic Topic, msg Msg)
	Subscribe(topic Topic, ch chan Msg) Subscription
	Unsubscribe(topic Topic, sub Subscription)
}

type emitter[Topic comparable, Msg any] struct {
	mu     sync.RWMutex
	subs   map[Topic]map[Subscription]chan<- Msg
	closed bool
}

func (e *emitter[_, _]) Ready() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	empty := true
	for _, subs := range e.subs {
		if len(subs) > 0 {
			empty = false
		}

		for _, sub := range subs {
			if len(sub) == 0 {
				return true
			}
		}
	}

	return !empty
}

// Close implements Emitter.
func (e *emitter[_, _]) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return nil
	}

	e.closed = true
	for _, subs := range e.subs {
		for _, sub := range subs {
			close(sub)
		}
	}

	return nil
}

// Publish implements Emitter.
func (e *emitter[Topic, Msg]) Publish(topic Topic, msg Msg) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.closed {
		return
	}

	for _, ch := range e.subs[topic] {
		go func(ch chan<- Msg) {
			ch <- msg
		}(ch)
	}
}

// Subscribe implements Emitter.
func (e *emitter[Topic, Msg]) Subscribe(topic Topic, ch chan Msg) Subscription {
	e.mu.Lock()
	defer e.mu.Unlock()

	token, _ := uuid.NewV4()
	sub := Subscription(token)
	if e.subs[topic] == nil {
		e.subs[topic] = make(map[Subscription]chan<- Msg)
	}

	e.subs[topic][sub] = ch

	return sub
}

// Unsubscribe implements Emitter.
func (e *emitter[Topic, _]) Unsubscribe(topic Topic, sub Subscription) {
	e.mu.Lock()
	defer e.mu.Unlock()

	close(e.subs[topic][sub])
	delete(e.subs[topic], sub)
}

func NewEmitter[Topic comparable, Msg any]() Emitter[Topic, Msg] {
	return &emitter[Topic, Msg]{
		subs: make(map[Topic]map[Subscription]chan<- Msg),
	}
}
