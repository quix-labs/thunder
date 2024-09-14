package thunder

import (
	"errors"
	"sync"
)

var (
	sources   = make(map[int]*Source)
	sourcesMu sync.RWMutex
)

func GetSources() map[int]*Source {
	sourcesMu.RLock()
	defer sourcesMu.RUnlock()
	return sources
}

func GetSource(id int) (*Source, error) {
	sourcesMu.RLock()
	defer sourcesMu.RUnlock()
	source, exists := sources[id]
	if !exists {
		return nil, errors.New("source not found")
	}
	return source, nil
}

func AddSource(s *Source) error {
	sourcesMu.Lock()
	defer sourcesMu.Unlock()

	// TODO UNIQUE ID INSTEAD OF INCREMENT
	// TODO ERROR ON DUPLICATE
	if s.ID == 0 {
		newId := len(sources) + 1
		s.ID = newId
	}

	sources[s.ID] = s

	return nil
}

func UpdateSource(id int, s *Source) error {
	sourcesMu.Lock()
	defer sourcesMu.Unlock()

	if id == 0 {
		return errors.New("cannot update source with id 0")
	}
	if _, exists := sources[id]; !exists {
		return errors.New("source not found")
	}

	s.ID = id
	sources[s.ID] = s

	return nil
}

func DeleteSource(id int) error {
	sourcesMu.Lock()
	defer sourcesMu.Unlock()
	_, exists := sources[id]
	if !exists {
		return errors.New("source not found")
	}
	delete(sources, id)
	return nil
}
