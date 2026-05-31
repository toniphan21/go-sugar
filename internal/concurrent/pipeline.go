package concurrent

import (
	"context"
	"sync"
)

func Init[T, V any](ctx context.Context, lists []T, initializer func(T) <-chan V) <-chan V {
	out := make([]<-chan V, len(lists))
	for i, item := range lists {
		out[i] = initializer(item)
	}
	return FanIn(ctx, out...)
}

func OrDone[T any](ctx context.Context, in <-chan T) <-chan T {
	stream := make(chan T)
	go func() {
		defer close(stream)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case stream <- v:
				case <-ctx.Done():
				}
			}
		}
	}()
	return stream
}

func Tee[T any](ctx context.Context, stream <-chan T) (<-chan T, <-chan T) {
	ch1 := make(chan T)
	ch2 := make(chan T)
	go func() {
		defer close(ch1)
		defer close(ch2)

		for v := range OrDone(ctx, stream) {
			var v1, v2 = ch1, ch2
			for i := 0; i < 2; i++ {
				select {
				case v1 <- v:
					v1 = nil
				case v2 <- v:
					v2 = nil
				}
			}
		}
	}()
	return ch1, ch2
}

func FanIn[T any](ctx context.Context, streams ...<-chan T) <-chan T {
	stream := make(chan T)
	var wg sync.WaitGroup
	for _, v := range streams {
		wg.Go(func() {
			for vv := range v {
				select {
				case <-ctx.Done():
					return
				case stream <- vv:
				}
			}
		})
	}
	go func() {
		wg.Wait()
		close(stream)
	}()
	return stream
}

func FanOut[T any](count int, worker func() <-chan T) []<-chan T {
	var outs = make([]<-chan T, count)
	for i := 0; i < count; i++ {
		outs[i] = worker()
	}
	return outs
}

func Parallelize[T any](ctx context.Context, count int, worker func() <-chan T) <-chan T {
	outs := FanOut(count, worker)

	return FanIn(ctx, outs...)
}
