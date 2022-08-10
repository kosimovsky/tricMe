package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kosimovsky/tricMe/config"
	"github.com/kosimovsky/tricMe/internal/handlers"
	"github.com/kosimovsky/tricMe/internal/logger"
	"github.com/kosimovsky/tricMe/internal/server"
	"github.com/kosimovsky/tricMe/internal/storage"
)

func main() {
	if err := config.InitServerConfig(); err != nil {
		fmt.Printf("error while reading config file %s", err.Error())
	}
	c := config.ServerConfig()

	logMe := logger.NewLogger()
	logMe.Default(c.Logfile)

	store, err := storage.NewStorage(&storage.Storage{StorageType: c.Storage})
	if err != nil {
		logMe.Error(err)
	}
	handler := handlers.NewHandler(store)
	err = store.Restore(c.Filename, c.Restore)
	if err != nil {
		logMe.Error(err.Error())
	}
	srv := server.NewServer()

	if c.Debug {
		logMe.SetWarnLevel()
		logMe.Printf("Server started in debug mode with loglevel: %v", logMe.GetLevel().String())
		ctx := context.Background()
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			for {
				select {
				case <-ticker.C:
					err = store.Output()
					if err != nil {
						logMe.Printf("error output: %s", err.Error())
					}
				case <-ctx.Done():
					ticker.Stop()
				}
			}
		}()
	} else {
		logMe.Printf("Server started on %s in silent mode with loglevel: %v", c.Address, logMe.GetLevel().String())
	}

	if c.Filename != "" && c.StoreInterval > 0 {
		ctx := context.Background()
		go func() {
			ticker := time.NewTicker(c.StoreInterval)
			select {
			case <-ticker.C:
				err = store.Keep(c.Filename)
				if err != nil {
					logMe.Printf("error storing metrics to file: %s", err.Error())
				}
			case <-ctx.Done():
				ticker.Stop()
			}
		}()
	}

	go func() {
		if err = srv.Run(c.Address, handler.MetricsRouter()); err != nil {
			logMe.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	logMe.Printf("Recieved a signal %v. Server is shutting down...", sig)
	if err = srv.Shutdown(context.Background()); err != nil {
		logMe.Printf("error occured while server shutting down : %s", err.Error())
	}
}
