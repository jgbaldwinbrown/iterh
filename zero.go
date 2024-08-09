package iterh

import (
	"iter"
)

func Repeat[T any](val T) iter.Seq[T] {
	return func(y func(T) bool) {
		for {
			if !y(val) {
				return
			}
		}
	}
}

func Head[T any](it iter.Seq[T], n int) iter.Seq[T] {
	return func(y func(T) bool) {
		if n < 1 {
			return
		}

		i := 0
		for val := range it {
			if !y(val) {
				return
			}
			i++
			if i >= n {
				return
			}
		}
	}
}

func CutHead[T any](it iter.Seq[T], n int) iter.Seq[T] {
	return func(y func(T) bool) {
		i := 0
		for val := range it {
			if i >= n {
				if !y(val) {
					return
				}
			}
			i++
		}
	}
}

func RepeatN[T any](val T, n int) iter.Seq[T] {
	return Head(Repeat(val), n)
}

func Zero[T any]() T {
	var t T
	return t
}

func RepeatSlice[T any](val T, n int) []T {
	out := make([]T, 0, n)
	return Append(out, RepeatN(val, n))
}
