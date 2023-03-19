package storage

import (
	"context"
	"sync"
	"sync/atomic"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{maxWorkersCount: 3}
}

func startWorker(ctx context.Context, in <-chan Dir, errch chan<- error, done chan struct{},
	wg *sync.WaitGroup, res *Result) {

	defer wg.Done()
	for input := range in {
		_, files, err := input.Ls(ctx)
		if err != nil {
			errch <- err
			return
		}
		atomic.AddInt64(&res.Count, int64(len(files)))
		for _, file := range files {
			size, err := file.Stat(ctx)
			if err != nil {
				errch <- err
				return
			}
			atomic.AddInt64(&res.Size, size)
		}
		done <- struct{}{}
	}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	res := Result{Size: 0, Count: 0}
	dirs, files, err := d.Ls(ctx)
	if err != nil {
		return Result{}, err
	}
	res.Count += int64(len(files))
	for _, file := range files {
		size, err := file.Stat(ctx)
		if err != nil {
			return Result{}, err
		}
		res.Size += size
	}
	for _, dir := range dirs {
		p, err := a.Size(ctx, dir)
		if err != nil {
			return Result{}, err
		}
		res.Size += p.Size
		res.Count += p.Count
	}
	return res, nil
}
