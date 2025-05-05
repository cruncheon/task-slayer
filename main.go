package main

import (
	"html/template"
	"log"
	"net/http"
)

// Define the Player struct for player data
type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	XP   int64  `json:"xp"`
	Gold int64  `json:"gold"`
}

// Define the Quest struct for quest data
type Quest struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	PlayerID string `json:"player_id"`
	XP       int64  `json:"xp"`
	Gold     int64  `json:"gold"`
	Complete bool   `json:"complete"`
}

type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

var players []Player
var quests []Quest
var items []Item

func main() {
	// Load players and quests
	loadPlayers()
	loadQuests()
	loadItems()

	// HTTP routes
	http.HandleFunc("/", homePage)

	http.HandleFunc("/quests", listQuests)
	http.HandleFunc("/quest/create", createQuest)
	http.HandleFunc("/quest/complete", completeQuest)

	http.HandleFunc("/players", listPlayers)
	http.HandleFunc("/player/create", createPlayer)

	http.HandleFunc("/items", listItems)
	http.HandleFunc("/item/create", createItem)

	// HTTP Server
	srv := &http.Server{Addr: ":8080"}
	log.Println("Server started at http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// Serve the home page
func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, "base.html")
	if err != nil {
		log.Fatal(err)
	}
}
