package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/ynuraddi/t-tsarka/config"
	"github.com/ynuraddi/t-tsarka/pkg/logger"
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

	service := service.New(config, logger)

	server := transport.New(config, logger, service)

	log.Fatalln(server.Start(ctx))
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
