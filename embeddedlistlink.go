package wlru

import (
	"unsafe"
)

// embeddedListLink is a link to the list container
type embeddedListLink struct {
	prev *node
	next *node
}

func (l *embeddedListLink) remove(head **node, tail **node) bool {
	if !l.isContained(*head) {
		return false
	}
	if l.prev == nil {
		*head = l.next
	} else {
		getListLink(l.prev).next = l.next
	}
	if l.next == nil {
		*tail = l.prev
	} else {
		getListLink(l.next).prev = l.prev
	}

	l.next = nil
	l.prev = nil
	return true
}

func (l *embeddedListLink) isContained(head *node) bool {
	return l.prev != nil || head == l.getItem(unsafe.Offsetof(head.link))
}

func (l *embeddedListLink) getItem(linkFieldOfs uintptr) *node {
	u := unsafe.Add(unsafe.Pointer(l), (^linkFieldOfs)+1)
	m := (*node)(u)
	return m
}

func getListLink(obj *node) *embeddedListLink {
	return &obj.link
}
