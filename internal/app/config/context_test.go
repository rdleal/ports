package config

import (
	"context"
	"os"
	"testing"
	"time"
)

func stubNotify(stub func(chan<- os.Signal, ...os.Signal)) func() {
	orig := notify
	notify = stub
	return func() { notify = orig }
}

func TestContextWithGracefulCancellation(t *testing.T) {
	defer stubNotify(func(c chan<- os.Signal, _ ...os.Signal) {
		select {
		case c <- os.Interrupt:
		case <-time.After(500 * time.Millisecond):
			t.Fatal("got timed out waiting for signal to be consumed")
		}
	})()

	ctx := ContextWithGracefulCancellation(context.Background())

	select {
	case <-ctx.Done():
	case <-time.After(500 * time.Millisecond):
		t.Fatal("got timed out waiting for context to be canceled")

	}
}
