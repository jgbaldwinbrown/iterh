package iterh

import (
	"iter"
)

type View[T any] interface {
	Len() int
	At(i int) T
}

func ViewIter[V View[T], T any](v V) iter.Seq[T] {
	return func(y func(T) bool) {
		l := v.Len()
		for i := 0; i < l; i++ {
			if !y(v.At(i)) {
				return
			}
		}
	}
}

func PulledIter[T any](next func() (T, bool)) iter.Seq[T] {
	return func(y func(T) bool) {
		for val, ok := next(); ok; val, ok = next() {
			if !y(val) {
				return
			}
		}
	}
}

func PulledIter2[T, U any](next func() (T, U, bool)) iter.Seq2[T, U] {
	return func(y func(T, U) bool) {
		for t, u, ok := next(); ok; t, u, ok = next() {
			if !y(t, u) {
				return
			}
		}
	}
}
