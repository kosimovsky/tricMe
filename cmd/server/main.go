package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kosimovsky/tricMe/config"
	"github.com/kosimovsky/tricMe/internal/handlers"
	"github.com/kosimovsky/tricMe/internal/server"
	"github.com/kosimovsky/tricMe/internal/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.InitServerConfig(); err != nil {
		_ = fmt.Errorf("error while reading config file %v", err.Error())
	}
	c := storage.ReadConfig()
	logfile, err := os.OpenFile(c.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error occured opening file : %v, %s", err, c.Logfile)
	}
	defer logfile.Close()
	logrus.New()
	logrus.SetOutput(logfile)
	logrus.SetFormatter(new(logrus.JSONFormatter))

	store, _ := storage.NewStorage(&storage.Storage{StorageType: c.Storage})
	handler := handlers.NewHandler(store)
	err = store.Restore(c.Filename, c.Restore)
	if err != nil {
		logrus.Error(err.Error())
	}
	srv := server.NewServer()

	if c.Debug {
		logrus.SetLevel(logrus.WarnLevel)
		logrus.Printf("Server started in debug mode with loglevel: %v", logrus.GetLevel().String())
		ctx := context.Background()
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			for {
				select {
				case <-ticker.C:
					err := store.Output()
					if err != nil {
						logrus.Printf("error output: %s", err.Error())
					}
				case <-ctx.Done():
					ticker.Stop()
				}
			}
		}()
	} else {
		logrus.Printf("Server started on %s in silent mode with loglevel: %v", c.Address, logrus.GetLevel().String())
	}

	if c.Filename != "" && c.StoreInterval > 0 {
		ctx := context.Background()
		go func() {
			ticker := time.NewTicker(c.StoreInterval)
			for {
				select {
				case <-ticker.C:
					err := store.Keep(c.Filename)
					if err != nil {
						logrus.Printf("error storing metrics to file: %s", err.Error())
					}
				case <-ctx.Done():
					ticker.Stop()
				}
			}
		}()
	}

	go func() {
		if err := srv.Run(c.Address, handler.MetricsRouter()); err != nil {
			log.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	log.Printf("Recieved a signal %v. Server is shutting down...", sig)
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("error occured while server shutting down : %s", err.Error())
	}
}
