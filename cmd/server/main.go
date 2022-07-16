package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/kosimovsky/tricMe/config"
	"github.com/kosimovsky/tricMe/internal/handlers"
	"github.com/kosimovsky/tricMe/internal/server"
	"github.com/kosimovsky/tricMe/internal/storage"
)

func main() {
	if err := config.InitConfig(); err != nil {
		_ = fmt.Errorf("error while reading config file %v", err.Error())
	}
	logfileFromConfig := viper.GetString("server.logfile")
	logfile, err := os.OpenFile(logfileFromConfig, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error occured opening file : %v, %s", err, logfileFromConfig)
	}
	defer logfile.Close()
	logrus.New()
	logrus.SetOutput(logfile)
	logrus.SetFormatter(new(logrus.JSONFormatter))

	store, _ := storage.NewStorage(&storage.Storage{StorageType: viper.GetString("server.storage")})
	handler := handlers.NewHandler(store)

	srv := server.NewServer()

	if viper.GetBool("server.debug") {
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
		logrus.Printf("Server started in silent mode with loglevel: %v", logrus.GetLevel().String())
	}

	go func() {
		if err := srv.Run(viper.GetString("server.port"), handler.ServeHTTP); err != nil {
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
