package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kosimovsky/tricMe/config"
	"github.com/kosimovsky/tricMe/internal/handlers"
	"github.com/kosimovsky/tricMe/internal/log"
	"github.com/kosimovsky/tricMe/internal/server"
	"github.com/kosimovsky/tricMe/internal/storage"
)

func main() {
	if err := config.InitServerConfig(); err != nil {
		panic(fmt.Sprintf("error while reading config file %s", err.Error()))
	}
	c := config.NewServerConfig()

	logger, err := log.NewLogger(c.Logfile)
	if err != nil {
		fmt.Printf("error occurs while logger initialization: %s", err.Error())
		os.Exit(1)
	}
	defer logger.Close()
	store, err := storage.NewStorage(c.Storage)
	if err != nil {
		logger.Error(err)
	}
	handler := handlers.NewHandler(store)
	err = store.Restore(c.Filename, c.Restore)
	if err != nil {
		logger.Error(err.Error())
	}
	srv := server.NewServer(c.Address, handler.MetricsRouter())

	if c.Debug {
		logger.SetWarnLevel()
		logger.Printf("Server started in debug mode with loglevel: %v", logger.GetLevel().String())
		go store.Output(logger)
	} else {
		logger.Printf("Server started on %s in silent mode with loglevel: %v", c.Address, logger.GetLevel().String())
	}

	if c.Filename != "" && c.StoreInterval > 0 {
		go store.Keep(c.Filename, c.StoreInterval, logger)
	}

	go srv.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	logger.Printf("Recieved a signal %v. Server is shutting down...", sig)
	if err = srv.Shutdown(context.Background()); err != nil {
		logger.Printf("error occured while server shutting down : %s", err.Error())
	}
}
