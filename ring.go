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
		if !y(n) {
			return
		}
		for n = n.Next(); n != r; n = n.Next() {
			if !y(n) {
				return
			}
		}
	}
}

func RingPtrs(r *ring.Ring) iter.Seq[*any] {
	return Transform(RingIter(r), func(r *ring.Ring) *any {
		return &r.Value
	})
}

func RingValues(r *ring.Ring) iter.Seq[any] {
	return Elems(RingPtrs(r))
}
