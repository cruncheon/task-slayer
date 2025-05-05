package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func listQuests(w http.ResponseWriter, r *http.Request) {
	// Parse the template files
	tmpl, err := template.ParseFiles("templates/base.html", "templates/list_quests.html")
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
		tmpl, err := template.ParseFiles("templates/base.html", "templates/create_quest.html")
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
