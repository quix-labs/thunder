package thunder

import (
	"errors"
	"sync"
)

var (
	targets   = make(map[int]*Target)
	targetsMu sync.RWMutex
)

func GetTargets() map[int]*Target {
	targetsMu.RLock()
	defer targetsMu.RUnlock()
	return targets
}

func GetTarget(id int) (*Target, error) {
	targetsMu.RLock()
	defer targetsMu.RUnlock()
	target, exists := targets[id]
	if !exists {
		return nil, errors.New("target not found")
	}
	return target, nil
}

func AddTarget(t *Target) error {
	targetsMu.Lock()
	defer targetsMu.Unlock()

	// TODO UNIQUE ID INSTEAD OF INCREMENT
	// TODO ERROR ON DUPLICATE

	if t.ID == 0 {
		newId := len(targets) + 1
		t.ID = newId
	}

	targets[t.ID] = t
	return nil
}

func UpdateTarget(id int, t *Target) error {
	targetsMu.Lock()
	defer targetsMu.Unlock()

	if id == 0 {
		return errors.New("cannot update target with id 0")
	}
	if _, exists := sources[id]; !exists {
		return errors.New("target not found")
	}

	t.ID = id
	targets[t.ID] = t

	return nil
}

func DeleteTarget(id int) error {
	targetsMu.Lock()
	defer targetsMu.Unlock()
	_, exists := targets[id]
	if !exists {
		return errors.New("target not found")
	}
	delete(targets, id)
	return nil
}
