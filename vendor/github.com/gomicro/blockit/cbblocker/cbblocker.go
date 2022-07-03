package cbblocker

import (
	"context"
	"time"
)

type Blocker struct {
	callback func() error
	duration time.Duration
}

func New(callback func() error, poll time.Duration) *Blocker {
	return &Blocker{
		callback: callback,
		duration: poll,
	}
}

func (g *Blocker) Blockit() <-chan struct{} {
	return g.BlockitWithContext(context.Background())
}

func (g *Blocker) BlockitWithContext(ctx context.Context) <-chan struct{} {
	ctxDone := ctx.Done()
	out := make(chan struct{})
	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			err := g.callback()
			if err != nil {
				<-time.After(g.duration)
				continue
			}

			done <- struct{}{}
			break
		}
	}()

	go func() {
		defer close(out)

		select {
		case <-ctxDone:
		case <-done:
		}

		out <- struct{}{}
	}()

	return out
}
