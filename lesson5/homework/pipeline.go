package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	channels := make([]In, len(stages)+1)
	res := make(chan any)

	for i := range channels {
		channels[i] = make(In)
	}
	channels[0] = in

	for i := 0; i < len(stages); i++ {
		channels[i+1] = stages[i](channels[i])
	}

	go func() {
	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			case data, ok := <-channels[len(stages)]:
				if !ok {
					break Loop
				}
				res <- data
			}
		}
		close(res)
	}()
	return res
}
