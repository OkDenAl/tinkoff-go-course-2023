package storage

import (
	"context"
	"log"
	"runtime"
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
	return &sizer{}
}

func startWorker(in <-chan Dir, errch chan<- error, wg *sync.WaitGroup, ctx context.Context, res *Result) {
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
		runtime.Gosched()
	}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	res := &Result{}
	worketInput := make(chan Dir)
	errorCh := make(chan error, 1)
	a.maxWorkersCount = 3
	var wg sync.WaitGroup
	for i := 0; i < a.maxWorkersCount; i++ {
		wg.Add(1)
		go startWorker(worketInput, errorCh, &wg, ctx, res)
	}
	que := make(chan Dir, 100)
	que <- d
	worketInput <- d
	for len(que) != 0 {
		curDir := <-que
		dirs, _, _ := curDir.Ls(ctx)
		for _, dir := range dirs {
			que <- dir
			worketInput <- dir
		}
	}
	close(worketInput)
	wg.Wait()
	select {
	case err := <-errorCh:
		log.Println(err)
		return Result{}, err
	default:
		return *res, nil
	}
}
