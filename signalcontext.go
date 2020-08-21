// Package signalcontext creates context.Contexts that cancel on os.Signals.
//
//     ctx, cancel := signalcontext.OnInterrupt()
//     defer cancel()
//
package signalcontext

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// OnInterrupt creates a new context that cancels on SIGINT or SIGTERM.
func OnInterrupt() (context.Context, func()) {
	return wrap(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}

// On creates a new context that cancels on the given signals.
func On(sig os.Signal, signals ...os.Signal) (context.Context, func()) {
	return wrap(context.Background(), append(signals, sig)...)
}

// OnAll creates a new context that cancels on all supported signals.
func OnAll() (context.Context, func()) {
	return wrap(context.Background())
}

// Wrap creates a new context that cancels on the given signals. It wraps the
// provided context.
func Wrap(ctx context.Context, sig os.Signal, signals ...os.Signal) (context.Context, func()) {
	return wrap(ctx, append(signals, sig)...)
}

// WrapAll creates a new context that cancels on all supported signals. It wraps the
// provided context.
func WrapAll(ctx context.Context) (context.Context, func()) {
	return wrap(ctx)
}

func wrap(ctx context.Context, signals ...os.Signal) (context.Context, func()) {
	ctx, closer := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)

	go func() {
		select {
		case <-c:
			closer()
		case <-ctx.Done():
		}
	}()

	return ctx, closer
}
