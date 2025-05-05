package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
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

var players []Player
var quests []Quest

func main() {
	// Load players and quests
	loadPlayers()
	loadQuests()

	// HTTP routes
	http.HandleFunc("/", homePage)                    // Handle home page
	http.HandleFunc("/quests", listQuests)            // Hanlde quest list
	http.HandleFunc("/quest/create", createQuest)     // Handle quest creation
	http.HandleFunc("/quest/complete", completeQuest) // Handle quest completion

	// HTTP Server
	srv := &http.Server{Addr: ":8080"}
	log.Println("Server started at http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// Load players
func loadPlayers() {
	file, err := os.Open("data/players.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&players)
	if err != nil {
		log.Fatal(err)
	}
}

// Load quests
func loadQuests() {
	file, err := os.Open("data/quests.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&quests)
	if err != nil {
		log.Fatal(err)
	}
}

// Serve the index page
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

// Get quest details by ID
func getQuest(id string) *Quest {
	for i := range quests {
		if quests[i].ID == id {
			return &quests[i]
		}
	}
	return nil
}

// Get player details by ID
func getPlayer(id string) *Player {
	for i := range players {
		if players[i].ID == id {
			return &players[i]
		}
	}
	return nil
}

// Save quests
func saveQuests() {
	questFile, err := os.Create("data/quests.json")
	if err != nil {
		log.Fatal(err)
	}
	defer questFile.Close()

	encoder := json.NewEncoder(questFile)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(quests)
	if err != nil {
		log.Fatal(err)
	}
}

// Save players
func savePlayers() {
	playerFile, err := os.Create("data/players.json")
	if err != nil {
		log.Fatal(err)
	}
	defer playerFile.Close()

	encoder := json.NewEncoder(playerFile)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(players)
	if err != nil {
		log.Fatal(err)
	}
}
