package main

import (
	"BloomFilter/handlers"
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	l := log.New(os.Stdout, "set-api", log.LstdFlags)

	sh:= handlers.NewSets(l)

	// create serve mux and register the handlers
	sm := mux.NewRouter()

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/sets/{set-name}", sh.AddSet)
	postRouter.Use(sh.MiddlewareValidateSet)

	sm.HandleFunc("/sets/{set-name}/{item-name}",sh.PutItem).Methods(http.MethodPut)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/sets/{set-name}/stats",sh.GetStats)
	getRouter.HandleFunc("/sets/{set-name}/items/{item-name}",sh.GetItem)

	// create the server
	s := http.Server{
		Addr:         "0.0.0.0:8080",      // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server on port 8080")
		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)


}
