package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/cruncheon/task-slayer/data"
	"github.com/cruncheon/task-slayer/templates"
)

func listPlayers(w http.ResponseWriter, r *http.Request) {
	// Parse the template files
	tmpl, err := template.ParseFiles(templates.Base, templates.ListPlayers)
	if err != nil {
		log.Fatal(err)
	}

	// Define data structure for list of players
	data := struct {
		Players []data.Player
		Quests  []data.Quest
	}{
		Players: data.Players,
		Quests:  data.Quests,
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
		tmpl, err := template.ParseFiles(templates.Base, templates.CreatePlayer)
		if err != nil {
		}

		// Define data structure for player creation form
		data := struct {
			Players []data.Player
		}{
			Players: data.Players,
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
		playerID := fmt.Sprintf("player%d", len(data.Players)+1)

		// Structure new player details
		newPlayer := data.Player{
			ID:   playerID,
			Name: name,
			XP:   xp,
			Gold: gold,
		}

		// Add new player to players slice
		data.Players = append(data.Players, newPlayer)

		// Save changes to JSON file
		data.SavePlayers()

		log.Printf("%v - %v created", newPlayer.ID, newPlayer.Name)
		http.Redirect(w, r, "/players", http.StatusSeeOther)
	}
}

// Edit player
func editPlayer(w http.ResponseWriter, r *http.Request) {
	// Extract player ID from URL path
	id := r.URL.Path[len("/player/edit/"):]

	// Get the player by ID
	player := data.GetPlayer(id)
	if player == nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		// Render edit player page
		tmpl, err := template.ParseFiles(templates.Base, templates.EditPlayer)
		if err != nil {
			log.Fatal(err)
		}

		data := struct {
			Player *data.Player
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

		// Update player details
		player.Name = name
		player.XP = xp
		player.Gold = gold

		// Save changes to JSON file
		data.SavePlayers()

		log.Printf("%v - %v updated", player.ID, player.Name)
		http.Redirect(w, r, "/players", http.StatusSeeOther)
	}
}
