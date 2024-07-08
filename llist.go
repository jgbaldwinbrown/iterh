package iterh

import (
	"iter"
	"container/list"
)

func ListElements(l *list.List) iter.Seq[*list.Element] {
	return func(y func(*list.Element) bool) {
		for e := l.Front(); e != nil; e = e.Next() {
			if ok := y(e); !ok {
				return
			}
		}
	}
}

func ListElementValues(it iter.Seq[*list.Element]) iter.Seq[any] {
	return func(y func(any) bool) {
		for e := range it {
			if ok := y(e.Value); !ok {
				return
			}
		}
	}
}

func ListValues(l *list.List) iter.Seq[any] {
	return ListElementValues(ListElements(l))
}
