package iterh

import (
	"iter"
)

func Sum[T Number](it iter.Seq[T]) T {
	var sum T
	for x := range it {
		sum += x
	}
	return sum
}
