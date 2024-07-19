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

func ListElementValuePtrs(it iter.Seq[*list.Element]) iter.Seq[*any] {
	return Transform(it, func(e *list.Element) *any {
		return &e.Value
	})
}
