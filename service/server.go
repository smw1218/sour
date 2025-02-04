package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
)

func Run(ctx context.Context, si ServiceInterface, handler http.Handler) {
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	//defer stop() // I don't think stop is needed in this context

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", si.DefaultPort()),
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		slog.Info("Got done signal, shutting down HTTP Server")
		err := srv.Shutdown(ctx)
		if err != nil {
			slog.Error("Server Shutdown Failed", "error", err)
		}
	}()

	// This call blocks until shutdown. After shutdown, we know we have finished processing
	// all open hanlders
	slog.Info(fmt.Sprintf("Running Server on %v", srv.Addr))
	err := srv.ListenAndServe()
	slog.Info("ListenAndServe exiting", "error", err)

	// Once the handlers have done their work, if there's any other waiting for work to complete,
	// it can happen in the service's Shutdown handler
	slog.Info("Calling Service Shutdown")
	err = si.Shutdown()
	if err != nil {
		slog.Error("Failed service Shutdown", "error", err)
	}
}
