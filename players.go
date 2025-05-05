package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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

		// Set player ID by getting current amount of quests and adding +1
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
