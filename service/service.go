package service

import (
	"context"
	"distributed/registry"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Start(ctx context.Context, port string, reg registry.Registration, registerHandlerFun func(mux *http.ServeMux)) (context.Context, error) {
	mux := http.NewServeMux()
	registerHandlerFun(mux)

	ctx = startService(ctx, port, reg, mux)
	err := registry.RegisterService(mux, reg)
	if err != nil {
		log.Fatal("RegisterService error: ", err)
	}
	return ctx, err
}

func startService(ctx context.Context, port string, reg registry.Registration, mux *http.ServeMux) context.Context {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	ctx, cancel := context.WithCancel(ctx)

	// start http server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%s Service error: %s", reg.ServiceName, err.Error())

			// deregister service
			if err := registry.DeregisterService(reg.ServiceUrl); err != nil {
				log.Fatal("DeregisterService error: ", err)
			}
		}

		cancel()
	}()

	go func() {
		// capture Ctrl+C
		fmt.Printf("%v started.\nPress Ctrl+C to stop.\n", reg.ServiceName)
		sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		<-sigCtx.Done()

		// deregister service
		if err := registry.DeregisterService(reg.ServiceUrl); err != nil {
			log.Fatal("DeregisterService error: ", err)
		}

		// graceful shutdown
		fmt.Printf("%v stopping ...\n", reg.ServiceName)
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
