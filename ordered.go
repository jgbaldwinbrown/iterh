package iterh

import (
	"iter"
	"sync"
	"golang.org/x/sync/errgroup"
)

type Ordered[T any] struct {
	buf map[int]T
	pos int
	mu sync.Mutex
	cond *sync.Cond
	closed bool
}

func NewOrdered[T any]() *Ordered[T] {
	o := new(Ordered[T])
	o.buf = map[int]T{}
	o.cond = sync.NewCond(&o.mu)
	return o
}

func (o *Ordered[T]) tryGet() (val T, open, got bool) {
	if val, ok := o.buf[o.pos]; ok {
		delete(o.buf, o.pos)
		o.pos++
		return val, true, true
	}
	var t T
	return t, !o.closed, false
}

func (o *Ordered[T]) TryGet() (val T, open, got bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.tryGet()
}

func (o *Ordered[T]) BlockingGet() (val T, got bool) {
	o.mu.Lock()
	defer o.mu.Unlock()

	for {
		val, open, got := o.tryGet()
		if got {
			return val, true
		}
		if !open {
			var t T
			return t, false
		}
		o.cond.Wait()
	}
}

func (o *Ordered[T]) All() iter.Seq[T] {
	return func(y func(T) bool) {
		for val, ok := o.BlockingGet(); ok; val, ok = o.BlockingGet() {
			if !y(val) {
				return
			}
		}
	}
}

func (o *Ordered[T]) Put(i int, val T) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.buf[i] = val
	o.cond.Broadcast()
}

func (o *Ordered[T]) Close() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.closed = true
	o.cond.Broadcast()
}

func OrderedParallel[T, U any](f func(T) U, in iter.Seq[T], threads int) *Ordered[U] {
	return OrderedParallelIndexed(f, Enumerate(in), threads)
}

func OrderedParallelIndexed[T, U any](f func(T) U, in iter.Seq2[int, T], threads int) *Ordered[U] {
	o := NewOrdered[U]()

	go func() {
		var g errgroup.Group
		if threads > 0 {
			g.SetLimit(threads)
		}
		for i, val := range in {
			g.Go(func() error {
				o.Put(i, f(val))
				return nil
			})
		}
		g.Wait()
		o.Close()
	}()

	return o
}
