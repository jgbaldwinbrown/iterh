package iterh

import (
	"iter"
	"cmp"
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

func MaxFunc[T any](it iter.Seq[T], cmpf func(T, T) int) T {
	var hi T
	started := false
	for val := range it {
		if !started || cmpf(hi, val) < 0 {
			hi = val
			started = true
		}
	}
	return hi
}

func MinFunc[T any](it iter.Seq[T], cmpf func(T, T) int) T {
	return MaxFunc(it, Negative(cmpf))
}

func Max[T cmp.Ordered](it iter.Seq[T]) T {
	return MaxFunc(it, cmp.Compare)
}

func Negative[T any](cmpf func(T, T) int) func(T, T) int {
	return func(a, b T) int {
		return -cmpf(a, b)
	}
}

func Min[T cmp.Ordered](it iter.Seq[T]) T {
	return MinFunc(it, cmp.Compare)
}

func RankFunc[T any](target T, it iter.Seq[T], cmpf func(T, T) int) (perc float64, nhigher, total int) {
	for val := range it {
		total++
		if cmpf(val, target) > 0 {
			nhigher++
		}
	}
	return float64(nhigher) / float64(total), nhigher, total
}

func Rank[T cmp.Ordered](target T, it iter.Seq[T]) (perc float64, nhigher, total int) {
	return RankFunc(target, it, cmp.Compare)
}

func IndexFunc[T any](it iter.Seq[T], idxf func(T) bool) (i int, val T) {
	for i, val := range Enumerate(it) {
		if idxf(val) {
			return i, val
		}
	}
	var t T
	return -1, t
}

func Index[T comparable](target T, it iter.Seq[T]) (i int) {
	for i, val := range Enumerate(it) {
		if val == target {
			return i
		}
	}
	return -1
}
