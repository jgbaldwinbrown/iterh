package iterh

import (
	"io"
	"bufio"
	"iter"
)

func LineIter(r io.Reader) iter.Seq2[string, error] {
	return func(y func(string, error) bool) {
		s := bufio.NewScanner(r)
		s.Buffer([]byte{}, 1e12)
		for s.Scan() {
			if ok := y(s.Text(), s.Err()); !ok {
				return
			}
		}
	}
}
