package sync

import "sync"

type KeyWaiter[T comparable] struct {
	mainMutex  *sync.Mutex
	keyMutexes map[T]*sync.Mutex
}

func NewWaiter[T comparable]() *KeyWaiter[T] {
	return &KeyWaiter[T]{
		mainMutex:  &sync.Mutex{},
		keyMutexes: make(map[T]*sync.Mutex),
	}
}

func (w *KeyWaiter[T]) Acquire(key T) {
	w.mainMutex.Lock()
	defer w.mainMutex.Unlock()

	foundMutex := w.keyMutexes[key]
	if foundMutex == nil {
		newMutex := &sync.Mutex{}
		newMutex.Lock()
		w.keyMutexes[key] = newMutex
		return
	}

	foundMutex.Lock()
}

func (w *KeyWaiter[T]) Release(key T) {
	w.mainMutex.Lock()
	defer w.mainMutex.Unlock()

	foundMutex := w.keyMutexes[key]
	if foundMutex == nil {
		panic("release of a non-existent key")
	}

	foundMutex.Unlock()
}
