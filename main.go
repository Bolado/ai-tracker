package main

import (
	"log"

	"github.com/Bolado/ai-tracker/database"
	router "github.com/Bolado/ai-tracker/router"
	"github.com/Bolado/ai-tracker/watcher"
)

func main() {
	// initialize the database
	if err := database.StartDabase(); err != nil {
		log.Fatalf("Failed to start database: %v\n", err)
	}
	log.Println("Database started")

	// load articles from the database
	if err := watcher.LoadArticles(); err != nil {
		log.Fatalf("Failed to load articles: %v\n", err)
	}
	log.Println("Articles loaded")

	// start the watcher
	if err := watcher.StartWatcher(); err != nil {
		log.Fatalf("Failed to start watcher: %v\n", err)
	}
	log.Println("Watcher started")

	// start the router and listen for requests
	if err := router.StartRouter(); err != nil {
		log.Fatalf("Failed to start router: %v\n", err)
	}
}
