package main

import (
	"fmt"
	"log"
	"net/http"
	"medium-go-project/internal/api/routes"
	"medium-go-project/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	router := routes.Setup(cfg)
	
	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}