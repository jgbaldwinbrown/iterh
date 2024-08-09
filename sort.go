package iterh

import (
	"slices"
	"iter"
	"cmp"
)

func SortFunc[T any](it iter.Seq[T], cmpf func(T, T) int) []T {
	s := Collect(it)
	slices.SortFunc(s, cmpf)
	return s
}

func Sort[T cmp.Ordered](it iter.Seq[T]) []T {
	return SortFunc(it, cmp.Compare)
}

func SortedFunc[T any](it iter.Seq[T], cmpf func(T, T) int) iter.Seq[T] {
	s := SortFunc(it, cmpf)
	return SliceIter(s)
}

func Sorted[T cmp.Ordered](it iter.Seq[T]) iter.Seq[T] {
	return SortedFunc(it, cmp.Compare)
}
