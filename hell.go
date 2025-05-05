package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
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
	file, err := os.Open("players.json")
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
	file, err := os.Open("quests.json")
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
	tmpl, err := template.ParseFiles("base.html", "index.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, "base.html")
	if err != nil {
		log.Fatal(err)
	}
}

func listQuests(w http.ResponseWriter, r *http.Request) {
	// Parse the template files
	tmpl, err := template.ParseFiles("base.html", "list_quests.html")
	if err != nil {
		log.Fatal(err)
	}

	// Define data structure for list of quests
	data := struct {
		Players []Player
		Quests  []Quest
	}{
		Players: players,
		Quests:  quests,
	}

	// Execute the template with the data structure
	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Fatal(err)
	}
}

// Create quest
func createQuest(w http.ResponseWriter, r *http.Request) {
	// If Get request, render create quest page
	if r.Method == http.MethodGet {
		// Parse the template files
		tmpl, err := template.ParseFiles("base.html", "create_quest.html")
		if err != nil {
		}

		// Define data structure for quest creation form
		data := struct {
			Players []Player
			Quests  []Quest
		}{
			Players: players,
			Quests:  quests,
		}

		// Execute template and render create quest page
		err = tmpl.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			log.Fatal(err)
		}

		// If Post request, create quest and redirect back to list quests page
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		title := r.FormValue("title")
		playerID := r.FormValue("player_id")

		// Convert xp form input from string to int
		xp, err := strconv.ParseInt(r.FormValue("xp"), 10, 64)
		if err != nil {
			log.Printf("Failed to parse XP: %v", err)
			http.Error(w, "Invalid XP value", http.StatusBadRequest)
			return
		}

		// Convert gold form input from string to int
		gold, err := strconv.ParseInt(r.FormValue("gold"), 10, 64)
		if err != nil {
			log.Printf("Failed to parse Gold: %v", err)
			http.Error(w, "Invalid Gold value", http.StatusBadRequest)
			return
		}

		// Set quest ID by getting current amount of quests and adding +1
		questID := fmt.Sprintf("quest%d", len(quests)+1)

		// Structure new quest details
		newQuest := Quest{
			ID:       questID,
			Title:    title,
			PlayerID: playerID,
			XP:       xp,
			Gold:     gold,
			Complete: false,
		}

		// Add new quest to quests slice
		quests = append(quests, newQuest)

		// Save changes to JSON file
		saveQuests()

		log.Printf("%v created for %v", newQuest.ID, newQuest.PlayerID)
		http.Redirect(w, r, "/quests", http.StatusSeeOther)
	}
}

// Complete quest
func completeQuest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		// Get quest ID
		questID := r.FormValue("quest_id")
		quest := getQuest(questID)
		if quest == nil {
			http.Error(w, "Quest not found", http.StatusNotFound)
			return
		}

		// Get player ID
		player := getPlayer(quest.PlayerID)
		if player == nil {
			http.Error(w, "Player not found", http.StatusNotFound)
			return
		}

		// Mark quest as complete
		quest.Complete = true

		// Add quest XP and Gold to player's XP and Gold
		player.XP += quest.XP
		player.Gold += quest.Gold

		// Save changes
		saveQuests()
		savePlayers()

		log.Printf("%s completed %v", player.ID, quest.ID)
		http.Redirect(w, r, "/quests", http.StatusSeeOther)
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
	questFile, err := os.Create("quests.json")
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
	playerFile, err := os.Create("players.json")
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
