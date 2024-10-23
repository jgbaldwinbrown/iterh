package iterh

import (
	"iter"
)

func Flatten[T any](itit iter.Seq[iter.Seq[T]]) iter.Seq[T] {
	return func(y func(T) bool) {
		for it := range itit {
			for x := range it {
				if !y(x) {
					return
				}
			}
		}
	}
}

func Concatenate[T any](itit iter.Seq[iter.Seq[T]]) iter.Seq[T] {
	return Flatten(itit)
}
