package iterh

import (
	"iter"
	"container/list"
)

func ListElements(l *list.List) iter.Seq[*list.Element] {
	return func(y func(*list.Element) bool) {
		for e := l.Front(); e != nil; e = e.Next() {
			if !y(e) {
				return
			}
		}
	}
}

func ListPtrs(l *list.List) iter.Seq[*any] {
	return Transform(ListElements(l), func(e *list.Element) *any {
		return &e.Value
	})
}

func ListValues(l *list.List) iter.Seq[any] {
	return Elems(ListPtrs(l))
}
