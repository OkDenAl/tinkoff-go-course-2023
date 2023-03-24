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
	var (
		inputs  = make([]chan any, 0)
		outputs = make([]Out, 0)
		res     = make(chan any)
		wait    = make(chan struct{})
	)

	j := 0
	for a := range in {
		inputs = append(inputs, make(chan any))
		outputs = append(outputs, make(chan any))

		go func(input In, j int) {
			for i := 0; i < len(stages); i++ {
				outputs[j] = stages[i](input)
				input = outputs[j]
				if i == len(stages)-1 {
					wait <- struct{}{}
				}
			}
		}(inputs[j], j)

		select {
		case <-ctx.Done():
			close(res)
			return res
		default:
			inputs[j] <- a
			close(inputs[j])
			j++
		}
		<-wait
	}

	go func() {
		for _, ch := range outputs {
			res <- <-ch
		}
		close(res)
	}()
	return res
}
