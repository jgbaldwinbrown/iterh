package iterh

import (
	"iter"
)

func Collect[T any](it iter.Seq[T]) []T {
	var out []T
	for val := range it {
		out = append(out, val)
	}
	return out
}

func BreakOnError[T any](it iter.Seq2[T, error], ep *error) iter.Seq[T] {
	return func(yield func(T) bool) {
		for val, err := range it {
			if err != nil {
				*ep = err
				break
			}
			if ok := yield(val); !ok {
				break
			}
		}
	}
}

func CollectWithError[T any](it iter.Seq2[T, error]) ([]T, error) {
	var e error
	it1 := BreakOnError(it, &e)
	s := Collect(it1)
	if e != nil {
		return nil, e
	}
	return s, nil
}

func AddDummy[T, D any](it iter.Seq[T]) iter.Seq2[T, D] {
	return func(yield func(T, D) bool) {
		for val := range it {
			var d D
			if ok := yield(val, d); !ok {
				break
			}
		}
	}
}

func Swap[T, U any](it iter.Seq2[T, U]) iter.Seq2[U, T] {
	return func(yield func(U, T) bool) {
		for t, u := range it {
			if ok := yield(u, t); !ok {
				break
			}
		}
	}
}

func Enumerate[T any](it iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for t := range it {
			if ok := yield(i, t); !ok {
				break
			}
			i++
		}
	}
}

func CollectMap[K comparable, V any](it iter.Seq2[K, V]) map[K]V {
	m := map[K]V{}
	for k, v := range it {
		m[k] = v
	}
	return m
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
			if ok := yield(t, u); !ok {
				return
			}
		}
	}
}

func ChannelIter[T any](c <-chan T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for val := range c {
			if ok := yield(val); !ok {
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
				if ok := y(i); !ok {
					return
				}
			}
		}
	}
	return func(y func(N) bool) {
		for i := start; i > end; i += step {
			if ok := y(i); !ok {
				return
			}
		}
	}
}

func SlicePointerIter[S ~[]T, T any](s []T) iter.Seq[*T] {
	return func(y func(*T) bool) {
		for i, _ := range s {
			if ok := y(&s[i]); !ok {
				return
			}
		}
	}
}