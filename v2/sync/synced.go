package sync

import "sync"

type Synced[T any] struct {
	mutex *sync.Mutex
	value T
}

func NewSynced[T any]() *Synced[T] {
	return &Synced[T]{
		mutex: &sync.Mutex{},
	}
}

func (s *Synced[T]) Get() T {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.value
}

func (s *Synced[T]) Set(value T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.value = value
}
