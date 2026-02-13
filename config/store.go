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

// Get returns the current configuration (read lock)
func (s *Store) Get() interface{} {
s.mu.RLock()
defer s.mu.RUnlock()
return s.config
}

// Set replaces the entire configuration (write lock)
func (s *Store) Set(cfg interface{}) {
s.mu.Lock()
defer s.mu.Unlock()
s.config = cfg
}

// Update atomically updates the configuration using a function (write lock)
func (s *Store) Update(fn func(interface{}) interface{}) {
s.mu.Lock()
defer s.mu.Unlock()
s.config = fn(s.config)
}

// GetCopy returns a deep copy of the configuration via JSON marshaling
func (s *Store) GetCopy(dest interface{}) error {
s.mu.RLock()
defer s.mu.RUnlock()

data, err := json.Marshal(s.config)
if err != nil {
return err
}

return json.Unmarshal(data, dest)
}
