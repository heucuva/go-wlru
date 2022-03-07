package wlru

import (
	"sync"
)

type safeList struct {
	// mu is a RWMutex to protect interactions with the list container
	mu sync.RWMutex
	// li is a double-ended list container
	li embeddedList
}

func (l *safeList) remove(n *node) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.li.Remove(n)
}

func (l *safeList) pushHead(n *node) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.li.InsertFirst(n)
}

func (l *safeList) tail() *node {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.li.Last()
}

func (l *safeList) head() *node {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.li.First()
}

func (l *safeList) prev(n *node) *node {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.li.Prev(n)
}

func (l *safeList) next(n *node) *node {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.li.Next(n)
}
