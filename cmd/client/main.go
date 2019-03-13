package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/amirhosseinab/sfs"
	"github.com/giornetta/ridl/server"
)

func main() {
	// Configuration flags for the application
	port := flag.String("port", "3000", "port where the application is exposed")
	flag.Parse()

	h := sfs.New(http.Dir("static"), indexHandler)

	srv := server.New(h, *port)
	// Correctly close the server when the program ends
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
