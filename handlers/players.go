package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cruncheon/task-slayer/data"
	"github.com/cruncheon/task-slayer/templates"
)

func listPlayers(w http.ResponseWriter, r *http.Request) {
	// Define data structure for list of players
	pageData := struct {
		Players []data.Player
		Quests  []data.Quest
	}{
		data.Players,
		data.Quests,
	}

	// Render list players page
	renderTemplate(w, templates.ListPlayers, pageData)
}

// Create player
func createPlayer(w http.ResponseWriter, r *http.Request) {
	// If Get request, render create player page
	if r.Method == http.MethodGet {
		// Define data structure for player creation form
		pageData := struct {
			Players []data.Player
		}{
			data.Players,
		}

		// Render create player page
		renderTemplate(w, templates.CreatePlayer, pageData)

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

	// Handle get requests
	if r.Method == http.MethodGet {
		pageData := struct {
			Player *data.Player
		}{
			player,
		}

		// Render edit player page
		renderTemplate(w, templates.EditPlayer, pageData)

		// Handle post requests
	} else if r.Method == http.MethodPost {
		// Update the player details
		r.ParseForm()

		name := r.FormValue("name")

		// Convert XP form input from string to int
		xp, err := parseFormInt(r, "xp")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Convert gold form input from string to int
		gold, err := parseFormInt(r, "gold")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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

func deletePlayer(w http.ResponseWriter, r *http.Request) {
	// Extract player ID from URL path
	id := r.URL.Path[len("/player/delete/"):]
	if err := data.DeletePlayer(id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete player: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("%v deleted", id)
	http.Redirect(w, r, "/players", http.StatusSeeOther)
}
