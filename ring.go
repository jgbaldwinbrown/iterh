package iterh

import (
	"iter"
	"container/ring"
)

func RingIter(r *ring.Ring) iter.Seq[*ring.Ring] {
	return func(y func(*ring.Ring) bool) {
		if r == nil {
			return
		}
		n := r
		if ok := y(n); !ok {
			return
		}
		for n = n.Next(); n != r; n = n.Next() {
			if ok := y(n); !ok {
				return
			}
		}
	}
}

func RingElementValuesIter(it iter.Seq[*ring.Ring]) iter.Seq[any] {
	return func(y func(any) bool) {
		for r := range it {
			if ok := y(r.Value); !ok {
				return
			}
		}
	}
}

func RingValuesIter(r *ring.Ring) iter.Seq[any] {
	return RingElementValuesIter(RingIter(r))
}
