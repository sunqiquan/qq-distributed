package main

import (
	"context"
	"distributed/registry"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// capture Ctrl+C
	fmt.Printf("Registry Service started. Press Ctrl+C to stop.\n")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// start http server
	mux := http.NewServeMux()
	mux.Handle("/services", &registry.RegistryService{})
	srv := &http.Server{
		Addr:    ":" + registry.ServerPort,
		Handler: mux,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Registry service failed to start: ", err)
			stop()
		}
	}()

	<-ctx.Done()
	fmt.Println("Shutting down registry service ...")

	// stop http server
	srv.Shutdown(ctx)
	fmt.Println("Registry service stopped.")
}
