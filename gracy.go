package gracy

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const defaultTimeout = 10 * time.Second

var defaultSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

type CallbackFunc func() error

var gracy *Gracy

type Gracy struct {
	stop      chan os.Signal
	mu        sync.RWMutex
	callbacks []CallbackFunc
}

func init() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, defaultSignals...)

	gracy = &Gracy{stop: stop}
}

func AddCallback(f CallbackFunc) {
	gracy.mu.Lock()
	gracy.callbacks = append(gracy.callbacks, f)
	gracy.mu.Unlock()
}

func Wait() error {
	select {
	case <-gracy.stop:
	}

	return GracefulShutdown()
}

func GracefulShutdown() error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, defaultSignals...)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for _, f := range gracy.callbacks {
			err := f()
			if err != nil {
				_ = err // todo handle err
			}
		}
	}()

	select {
	case <-done:
		return nil
	case <-stop:
		return errors.New("gracy force stopped")
	case <-ctx.Done():
		return errors.New("gracy waiting timeout")
	}
}
