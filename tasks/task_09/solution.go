package main

import (
	"context"
	"errors"
	"sync"
)

type Result[R any] struct {
	Value R
	Err   error
}

var ErrInvalidWorkers = errors.New("workers must be positive")

func ParallelMap[T any, R any](ctx context.Context, workers int, in []T, fn func(context.Context, T) (R, error)) ([]Result[R], error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkers
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if len(in) == 0 {
		return make([]Result[R], 0), nil
	}

	result := make([]Result[R], len(in))

	wg := sync.WaitGroup{}

	count := min(workers, len(in))

	ch := make(chan int)

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for x := range ch {
				v, err := fn(ctx, in[x])
				if err != nil {
					var zero R
					result[x] = Result[R]{Value: zero, Err: err}
					continue
				}
				result[x] = Result[R]{Value: v, Err: nil}
			}
		}()
	}

	cancelled := false
	for i := range in {
		if ctx.Err() != nil {
			cancelled = true
			break
		}
		select {
		case <-ctx.Done():
			cancelled = true
		case ch <- i:
		}
		if cancelled {
			break
		}
	}
	close(ch)
	wg.Wait()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
