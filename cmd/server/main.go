package main

import (
	"catecard/pkg/config"
	"catecard/pkg/infrastructure/web"
	"flag"
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
