package iterh

import (
	"iter"
)

func Elems[T any](it iter.Seq[*T]) iter.Seq[T] {
	return func(y func(T) bool) {
		for p := range it {
			if ok := y(*p); !ok {
				return
			}
		}
	}
}
