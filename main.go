package main

import (
	"log"
	"net/http"

	"github.com/cruncheon/task-slayer/data"
	"github.com/cruncheon/task-slayer/handlers"
)

func main() {
	// JSON data files
	data.LoadData()

	// HTTP routes
	handlers.LoadRoutes()

	// HTTP Server
	srv := &http.Server{Addr: ":8080"}
	log.Println("Server started on port 8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
