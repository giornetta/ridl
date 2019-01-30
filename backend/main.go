package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/giornetta/ridl/ridl"
	"github.com/giornetta/ridl/server"
)

func main() {
	svc := ridl.New()

	srv := server.New(svc)

	// Launch the server in a goroutine
	go func() {
		log.Println("listening on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("could not listen: %v", err)
		}
	}()

	// Create a signal channel to listen to Interrupt signals
	// in order to correctly shutdown the server
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig
	log.Println("Shutting down server...")
	srv.Shutdown(nil)
}
