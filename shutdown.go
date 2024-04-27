package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// OnShutdown is a callback function that is called on shutdown.
type OnShutdown func(ctx context.Context)

// shutdownHooks
// Note: you can call AddShutdownHook to add an OnShutdown function, which will be called on shutdown.
var shutdownHooks []OnShutdown

func AddShutdownHook(hook OnShutdown) {
	shutdownHooks = append(shutdownHooks, hook)
}

// waitSignal
// Node: Listens for operating system shutdown signals such as kill -15 (ctrl + c), kill -9, or
// listens for error signals sent through errCh.
func waitSignal(errCh chan error) error {
	signalToNotify := []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, signalToNotify...)
	select {
	case sig := <-signals:
		switch sig {
		case syscall.SIGTERM:
			// handle force shutdown
			return errors.New(sig.String())
		case syscall.SIGHUP, syscall.SIGINT:
			// handle standard interrupt signals
			return errors.New(sig.String())
		}
	case err := <-errCh:
		// handle error received on errCh
		return err
	}
	return nil
}

func Shutdown(ctx context.Context, srv *http.Server) (err error) {

	ch := make(chan struct{})
	// execute shutdown hook
	go executeShutdownHook(ctx, ch)
	defer func() {
		select {
		case <-ctx.Done():
			log.Printf("Execute ShutdownHooks timeout: error=%v\n", ctx.Err())
			return
		case <-ch:
			log.Println("Execute ShutdownHooks finish")
			return
		}
	}()

	// shutdown http server
	if err = srv.Shutdown(ctx); err != nil {
		return err
	}
	return
}

// executeShutdownHook
// Note: Execute each OnShutdown function in sequence.
func executeShutdownHook(ctx context.Context, ch chan struct{}) {
	wg := sync.WaitGroup{}
	for i := range shutdownHooks {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			shutdownHooks[index](ctx)
		}(i)
	}
	wg.Wait()
	ch <- struct{}{}
}
