package main

import (
	"IB2/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	serv := service.New()
	s := &http.Server{
		Addr:    ":13999", //"0.0.0.0:%d", cfg.Port),
		Handler: serv.GetHandler(),
	}
	s.SetKeepAlivesEnabled(true)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		log.Printf("starting http server at %d", 13999)
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}

	}()

	gracefullyShutdown(ctx, cancel, s)
}

func gracefullyShutdown(ctx context.Context, cancel context.CancelFunc, server *http.Server) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)
	<-ch
	if err := server.Shutdown(ctx); err != nil {
		log.Print(err)
	}
	cancel()
}
