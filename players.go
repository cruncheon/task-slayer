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

func listPlayers(w http.ResponseWriter, r *http.Request) {
	// Parse the template files
	tmpl, err := template.ParseFiles("templates/base.html", "templates/list_players.html")
	if err != nil {
		log.Fatal(err)
	}

	// Define data structure for list of players
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

// Create player
func createPlayer(w http.ResponseWriter, r *http.Request) {
	// If Get request, render create player page
	if r.Method == http.MethodGet {
		// Parse the template files
		tmpl, err := template.ParseFiles("templates/base.html", "templates/create_player.html")
		if err != nil {
		}

		// Define data structure for player creation form
		data := struct {
			Players []Player
		}{
			Players: players,
		}

		// Execute template and render create player page
		err = tmpl.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			log.Fatal(err)
		}

		// If Post request, create player and redirect back to list players page
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		name := r.FormValue("name")
		xp := int64(0)
		gold := int64(0)

		// Set player ID by getting current amount of players and adding +1
		playerID := fmt.Sprintf("player%d", len(players)+1)

		// Structure new player details
		newPlayer := Player{
			ID:   playerID,
			Name: name,
			XP:   xp,
			Gold: gold,
		}

		// Add new player to players slice
		players = append(players, newPlayer)

		// Save changes to JSON file
		savePlayers()

		log.Printf("%v - %v created", newPlayer.ID, newPlayer.Name)
		http.Redirect(w, r, "/players", http.StatusSeeOther)
	}
}

// Edit player
func editPlayer(w http.ResponseWriter, r *http.Request) {
	// Extract player ID from URL path
	id := r.URL.Path[len("/player/edit/"):]

	// Get the player by ID
	player := getPlayer(id)
	if player == nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		// Render edit player page
		tmpl, err := template.ParseFiles("templates/base.html", "templates/edit_player.html")
		if err != nil {
			log.Fatal(err)
		}

		data := struct {
			Player *Player
		}{
			Player: player,
		}

		err = tmpl.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			log.Fatal(err)
		}
	} else if r.Method == http.MethodPost {
		// Update the player details
		r.ParseForm()

		name := r.FormValue("name")

		// Convert XP form input from string to int
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

		// Update item details
		player.Name = name
		player.XP = xp
		player.Gold = gold

		// Save changes to JSON file
		savePlayers()

		log.Printf("%v - %v updated", player.ID, player.Name)
		http.Redirect(w, r, "/players", http.StatusSeeOther)
	}
}

// Load players
func loadPlayers() {
	filePath := "data/players.json"

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the file if it does not exist
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Initialize an empty slice of players and write to the file
		var initialPlayers []Player
		err = json.NewEncoder(file).Encode(initialPlayers)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// Open the file if it exists
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decode JSON from player file
	err = json.NewDecoder(file).Decode(&players)
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

// Get player details by ID
func getPlayer(id string) *Player {
	for i := range players {
		if players[i].ID == id {
			return &players[i]
		}
	}
	return nil
}
