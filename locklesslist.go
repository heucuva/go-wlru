package wlru

import (
	"sync/atomic"
	"unsafe"
)

type locklessList struct {
	h unsafe.Pointer // *node
	t unsafe.Pointer // *node
}

func (l *locklessList) head() *node {
	return (*node)(atomic.LoadPointer(&l.h))
}

func (l *locklessList) tail() *node {
	return (*node)(atomic.LoadPointer(&l.t))
}

func (l *locklessList) pushHead(n *node) {
	newHeadPtr := unsafe.Pointer(n)
	atomic.StorePointer(&n.left, nil)
	for {
		oldHead := l.head()
		oldHeadPtr := unsafe.Pointer(oldHead)
		if oldHeadPtr != newHeadPtr {
			atomic.StorePointer(&n.right, oldHeadPtr)
		}
		if atomic.CompareAndSwapPointer(&l.h, oldHeadPtr, newHeadPtr) {
			// we've completed the swap!
			if oldHead != nil && oldHeadPtr != newHeadPtr {
				atomic.CompareAndSwapPointer(&oldHead.left, nil, newHeadPtr)
			}
			atomic.CompareAndSwapPointer(&l.t, nil, newHeadPtr)
			return
		}
	}
}

func (l *locklessList) remove(n *node) {
	curNode := unsafe.Pointer(n)
	left, right := n.prev(), n.next()
	if left != nil {
		atomic.CompareAndSwapPointer(&left.right, curNode, unsafe.Pointer(right))
	}
	if right != nil {
		atomic.CompareAndSwapPointer(&right.left, curNode, unsafe.Pointer(left))
	}
	atomic.CompareAndSwapPointer(&l.h, curNode, unsafe.Pointer(right))
	atomic.CompareAndSwapPointer(&l.t, curNode, unsafe.Pointer(left))
	atomic.StorePointer(&n.left, nil)
	atomic.StorePointer(&n.right, nil)
}
