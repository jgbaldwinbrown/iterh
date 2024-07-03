package iterh

import (
	"io"
	"os"
	"iter"
	"github.com/jgbaldwinbrown/csvh"
	"encoding/json"
	"encoding/gob"
)

func PathIter[T any](path string, f func(io.Reader) iter.Seq2[T, error]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		r, err := os.Open(path)
		if err != nil {
			var t T
			if ok := yield(t, err); !ok {
				return
			}
		}
		defer r.Close()
		it := f(r)
		it(yield)
	}
}

func MaybeGzPathIter[T any](path string, f func(io.Reader) iter.Seq2[T, error]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		r, err := csvh.OpenMaybeGz(path)
		if err != nil {
			var t T
			if ok := yield(t, err); !ok {
				return
			}
		}
		defer r.Close()
		it := f(r)
		it(yield)
	}
}

func GzPathIter[T any](path string, f func(io.Reader) iter.Seq2[T, error]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		r, err := csvh.GzOpen(path)
		if err != nil {
			var t T
			if ok := yield(t, err); !ok {
				return
			}
		}
		defer r.Close()
		it := f(r)
		it(yield)
	}
}

func CsvIter(r io.Reader) iter.Seq2[[]string, error] {
	cr := csvh.CsvIn(r)
	return func(yield func([]string, error) bool) {
		for l, e := cr.Read(); e != io.EOF; l, e = cr.Read() {
			if ok := yield(l, e); !ok {
				return
			}
		}
	}
}

func JsonIter[T any](r io.Reader) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		dec := json.NewDecoder(r)
		var t T
		for e := dec.Decode(&t); e != io.EOF; e = dec.Decode(&t) {
			if ok := yield(t, e); !ok {
				return
			}
		}
	}
}

func WriteJson[T any](w io.Writer, it iter.Seq2[T, error]) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	for t, err := range it {
		if err != nil {
			return err
		}
		err = enc.Encode(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func GobIter[T any](r io.Reader) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		dec := gob.NewDecoder(r)
		var t T
		for e := dec.Decode(&t); e != io.EOF; e = dec.Decode(&t) {
			if ok := yield(t, e); !ok {
				return
			}
		}
	}
}

func WriteGob[T any](w io.Writer, it iter.Seq2[T, error]) error {
	enc := gob.NewEncoder(w)
	for t, err := range it {
		if err != nil {
			return err
		}
		err = enc.Encode(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func WritePath[T any](path string, f func(io.Writer) error) (err error) {
	w, e := os.Open(path)
	if e != nil {
		return e
	}
	defer func() { csvh.DeferE(&err, w.Close()) }()
	return f(w)
}

func WriteGzPath[T any](path string, f func(io.Writer) error) (err error) {
	w, e := csvh.GzCreate(path)
	if e != nil {
		return e
	}
	defer func() { csvh.DeferE(&err, w.Close()) }()
	return f(w)
}

func WriteMaybeGzPath[T any](path string, f func(io.Writer) error) (err error) {
	w, e := csvh.CreateMaybeGz(path)
	if e != nil {
		return e
	}
	defer func() { csvh.DeferE(&err, w.Close()) }()
	return f(w)
}
