package wlru

import (
	"context"
	"sync/atomic"

	"github.com/pkg/errors"
)

type node struct {
	key       interface{}  // immutable data
	value     atomic.Value // interface{}
	ctx       atomic.Value // context.Context
	permanent atomic.Value // bool

	link embeddedListLink
}

func newNode(ctx context.Context, key, value interface{}, permanent bool) (*node, error) {
	if key == nil {
		return nil, errors.Wrap(ErrNilArgument, "key")
	}

	if value == nil {
		return nil, errors.Wrap(ErrNilArgument, "value")
	}

	n := &node{
		key: key,
	}

	n.value.Store(value)
	n.ctx.Store(ctx)
	n.permanent.Store(permanent)

	return n, nil
}

func (n *node) isExpired() bool {
	ctx, _ := n.ctx.Load().(context.Context)
	if ctx == nil {
		return false
	}

	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func (n *node) isPermanent() bool {
	v, _ := n.permanent.Load().(bool)
	return v
}

func (n *node) update(other *node) {
	n.value.Store(other.value.Load())
	n.ctx.Store(other.ctx.Load())
	n.permanent.Store(other.permanent.Load())
}

func (n *node) toEntry() Entry {
	e := Entry{
		Key:         n.key,
		Value:       n.value.Load(),
		IsPermanent: n.isPermanent(),
		IsExpired:   n.isExpired(),
	}
	return e
}
