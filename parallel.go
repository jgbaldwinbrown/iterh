package iterh

import (
	"golang.org/x/sync/errgroup"
	"iter"
)

func ParallelRun(fs iter.Seq[func() error], threads int) error {
	var g errgroup.Group
	if threads > 0 {
		g.SetLimit(threads)
	}
	for f := range fs {
		g.Go(f)
	}
	return g.Wait()
}
