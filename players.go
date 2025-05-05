package main

import (
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
