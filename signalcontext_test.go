package signalcontext_test

import (
	"context"
	"log"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/sethvargo/go-signalcontext"
)

func TestWrap(t *testing.T) {
	ctx, cancel := signalcontext.Wrap(context.Background(), syscall.SIGUSR1)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("context should not be done")
	case <-time.After(10 * time.Millisecond):
	}

	if err := syscall.Kill(syscall.Getpid(), syscall.SIGUSR1); err != nil {
		t.Fatal("failed to signal")
	}

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(10 * time.Millisecond):
		t.Fatal("context should have been done")
	}
}

func TestWrapAll(t *testing.T) {
	ctx, cancel := signalcontext.WrapAll(context.Background())
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("context should not be done")
	case <-time.After(10 * time.Millisecond):
	}

	if err := syscall.Kill(syscall.Getpid(), syscall.SIGINFO); err != nil {
		t.Fatal("failed to signal")
	}

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(10 * time.Millisecond):
		t.Fatal("context should have been done")
	}
}

func ExampleOnInterrupt() {
	ctx, cancel := signalcontext.OnInterrupt()
	defer cancel()

	s := &http.Server{
		Addr: ":8080",
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for CTRL+C
	<-ctx.Done()

	// Stop the server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
