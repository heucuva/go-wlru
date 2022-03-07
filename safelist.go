package wlru

import (
	"sync"
	"unsafe"
)

type safeList struct {
	// mu is a RWMutex to protect interactions with the list container
	mu sync.RWMutex
	// li is a double-ended list container
	li *embeddedList
	// once is an initialization gate
	once sync.Once
}

func (l *safeList) remove(n *node) {
	l.init()
	l.mu.Lock()
	defer l.mu.Unlock()

	l.li.Remove(n)
}

func (l *safeList) pushHead(n *node) {
	l.init()
	l.mu.Lock()
	defer l.mu.Unlock()

	l.li.InsertFirst(n)
}

func (l *safeList) tail() *node {
	l.init()
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.li.Last()
}

func (l *safeList) head() *node {
	l.init()
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.li.First()
}

func (l *safeList) prev(n *node) *node {
	l.init()
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.li.Prev(n)
}

func (l *safeList) next(n *node) *node {
	l.init()
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.li.Next(n)
}

func (l *safeList) init() {
	l.once.Do(func() {
		l.mu.Lock()
		defer l.mu.Unlock()
		var n node
		l.li = newEmbeddedList(unsafe.Offsetof(n.link))
	})
}
