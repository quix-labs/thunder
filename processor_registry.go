package thunder

import (
	"errors"
	"sync"
)

var (
	processors   = make(map[int]*Processor)
	processorsMu sync.RWMutex
)

func GetProcessors() map[int]*Processor {
	processorsMu.RLock()
	defer processorsMu.RUnlock()
	return processors
}

func GetProcessor(id int) (*Processor, error) {
	processorsMu.RLock()
	defer processorsMu.RUnlock()
	processor, exists := processors[id]
	if !exists {
		return nil, errors.New("processor not found")
	}
	return processor, nil
}

func AddProcessor(p *Processor) error {
	processorsMu.Lock()
	defer processorsMu.Unlock()

	// TODO UNIQUE ID INSTEAD OF INCREMENT
	// TODO ERROR ON DUPLICATE

	if p.ID == 0 {
		newId := 1
		for _, exists := processors[newId]; exists; _, exists = processors[newId] {
			newId++
		}
		p.ID = newId
	}

	processors[p.ID] = p
	GetBroadcaster().Dispatch("processor-created", p.ID)
	return nil
}

func UpdateProcessor(id int, p *Processor) error {
	processorsMu.Lock()
	defer processorsMu.Unlock()

	if id == 0 {
		return errors.New("cannot update processor with id 0")
	}
	if _, exists := processors[id]; !exists {
		return errors.New("processor not found")
	}

	p.ID = id
	processors[p.ID] = p
	GetBroadcaster().Dispatch("processor-updated", p.ID)

	return nil
}

func DeleteProcessor(id int) error {
	processorsMu.Lock()
	defer processorsMu.Unlock()
	_, exists := processors[id]
	if !exists {
		return errors.New("processor not found")
	}
	delete(processors, id)
	GetBroadcaster().Dispatch("processor-deleted", id)

	return nil
}
