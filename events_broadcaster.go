package thunder

type EventsBroadcasterAllCallback func(event string, data any)
type EventsBroadcasterEventCallback func(data any)
type EventsBroadCasterOffFunction func()

type EventsBroadcaster struct {
	listeners       map[string][]*EventsBroadcasterEventCallback
	globalListeners []*EventsBroadcasterAllCallback
}

func (b *EventsBroadcaster) On(event string, callback EventsBroadcasterEventCallback) EventsBroadCasterOffFunction {
	if b.listeners[event] == nil {
		b.listeners[event] = make([]*EventsBroadcasterEventCallback, 0)
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

func (b *EventsBroadcaster) OnAll(callback EventsBroadcasterAllCallback) EventsBroadCasterOffFunction {
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

func (b *EventsBroadcaster) Dispatch(event string, data any) {
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
	broadcaster = &EventsBroadcaster{
		listeners:       make(map[string][]*EventsBroadcasterEventCallback),
		globalListeners: make([]*EventsBroadcasterAllCallback, 0),
	}
)

func GetEventBroadcaster() *EventsBroadcaster {
	return broadcaster
}
