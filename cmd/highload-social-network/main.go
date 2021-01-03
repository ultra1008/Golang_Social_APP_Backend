package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/niklod/highload-social-network/config"
	"github.com/niklod/highload-social-network/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewHTTPServer(cfg.Server)
	srv.Start()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	log.Printf("received signal %s, stopping program...", sig)
	srv.Shutdown()
	signal.Stop(sigCh)
	log.Println("program stopped")
}
