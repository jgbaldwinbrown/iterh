package iterh

import (
	"iter"
	"slices"
)

func Repeat[T any](vals ...T) iter.Seq[T] {
	return func(y func(T) bool) {
		for {
			for _, val := range vals {
				if !y(val) {
					return
				}
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

func RepeatN[T any](n int, vals ...T) iter.Seq[T] {
	return Head(Repeat(vals...), n * len(vals))
}

func Zero[T any]() T {
	var t T
	return t
}

func RepeatSlice[T any](n int, vals ...T) []T {
	out := make([]T, 0, n * len(vals))
	return slices.AppendSeq(out, RepeatN(n, vals...))
}
