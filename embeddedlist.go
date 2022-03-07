package wlru

// This is a double-linked list container - it allows for linear iteration over
// its contents.
// This cointainer does not take ownership of its contents, so the application
// must remove items manually.

type embeddedList struct {
	head  *node
	tail  *node
	count int
}

func (c *embeddedList) getLink(obj *node) *embeddedListLink {
	return getListLink(obj)
}

func (c *embeddedList) First() *node {
	return c.head
}

func (c *embeddedList) Last() *node {
	return c.tail
}

func (c *embeddedList) Next(cur *node) *node {
	return c.getLink(cur).next
}

func (c *embeddedList) Prev(cur *node) *node {
	return c.getLink(cur).prev
}

func (c *embeddedList) Position(index int) *node {
	cur := c.head
	for cur != nil && index > 0 {
		cur = c.Next(cur)
		index--
	}
	return cur
}

func (c *embeddedList) Count() int {
	return c.count
}

func (c *embeddedList) Remove(obj *node) *node {
	if c.getLink(obj).remove(&c.head, &c.tail) {
		c.count--
	}
	return obj
}

func (c *embeddedList) RemoveFirst() *node {
	if c.head == nil {
		return nil
	}
	return c.Remove(c.head)
}

func (c *embeddedList) RemoveLast() *node {
	if c.tail == nil {
		return nil
	}
	return c.Remove(c.tail)
}

func (c *embeddedList) RemoveAll() {
	for cur := c.tail; cur != nil; cur = c.tail {
		c.Remove(cur)
	}
}

func (c *embeddedList) InsertFirst(cur *node) *node {
	c.getLink(cur).next = c.head
	if c.head != nil {
		c.getLink(c.head).prev = cur
		c.head = cur
	} else {
		c.head = cur
		c.tail = cur
	}
	c.count++
	return cur
}

func (c *embeddedList) InsertLast(cur *node) *node {
	c.getLink(cur).prev = c.tail
	if c.tail != nil {
		c.getLink(c.tail).next = cur
		c.tail = cur
	} else {
		c.head = cur
		c.tail = cur
	}
	c.count++
	return cur
}

func (c *embeddedList) InsertAfter(prev, cur *node) *node {
	if prev == nil {
		return c.InsertFirst(cur)
	}
	curU := c.getLink(cur)
	prevU := c.getLink(prev)
	curU.prev = prev
	curU.next = prevU.next
	prevU.next = cur

	if curU.next != nil {
		c.getLink(curU.next).prev = cur
	} else {
		c.tail = cur
	}

	c.count++
	return cur
}

func (c *embeddedList) InsertBefore(after, cur *node) *node {
	if after == nil {
		return c.InsertLast(cur)
	}
	curU := c.getLink(cur)
	afterU := c.getLink(after)
	curU.next = after
	curU.prev = afterU.prev
	afterU.prev = cur

	if curU.prev != nil {
		c.getLink(curU.prev).next = cur
	} else {
		c.head = cur
	}

	c.count++
	return cur
}

func (c *embeddedList) MoveFirst(cur *node) {
	c.Remove(cur)
	c.InsertFirst(cur)
}

func (c *embeddedList) MoveLast(cur *node) {
	c.Remove(cur)
	c.InsertLast(cur)
}

func (c *embeddedList) MoveAfter(dest, cur *node) {
	c.Remove(cur)
	c.InsertAfter(dest, cur)
}

func (c *embeddedList) MoveBefore(dest, cur *node) {
	c.Remove(cur)
	c.InsertBefore(dest, cur)
}

func (c *embeddedList) IsEmpty() bool {
	return c.count == 0
}

func (c *embeddedList) IsContained(cur *node) bool {
	return c.getLink(cur).isContained(c.head)
}
