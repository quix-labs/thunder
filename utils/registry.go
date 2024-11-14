package utils

import (
	"fmt"
	"sync"
)

type Registrable interface{}

type RegistryValidateCallback[T Registrable] func(ID string, item T) error
type RegistryIdGenerator[T Registrable] func(item *T) (string, error)

type Registry[T Registrable] struct {
	mu    sync.RWMutex
	items map[string]T

	_entityName    string
	_idGenerator   RegistryIdGenerator[T]
	_validateUsing RegistryValidateCallback[T]
}

func NewRegistry[T Registrable](entityName string) *Registry[T] {
	if entityName == "" {
		panic("please provide entity name to create registry")
	}

	return &Registry[T]{
		mu:             sync.RWMutex{},
		items:          make(map[string]T),
		_entityName:    entityName,
		_validateUsing: nil,
	}
}

func (r *Registry[T]) SetIdGenerator(generator RegistryIdGenerator[T]) *Registry[T] {
	r._idGenerator = generator
	return r
}
func (r *Registry[T]) ValidateUsing(callback RegistryValidateCallback[T]) *Registry[T] {
	r._validateUsing = callback
	return r
}

func (r *Registry[T]) Register(ID string, item T) error {
	if ID == "" && r._idGenerator == nil {
		return fmt.Errorf("provide an id to register %s", r._entityName)
	}

	if ID == "" && r._idGenerator != nil {
		var err error
		if ID, err = r._idGenerator(&item); err != nil {
			return fmt.Errorf("failed to generate an id for %s", r._entityName)
		}
	}

	// Validate if callback specified
	if r._validateUsing != nil {
		if err := r._validateUsing(ID, item); err != nil {
			return err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[ID]; ok {
		return fmt.Errorf("%s already registered: %s", r._entityName, ID)
	}
	r.items[ID] = item
	return nil
}

func (r *Registry[T]) Update(ID string, item T) error {
	if _, err := r.Get(ID); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[ID] = item
	return nil
}

func (r *Registry[T]) Delete(ID string) error {
	if _, err := r.Get(ID); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, ID)
	return nil
}

func (r *Registry[T]) Get(ID string) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[ID]
	if !ok {
		return *new(T), fmt.Errorf("%s with ID: %s not found", r._entityName, ID)
	}

	return item, nil
}

func (r *Registry[T]) All() map[string]T {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.items
}

/** EVENT HANDLING **/

func (r *Registry[T]) AfterUpdated(callback RegistryValidateCallback[T]) *Registry[T] {
	//TODO
	return r
}

func (r *Registry[T]) AfterDeleted(callback RegistryValidateCallback[T]) *Registry[T] {
	//TODO
	return r
}

func (r *Registry[T]) AfterRegistered(callback RegistryValidateCallback[T]) *Registry[T] {
	//TODO
	return r
}
