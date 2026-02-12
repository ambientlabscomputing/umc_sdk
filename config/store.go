package config
package config

import (
	"encoding/json"
	"sync"
)

// Store provides thread-safe configuration storage
type Store struct {
	mu     sync.RWMutex
	config interface{}
}

// NewStore creates a new configuration store
func NewStore(initial interface{}) *Store {
	return &Store{
		config: initial,
	}
}


































}	return json.Unmarshal(data, dest)		}		return err	if err != nil {	data, err := json.Marshal(s.config)		defer s.mu.RUnlock()	s.mu.RLock()func (s *Store) GetCopy(dest interface{}) error {// GetCopy returns a deep copy of the configuration via JSON marshaling}	s.config = fn(s.config)	defer s.mu.Unlock()	s.mu.Lock()func (s *Store) Update(fn func(interface{}) interface{}) {// Update atomically updates the configuration using a function (write lock)}	s.config = cfg	defer s.mu.Unlock()	s.mu.Lock()func (s *Store) Set(cfg interface{}) {// Set replaces the entire configuration (write lock)}	return s.config	defer s.mu.RUnlock()	s.mu.RLock()func (s *Store) Get() interface{} {// Get returns the current configuration (read lock)