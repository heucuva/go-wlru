package wlru

import (
	"sync"
	"sync/atomic"
)

type safeMap struct {
	e sync.Map
	c int64
}

func (m *safeMap) set(n *node) {
	if _, found := m.e.LoadOrStore(n.key, n); !found {
		atomic.AddInt64(&m.c, 1)
	}
}

func (m *safeMap) size() int {
	return int(atomic.LoadInt64(&m.c))
}

func (m *safeMap) remove(key interface{}) (*node, bool) {
	lhs, found := m.e.LoadAndDelete(key)
	if !found {
		return nil, false
	}
	atomic.AddInt64(&m.c, -1)
	return lhs.(*node), true
}

func (m *safeMap) get(key interface{}) (*node, bool) {
	lhs, found := m.e.Load(key)
	if !found {
		return nil, false
	}
	return lhs.(*node), true
}

func (m *safeMap) clear() {
	m.e = sync.Map{}
	atomic.StoreInt64(&m.c, 0)
}
