package main

import (
	"log"
	"net/http"

	"github.com/cruncheon/task-slayer/handlers"
)

func main() {
	// Load players and quests
	handlers.LoadData()

	// HTTP routes
	handlers.LoadRoutes()

	// HTTP Server
	srv := &http.Server{Addr: ":8080"}
	log.Println("Server started at http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
