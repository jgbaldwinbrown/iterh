package iterh

import (
	"iter"
	"slices"
)

func BreakOnError[T any](it iter.Seq2[T, error], ep *error) iter.Seq[T] {
	return func(yield func(T) bool) {
		for val, err := range it {
			if err != nil {
				*ep = err
				break
			}
			if !yield(val) {
				break
			}
		}
	}
}

func BreakWithError[T any](it iter.Seq2[T, error]) (iter.Seq[T], *error) {
	ep := new(error)
	return BreakOnError(it, ep), ep
}

func BreakOnErrorMulti[T any](its iter.Seq[iter.Seq2[T, error]], ep *error) iter.Seq[iter.Seq[T]] {
	return func(y func(iter.Seq[T]) bool) {
		for it := range its {
			it2 := BreakOnError(it, ep)
			if !y(it2) {
				return
			}
			if *ep != nil {
				return
			}
		}
	}
}

func CollectWithError[T any](it iter.Seq2[T, error]) ([]T, error) {
	var e error
	it1 := BreakOnError(it, &e)
	s := slices.Collect(it1)
	if e != nil {
		return nil, e
	}
	return s, nil
}

func AddDummy[T, D any](it iter.Seq[T], dummy D) iter.Seq2[T, D] {
	return Zip(it, RepeatForever(dummy))
}

func AddNilError[T any](it iter.Seq[T]) iter.Seq2[T, error] {
	return Zip(it, RepeatForever(Zero[error]()))
}

func Swap[T, U any](it iter.Seq2[T, U]) iter.Seq2[U, T] {
	return func(yield func(U, T) bool) {
		for t, u := range it {
			if !yield(u, t) {
				break
			}
		}
	}
}

func Enumerate[T any](it iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for t := range it {
			if !yield(i, t) {
				break
			}
			i++
		}
	}
}

func Zip[T, U any](ti iter.Seq[T], ui iter.Seq[U]) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		up, ucancel := iter.Pull(ui)
		defer ucancel()
		for t := range ti {
			u, ok := up()
			if !ok {
				return
			}
			if !yield(t, u) {
				return
			}
		}
	}
}

func ChannelIter[T any](c <-chan T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for val := range c {
			if !yield(val) {
				return
			}
		}
	}
}

func IterChannel[T any](it iter.Seq[T], chanLen int) (c <-chan T, cancel func()) {
	c1 := make(chan T, chanLen)
	cancelc := make(chan struct{})
	go func() {
		defer close(c1)
		for t := range it {
			select {
			case <-cancelc:
				return
			case c1 <- t:
			}
		}
	}()
	return c1, func() { close(cancelc) }
}

type Number interface {
	int | int8 | int16 | int32 | int64 |
	uint | uint8 | uint16 | uint32 | uint64 |
	float64 | float32
}

func Range[N Number](start, end, step N) iter.Seq[N] {
	var z N
	if step >= z {
		return func(y func(N) bool) {
			for i := start; i < end; i += step {
				if !y(i) {
					return
				}
			}
		}
	}
	return func(y func(N) bool) {
		for i := start; i > end; i += step {
			if !y(i) {
				return
			}
		}
	}
}

func SlicePtrs[S ~[]T, T any](s S) iter.Seq[*T] {
	return func(y func(*T) bool) {
		for i, _ := range s {
			if !y(&s[i]) {
				return
			}
		}
	}
}

func First[T, U any](it iter.Seq2[T, U]) iter.Seq[T] {
	return func(y func(T) bool) {
		for t, _ := range it {
			if !y(t) {
				return
			}
		}
	}
}

func Second[T, U any](it iter.Seq2[T, U]) iter.Seq[U] {
	return func(y func(U) bool) {
		for _, u := range it {
			if !y(u) {
				return
			}
		}
	}
}
