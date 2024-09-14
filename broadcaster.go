package thunder

type BroadcasterAllCallback func(event string, data any)
type BroadcasterEventCallback func(data any)
type BroadCasterOffFunction func()

type Broadcaster struct {
	listeners       map[string][]*BroadcasterEventCallback
	globalListeners []*BroadcasterAllCallback
}

func (b *Broadcaster) On(event string, callback BroadcasterEventCallback) BroadCasterOffFunction {
	if b.listeners[event] == nil {
		b.listeners[event] = make([]*BroadcasterEventCallback, 0)
	}
	b.listeners[event] = append(b.listeners[event], &callback)

	return func() {
		if listeners, ok := b.listeners[event]; ok {
			for i, cb := range listeners {
				if cb == &callback {
					b.listeners[event] = append(listeners[:i], listeners[i+1:]...)
					break
				}
			}
		}
	}
}

func (b *Broadcaster) OnAll(callback BroadcasterAllCallback) BroadCasterOffFunction {
	b.globalListeners = append(b.globalListeners, &callback)
	return func() {
		for i, cb := range b.globalListeners {
			if cb == &callback {
				b.globalListeners = append(b.globalListeners[:i], b.globalListeners[i+1:]...)
				break
			}
		}
	}
}

func (b *Broadcaster) Dispatch(event string, data any) {
	if listeners, ok := b.listeners[event]; ok {
		for _, callback := range listeners {
			(*callback)(data)
		}
	}
	for _, callback := range b.globalListeners {
		(*callback)(event, data)
	}
}

var (
	broadcaster = &Broadcaster{
		listeners:       make(map[string][]*BroadcasterEventCallback),
		globalListeners: make([]*BroadcasterAllCallback, 0),
	}
)

func GetBroadcaster() *Broadcaster {
	return broadcaster
}
