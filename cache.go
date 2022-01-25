package wlru

import (
	"context"

	"github.com/pkg/errors"
)

// Cache is a fast, lock-free Least Recently Used cache with partitioning and weighting.
type Cache struct {
	// Capacity if > 0 will limit the number of entries the map Cache will contain
	// if == 0, then no limit
	Capacity int

	// cache is a doubly-linked list which tracks order of the all entries.
	// head is most recent, tail is least.
	cache locklessList
	// items holds items that are O(1) for key/value lookup and O(1) for removal
	items safeMap
}

// Size returns the number of items in the map, which should reflect all items in the DLL.
func (c *Cache) Size() int {
	return c.items.size()
}

// Get returns the value and a boolean true value, if found.
// If not, then nil and false.
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	if key == nil {
		return nil, false
	}

	entry, found := c.items.get(key)
	if !found {
		return nil, false
	}

	// remove it from the spot it's in the cache, wherever it might be.
	c.cache.remove(entry)

	if entry.isExpired() {
		c.items.remove(key)
		return nil, false
	}

	// push it onto the head (MRU)
	c.cache.pushHead(entry)

	return entry.value.Load(), true
}

// Set adds an item into the cache or replaces the value if an item with
// the same key already exists there.
func (c *Cache) Set(key, value interface{}, permanent bool) error {
	n, err := newNode(context.Background(), key, value, permanent)
	if err != nil {
		return err
	}
	return c.set(n)
}

// SetWithContext adds an item into the cache or replaces the value if an
// item with the same key already exists there, but this time with a context
func (c *Cache) SetWithContext(ctx context.Context, key, value interface{}, permanent bool) error {
	n, err := newNode(ctx, key, value, permanent)
	if err != nil {
		return err
	}
	return c.set(n)
}

func (c *Cache) set(n *node) error {
	if n == nil {
		panic("node is nil")
	}

	key := n.key

	// does key already exist?
	entry, found := c.items.get(key)
	if found {
		// update the entry
		entry.update(n)

		// remove it from the spot it's in the cache, wherever it might be.
		c.cache.remove(entry)
		c.cache.pushHead(entry)
		return nil
	}

	entry = n

	// add a new node to the list
	c.cache.pushHead(entry)

	// add the node to the map
	c.items.set(entry)

	// check to see if the cache capacity has been exceeded
	// and seek to the oldest item in the cache that we can prune
	entry = c.cache.tail()
	for c.Capacity > 0 && c.Size() > c.Capacity && entry != nil {
		// take note of our next entry
		prev := entry.prev()
		// are we looking at a permanent entry?
		if entry.isPermanent() && !entry.isExpired() {
			// we've found a permanent entry (that's not expired), skip it.
			// seek to the next in line
			entry = prev
			continue
		}

		// we have a mutable (and/or expired) item that can be pruned
		// yank it out of the map
		c.items.remove(entry.key)
		// remove it from the list
		c.cache.remove(entry)
		// seek to the next in line
		entry = prev
	}

	return nil
}

// Remove removes an item from the cache by key
func (c *Cache) Remove(key interface{}) (bool, error) {
	if key == nil {
		return false, errors.Wrap(ErrNilArgument, "key")
	}

	if entry, found := c.items.remove(key); found {
		// remove it from the spot it's in the cache, wherever it might be.
		c.cache.remove(entry)
		return true, nil
	}

	return false, nil
}

// SnapshotList returns a snapshot of the list entries along with
// pertenent cache information. Entries are ordered from MRU..LRU (0..n)
func (c *Cache) SnapshotList() []Entry {
	var entries []Entry
	for entry := c.cache.head(); entry != nil; entry = entry.next() {
		entries = append(entries, entry.toEntry())
	}
	return entries
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	entry := c.cache.head()
	for entry != nil {
		prev := entry.next()
		c.cache.remove(entry)
		entry = prev
	}
	c.items.clear()
}

// RemoveExpired removes the expired entries from the cache, then
// returns the number of removed items.
func (c *Cache) RemoveExpired() int {
	var numExpired int
	entry := c.cache.head()
	for entry != nil {
		next := entry.next()
		if entry.isExpired() {
			c.cache.remove(entry)
			c.items.remove(entry.key)
			numExpired++
		}
		entry = next
	}

	return numExpired
}
