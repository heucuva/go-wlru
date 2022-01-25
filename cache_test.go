package wlru_test

import (
	"math/rand"
	"testing"
	"context"

	"github.com/heucuva/go-wlru"
)

func TestPruning(t *testing.T) {
	c := &wlru.Cache{Capacity: 2}

	input := []wlru.Entry{
		{"a", "hello", false, false},
		{"b", "world", false, false},
		{"c", "matey", false, false},
		{"d", "ahoy", false, false},
	}

	expected := []wlru.Entry{
		input[3],
		input[2],
	}

	for _, entry := range input {
		if err := c.Set(entry.Key, entry.Value, entry.IsPermanent); err != nil {
			t.Fatal(err)
		}
	}

	compareSnapshot(t, c.SnapshotList(), expected)
}

func TestPermanency(t *testing.T) {
	c := &wlru.Cache{Capacity: 2}

	input := []wlru.Entry{
		{"a", "hello", true, false},
		{"b", "world", true, false},
		{"c", "matey", false, false},
		{"d", "ahoy", false, false},
	}

	expected := []wlru.Entry{
		input[1],
		input[0],
	}

	for _, entry := range input {
		if err := c.Set(entry.Key, entry.Value, entry.IsPermanent); err != nil {
			t.Fatal(err)
		}
	}

	compareSnapshot(t, c.SnapshotList(), expected)
}

func TestExpiration(t *testing.T) {
	c := &wlru.Cache{}

	input := []wlru.Entry{
		{"a", "hello", false, false},
		{"b", "world", false, true},
		{"c", "matey", false, false},
		{"d", "ahoy", false, true},
	}

	expectedBeforeRemove := []wlru.Entry{
		input[3],
		input[2],
		input[1],
		input[0],
	}

	expectedAfterPrune := []wlru.Entry{
		input[2],
		input[0],
	}

	expiredCtx, cancel := context.WithCancel(context.Background())

	for _, entry := range input {
		ctx := context.Background()
		if entry.IsExpired {
			ctx = expiredCtx
		}
		if err := c.SetWithContext(ctx, entry.Key, entry.Value, entry.IsPermanent); err != nil {
			t.Fatal(err)
		}
	}

	cancel()

	compareSnapshot(t, c.SnapshotList(), expectedBeforeRemove)

	c.RemoveExpired()

	compareSnapshot(t, c.SnapshotList(), expectedAfterPrune)
}

func TestUnbounded(t *testing.T) {
	c := &wlru.Cache{}

	input := []wlru.Entry{
		{"a", "hello", false, false},
		{"b", "world", false, false},
		{"c", "matey", false, false},
		{"d", "ahoy", false, false},
	}

	expected := []wlru.Entry{
		input[3],
		input[2],
		input[1],
		input[0],
	}

	for _, entry := range input {
		if err := c.Set(entry.Key, entry.Value, entry.IsPermanent); err != nil {
			t.Fatal(err)
		}
	}

	compareSnapshot(t, c.SnapshotList(), expected)
}

func compareSnapshot(t *testing.T, snapshot, expected []wlru.Entry) {
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
