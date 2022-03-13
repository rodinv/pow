package storage

import (
	"sync"
)

// Storage is a hash storage for hashcash algorithm
type Storage struct {
	m sync.Map
}

func New() *Storage {
	return &Storage{
		m: sync.Map{},
	}
}

// IsSpent if not exists - saves value and return false, if exists - return true
func (s *Storage) IsSpent(key string) bool {
	_, loaded := s.m.LoadOrStore(key, struct{}{})
	return loaded
}
