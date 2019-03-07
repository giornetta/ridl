package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/giornetta/ridl/repository"

	"github.com/giornetta/ridl/cipher"

	"github.com/dgraph-io/badger"

	"github.com/giornetta/ridl/ridl"
	"github.com/giornetta/ridl/server"
)

func main() {
	opts := badger.DefaultOptions
	opts.Dir = "./badger_db"
	opts.ValueDir = "./badger_db"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.New(db)

	svc := ridl.NewService(cipher.NewAES(), repo)

	srv := server.New(svc)
	defer srv.Shutdown(nil)

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
}
