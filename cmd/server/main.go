package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/giornetta/ridl/service"

	"github.com/giornetta/ridl/repository"

	"github.com/giornetta/ridl/cipher"

	"github.com/dgraph-io/badger"

	"github.com/giornetta/ridl/server"
)

func main() {
	// Configuration flags for the application
	port := flag.String("port", "8080", "port where the application is exposed")
	badgerPath := flag.String("badger", "./badger_db", "path for badger database")
	flag.Parse()

	// Open the badger database
	opts := badger.DefaultOptions
	opts.Dir = *badgerPath
	opts.ValueDir = *badgerPath
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	// Close the DB connection when the program ends
	defer db.Close()

	repo := repository.NewBadger(db)

	svc := service.New(cipher.NewAES(), repo)

	srv := server.New(server.Router(svc), *port)
	// Correctly close the server when the program ends
	defer srv.Shutdown(nil)

	// Concurrently run a task to delete expired riddles every hour
	go func() {
		if err := repo.DeleteExpired(); err != nil {
			log.Printf("could not delete expired riddles: %v", err)
		}
		time.Sleep(time.Hour)
	}()

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
