package main

import (
	"context"
	"fmt"
	"github.com/kosimovsky/tricMe/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const serverPort = 8080

func main() {

	srv := server.NewServer()

	go func() {
		mux := http.DefaultServeMux
		mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("server: %s /\n", r.Method)
		})
		if err := srv.Run("8080", mux); err != nil {
			log.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	log.Printf("Recieved a signal %v. Server is shutting down...", sig)
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occured while server shutting down : %s", err.Error())
	}
}
