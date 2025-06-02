package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cruncheon/task-slayer/data"
	"github.com/cruncheon/task-slayer/templates"
)

func listPlayers(w http.ResponseWriter, r *http.Request) {
	pageData := struct {
		Players []data.Player
		Quests  []data.Quest
	}{
		data.Players,
		data.Quests,
	}

	renderTemplate(w, templates.ListPlayers, pageData)
}

// Create player
func createPlayer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		pageData := struct {
			Players []data.Player
		}{
			data.Players,
		}

		renderTemplate(w, templates.CreatePlayer, pageData)

	case "POST":
		r.ParseForm()

		name := r.FormValue("name")
		xp := int64(0)
		gold := int64(0)

		playerID := fmt.Sprintf("player%d", len(data.Players)+1)

		newPlayer := data.Player{
			ID:   playerID,
			Name: name,
			XP:   xp,
			Gold: gold,
		}

		data.Players = append(data.Players, newPlayer)

		data.SavePlayers()

		log.Printf("%v - %v created", newPlayer.ID, newPlayer.Name)
		http.Redirect(w, r, "/players", http.StatusSeeOther)

	default:
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
	}
}

// Edit player
func editPlayer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/player/edit/"):]

	player := data.GetPlayer(id)
	if player == nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		pageData := struct {
			Player *data.Player
		}{
			player,
		}

		renderTemplate(w, templates.EditPlayer, pageData)

	case "POST":
		r.ParseForm()

		name := r.FormValue("name")

		xp, err := parseFormInt(r, "xp")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		gold, err := parseFormInt(r, "gold")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		player.Name = name
		player.XP = xp
		player.Gold = gold

		data.SavePlayers()

		log.Printf("%v - %v updated", player.ID, player.Name)
		http.Redirect(w, r, "/players", http.StatusSeeOther)

	default:
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
	}
}

// Delete player
func deletePlayer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/player/delete/"):]
	if err := data.DeletePlayer(id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete player: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("%v deleted", id)
	http.Redirect(w, r, "/players", http.StatusSeeOther)
}
