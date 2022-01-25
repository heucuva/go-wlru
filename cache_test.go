package wlru_test

import (
	"math/rand"
	"testing"

	"github.com/heucuva/go-wlru"
)

func TestCache(t *testing.T) {
	c := &wlru.Cache{Capacity: 2}

	if err := c.Set("a", "hello", true); err != nil {
		t.Fatal(err)
	}
	if err := c.Set("b", "world", false); err != nil {
		t.Fatal(err)
	}
	if err := c.Set("a", "ahoy", true); err != nil {
		t.Fatal(err)
	}
	if err := c.Set("b", "planet", false); err != nil {
		t.Fatal(err)
	}
	if err := c.Set("c", "earth", false); err != nil {
		t.Fatal(err)
	}
	if err := c.Set("b", "space", false); err != nil {
		t.Fatal(err)
	}
	expected := []wlru.Entry{
		{Key: "b", Value: "space", IsPermanent: false, IsExpired: false},
		//{Key: "c", Value: "earth", IsPermanent: false, IsExpired: false}, // this item will be pruned because capacity is 2
		{Key: "a", Value: "ahoy", IsPermanent: true, IsExpired: false},
	}
	snapshot := c.SnapshotList()
	if len(snapshot) != len(expected) {
		t.Fatalf("unexpected snapshot size %d != expected %d", len(snapshot), len(expected))
	}
	for i, entry := range snapshot {
		ex := expected[i]
		if entry.Key != ex.Key {
			t.Fatalf("[%d] unexpected key %q != expected %q", i, entry.Key, ex.Key)
		}
		if entry.Value != ex.Value {
			t.Fatalf("[%d] unexpected value %q != expected %q", i, entry.Value, ex.Value)
		}
		if entry.IsPermanent != ex.IsPermanent {
			t.Fatalf("[%d] unexpected isPermanent %v != expected %v", i, entry.IsPermanent, ex.IsPermanent)
		}
		if entry.IsExpired != ex.IsExpired {
			t.Fatalf("[%d] unexpected isExpired %v != expected %v", i, entry.IsExpired, ex.IsExpired)
		}
	}
}

func BenchmarkSet(b *testing.B) {
	c := &wlru.Cache{}

	for i := 0; i < 100000; i++ {
		if err := c.Set(i, rand.Int63(), false); err != nil {
			b.Fatal(err)
		}
	}
}
