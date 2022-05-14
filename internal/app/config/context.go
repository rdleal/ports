package config

import (
	"context"
	"os"
	"os/signal"
)

// Pulled out for testing.
var notify = signal.Notify

// ContextWithGracefulCancellation returns a context.Context
// that is canceled when it receives an os interruption signal.
func ContextWithGracefulCancellation(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	signalCh := make(chan os.Signal, 2)
	notify(signalCh, os.Interrupt, os.Kill)

	go func() {
		defer cancel()

		select {
		case <-signalCh:
		case <-ctx.Done():
		}
	}()

	return ctx
}
