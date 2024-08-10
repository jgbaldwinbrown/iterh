package iterh

import (
	"iter"
)

func ReverseSliceIter[S ~[]T, T any](s S) iter.Seq[T] {
	return func(y func(T) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !y(s[i]) {
				return
			}
		}
	}
}

func Reverse[T any](it iter.Seq[T]) iter.Seq[T] {
	s := Collect(it)
	return ReverseSliceIter(s)
}

type LenSwapper interface {
	Len() int
	Swap(i, j int)
}

func ReverseInPlace[S LenSwapper](s S) {
	l := s.Len()
	mid := l / 2
	for i := 0; i < mid; i++ {
		j := l - i
		s.Swap(i, j)
	}
}

type SliceLenSwapper[T any] []T

func (s SliceLenSwapper[T]) Len() int {
	return len(s)
}

func (s SliceLenSwapper[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func ReverseSlice[S ~[]T, T any](s S) {
	ReverseInPlace(SliceLenSwapper[T](s))
}
