package utils

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type BroadcasterIn[T any] interface {
	Broadcast(value T) (emitted bool)
	Finish()
	Closed() bool
}

type broadcasterIn[T any] struct {
	_inputMu   sync.Mutex
	_inputChan chan T
	_closed    atomic.Bool
}

func (i *broadcasterIn[T]) Closed() bool {
	return i._closed.Load()
}

func (i *broadcasterIn[T]) Broadcast(value T) bool {
	i._inputMu.Lock()
	defer i._inputMu.Unlock()
	if i._closed.Load() {
		return false
	}
	i._inputChan <- value
	return true
}

func (i *broadcasterIn[T]) Finish() {
	if i._closed.Load() {
		return
	}
	i._closed.Store(true)

	// Close channel when event are not writing
	i._inputMu.Lock()
	close(i._inputChan)
	i._inputMu.Unlock()
}

type BroadcastProcessor[T any, V any] func(in T) (out V)

type Broadcaster[T any, V any] struct {
	in        *broadcasterIn[T]
	processor BroadcastProcessor[T, V]

	listenerCtx    context.Context
	listenerCancel context.CancelFunc

	listenChannels sync.Map

	_once sync.Once
}

func (b *Broadcaster[T, V]) NewListenChan() (<-chan V, func()) {
	listenChan := make(chan V)

	b.listenChannels.Store(listenChan, true)
	return listenChan, func() {
		b.listenChannels.Delete(listenChan)
		close(listenChan)
	}
}

func (b *Broadcaster[T, V]) In() BroadcasterIn[T] {
	return b.in
}

func (b *Broadcaster[T, V]) Start() {
	go func() {
		for data := range b.in._inputChan {
			transformedData := b.processor(data)
			var wg sync.WaitGroup
			b.listenChannels.Range(func(key, value any) bool {
				wg.Add(1)
				//TODO USE MUTEX INSTEAD and avoid recover
				go func(key any) {
					defer wg.Done()
					defer func() {
						if r := recover(); r != nil {
						}
					}()

					if ch, ok := key.(chan V); ok {
						ch <- transformedData
					}
				}(key)
				return true
			})
			wg.Wait()
		}
		b.listenerCancel()
	}()
}

func (b *Broadcaster[T, V]) Close() {
	b._once.Do(func() {
		b.in.Finish()
		b.Wait()
	})
}

func (b *Broadcaster[T, V]) Wait() {

	// Wait for input end
	for !b.in.Closed() {
		time.Sleep(time.Millisecond * 10)
	}

	// Wait for targets end
	<-b.listenerCtx.Done()
}

func NewBroadcaster[T any, V any](processor BroadcastProcessor[T, V]) *Broadcaster[T, V] {
	in := &broadcasterIn[T]{_inputChan: make(chan T)}
	ctx, cancel := context.WithCancel(context.Background())
	return &Broadcaster[T, V]{in: in, listenerCtx: ctx, listenerCancel: cancel, processor: processor}
}
