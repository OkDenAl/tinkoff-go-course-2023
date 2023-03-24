package executor

import (
	"context"
	"sync"
)

type any = interface {
}

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	var (
		wg      = sync.WaitGroup{}
		inputs  = make([]chan any, 0)
		outputs = make([]Out, 0)
		res     = make(chan any)
		wait    = make(chan struct{})
		//mu      = sync.Mutex{}
	)
	j := 0

	for a := range in {
		inputs = append(inputs, make(chan any))
		outputs = append(outputs, make(chan any))

		wg.Add(1)

		go func(input In, j int) {
			defer wg.Done()
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

	wg.Wait()
	go func() {
		for _, ch := range outputs {
			res <- <-ch
		}
		close(res)
	}()
	return res
}
