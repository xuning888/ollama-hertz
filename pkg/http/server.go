package http

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xuning888/yoyoyo/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var DefaultShutdownTimeout = time.Minute

type Options func(*Server)

// OnShutdown is a callback function that is called on shutdown.
type OnShutdown func(ctx context.Context)

// Server is an HTTP server which implements graceful shutdown and shutdown hooks
type Server struct {
	shutdownTimeout time.Duration
	srv             *http.Server
	// shutdownHooks
	// Note: you can call AddShutdownHook to add an OnShutdown function, which will be called on shutdown.
	shutdownHooks []OnShutdown
	lg            logger.Logger
}

func WithShutdownTimout(shutdownTimeout time.Duration) Options {
	return func(server *Server) {
		server.shutdownTimeout = shutdownTimeout
	}
}

func NewServer(router *gin.Engine, addr string, opts ...Options) *Server {

	// create a http Server for the purpose of implementing graceful shutdown
	server := &Server{
		shutdownTimeout: DefaultShutdownTimeout,
		srv: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		shutdownHooks: make([]OnShutdown, 0),
		lg:            logger.Named("http-server"),
	}

	for _, opt := range opts {
		opt(server)
	}
	return server
}

func (s *Server) AddShutdownHook(shutdownHook OnShutdown) {
	s.shutdownHooks = append(s.shutdownHooks, shutdownHook)
}

func (s *Server) Serve() {
	defer s.lg.Sync()

	errCh := make(chan error)
	go func() {
		errCh <- s.srv.ListenAndServe()
	}()

	// listen os shutdown signals such as kill -15 (ctrl + c), kill -9
	if err := waitSignal(errCh); err != nil {
		s.lg.Warnf("Received SIGINT %v scheduling shutdown...", err)
	}

	// shutdown gracefully timeout
	timeout, cancelFunc := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancelFunc()

	// call shutdown and execute onShutdownHooks
	if err := s.Shutdown(timeout); err != nil {
		s.lg.Errorf("Server shutdown error: %v\n", err)
		return
	}
	s.lg.Infof("Server shutdown success")
}

func (s *Server) Shutdown(ctx context.Context) (err error) {
	defer s.lg.Sync()

	ch := make(chan struct{})
	// execute shutdown hook
	go s.executeShutdownHook(ctx, ch)
	defer func() {
		select {
		case <-ctx.Done():
			s.lg.Infof("Execute ShutdownHooks timeout: error=%v", ctx.Err())
			return
		case <-ch:
			s.lg.Infof("Execute ShutdownHooks finish")
			return
		}
	}()

	// shutdown http server
	if err = s.srv.Shutdown(ctx); err != nil {
		return err
	}
	return
}

func (s *Server) executeShutdownHook(ctx context.Context, ch chan struct{}) {
	var wg sync.WaitGroup
	for i := range s.shutdownHooks {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			s.shutdownHooks[index](ctx)
		}(i)
	}
	wg.Wait()
	ch <- struct{}{}
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
