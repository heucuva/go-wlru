package wlru

import (
	"unsafe"
)

// embeddedListLink is a link to the list container
type embeddedListLink struct {
	prev *node
	next *node
}

func (l *embeddedListLink) remove(linkFieldOfs uintptr, head **node, tail **node) bool {
	if !l.isContained(linkFieldOfs, *head) {
		return false
	}
	if l.prev == nil {
		*head = l.next
	} else {
		getListLink(unsafe.Pointer(l.prev), linkFieldOfs).next = l.next
	}
	if l.next == nil {
		*tail = l.prev
	} else {
		getListLink(unsafe.Pointer(l.next), linkFieldOfs).prev = l.prev
	}

	l.next = nil
	l.prev = nil
	return true
}

func (l *embeddedListLink) isContained(linkFieldOfs uintptr, head *node) bool {
	return l.prev != nil || head == l.getItem(linkFieldOfs)
}

func (l *embeddedListLink) getItem(linkFieldOfs uintptr) *node {
	u := unsafe.Add(unsafe.Pointer(l), (^linkFieldOfs)+1)
	m := (*node)(u)
	return m
}

func getListLink(obj unsafe.Pointer, linkFieldOfs uintptr) *embeddedListLink {
	return (*embeddedListLink)(unsafe.Add(obj, linkFieldOfs))
}
