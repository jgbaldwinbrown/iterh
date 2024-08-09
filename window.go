package iterh

import (
	"iter"
	"github.com/gammazero/deque"
)

type WinView[T any] struct {
	d *deque.Deque[T]
}

func (w WinView[T]) Len() int {
	return w.d.Len()
}

func (w WinView[T]) At(i int) T {
	return w.d.At(i)
}

func Window[T any](it iter.Seq[T], winsize int, winstep int) iter.Seq[WinView[T]] {
	return func(y func(WinView[T]) bool) {
		d := deque.New[T](winsize + winstep, winsize)
		i := -1
		started := false
		for val := range it {
			if !started && d.Len() < winsize {
				d.PushBack(val)
				continue
			}

			started = true
			i++
			if i % winstep == 0 {
				for d.Len() > winsize {
					d.PopFront()
				}
				if !y(WinView[T]{d}) {
					return
				}
			}
		}
		if i % winstep != 0 {
			for d.Len() > winsize {
				d.PopFront()
			}
			if !y(WinView[T]{d}) {
				return
			}
		}
	}
}
