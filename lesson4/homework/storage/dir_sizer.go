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
	res := Result{}
	workerInput := make(chan Dir)
	errorCh := make(chan error)
	done := make(chan struct{})
	wg := sync.WaitGroup{}

	for i := 0; i < a.maxWorkersCount; i++ {
		wg.Add(1)
		go startWorker(ctx, workerInput, errorCh, done, &wg, &res)
	}

	que := make([]Dir, 0) // queue ds
	que = append(que, d)
	workerInput <- d
	for {
		if len(done) == 0 && len(que) == 0 {
			close(workerInput)
			break
		}
		select {
		case err := <-errorCh:
			return Result{}, err
		case <-done:
			curDir := que[0]
			que = que[1:]
			dirs, _, _ := curDir.Ls(ctx)
			for _, dir := range dirs {
				que = append(que, dir)
				workerInput <- dir
			}
		}
	}
	wg.Wait()
	return res, nil
}
