package main

import (
	"flag"
	"goapp/pkg/config"
	"goapp/pkg/infrastructure/web"
	"log"
)

func main() {

	// LOAD DEFAULT CONFIG
	cfg := config.LoadConfig()

	flag.StringVar(&cfg.ServerPort, "port", cfg.ServerPort, "Port to listen on")
	flag.Parse()

	//Start server
	if err := web.StartWebServer(cfg); err != nil {
		log.Fatalf("server exited %v", err)
	}

}
