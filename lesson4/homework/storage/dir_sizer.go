package storage

import (
	"context"
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

	// TODO: add other fields as you wish
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
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
