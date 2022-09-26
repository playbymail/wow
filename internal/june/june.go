/*
 * june - a nice package to gracefully shutdown an HTTP server.
 *
 * Copyright (c) 2022 Clavin June
 */

// Package june implements the graceful HTTP server shutdown from
// https://clavinjune.dev/en/blogs/golang-http-server-graceful-shutdown/.
//
// Expects to be called something like this:
//    func main() {
//        log.SetFlags(log.Lshortfile)
//
//        s := &http.Server{}
//        go june.StartServer(s)
//
//        stopCh, closeCh := june.CreateChannel()
//        defer closeCh()
//        log.Println("notified:", <-stopCh)
//
//        june.Shutdown(context.Background(), s)
//    }
//
package june

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// StartServer runs a server.
// It is safe to run from a go routine.
func StartServer(server *http.Server) {
	log.Println("application started")
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
	log.Println("application stopped gracefully")
}

// ShutdownServer gracefully stops a running server.
// Timeout should be a duration like `5*time.Second`.
func ShutdownServer(c context.Context, server *http.Server, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("application has shut down")
}

// CreateChannel returns channels for catching signals.
// (sort of like trap?)
func CreateChannel() (chan os.Signal, func()) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	return stopCh, func() {
		close(stopCh)
	}
}
