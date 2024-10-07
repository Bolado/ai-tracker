package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Bolado/ai-tracker/database"
	router "github.com/Bolado/ai-tracker/router"
	"github.com/Bolado/ai-tracker/watcher"
)

func main() {

	//get env variables
	interval, _ := strconv.Atoi(os.Getenv("WATCHER_INTERVAL"))

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
	go func() {
		for {
			if err := watcher.StartWatcher(); err != nil {
				log.Fatalf("Failed to start watcher: %v\n", err)
			}
			if interval == 0 {
				time.Sleep(15 * time.Minute)
			} else {
				time.Sleep(time.Duration(interval) * time.Minute)
			}
		}
	}()
	log.Println("Watcher started")

	// start the router and listen for requests
	if err := router.StartRouter(); err != nil {
		log.Fatalf("Failed to start router: %v\n", err)
	}
}
