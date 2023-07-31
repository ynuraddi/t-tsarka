package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/ynuraddi/t-tsarka/config"
	"github.com/ynuraddi/t-tsarka/pkg/logger"
	"github.com/ynuraddi/t-tsarka/repository"
	"github.com/ynuraddi/t-tsarka/service"
	"github.com/ynuraddi/t-tsarka/transport"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	osC := gracefullShutdown(cancel)

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln(err)
	}

	fileLogs, err := os.OpenFile(config.LogPath, os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatalln(err)
	}
	logger := logger.NewLogger(fileLogs, logger.Level(config.LogLevel), osC)

	repository, err := repository.New(config, logger)
	if err != nil {
		logger.Fatal("failed init repo", err)
		return
	}

	service, err := service.New(config, logger, repository)
	if err != nil {
		logger.Fatal("failed init service", err)
		return
	}
	logger.Debug("service succes inited")

	server := transport.New(config, logger, service)
	logger.Debug("server success inited")

	logger.Info("start server")
	if err := server.Start(ctx); err != nil {
		log.Println(err)
	}
}

func gracefullShutdown(c context.CancelFunc) chan os.Signal {
	osC := make(chan os.Signal, 2)
	signal.Notify(osC, os.Interrupt)

	go func() {
		log.Println(<-osC)
		c()
	}()

	return osC
}
