package iterh

import (
	"iter"
)

func Filter[T any](it iter.Seq[T], filt func(T) bool) iter.Seq[T] {
	return func(y func(T) bool) {
		for val := range it {
			if filt(val) {
				if !y(val) {
					return
				}
			}
		}
	}
}

func Transform[T, U any](it iter.Seq[T], trans func(T) U) iter.Seq[U] {
	return func(y func(U) bool) {
		for val := range it {
			if !y(trans(val)) {
				return
			}
		}
	}
}

func Reduce[T any](it iter.Seq[T], red func(T, T) T) T {
	var sum T
	i := 0
	for val := range it {
		if i == 0 {
			sum = val
		} else {
			sum = red(sum, val)
		}
		i++
	}
	return sum
}

func Max[T any](it iter.Seq[T], cmp func(T, T) int) T {
	var hi T
	started := false
	for val := range it {
		if !started || cmp(hi, val) < 0 {
			hi = val
			started = true
		}
	}
	return hi
}
