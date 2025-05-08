package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cruncheon/task-slayer/data"
	"github.com/cruncheon/task-slayer/templates"
)

func listQuests(w http.ResponseWriter, r *http.Request) {
	// Define data structure for list of quests
	pageData := struct {
		Players []data.Player
		Quests  []data.Quest
	}{
		Players: data.Players,
		Quests:  data.Quests,
	}

	// Render list quests page
	renderTemplate(w, templates.ListQuests, pageData)
}

// Create quest
func createQuest(w http.ResponseWriter, r *http.Request) {
	// If Get request, render create quest page
	if r.Method == http.MethodGet {
		// Define data structure for quest creation form
		pageData := struct {
			Players []data.Player
			Quests  []data.Quest
		}{
			Players: data.Players,
			Quests:  data.Quests,
		}

		// Render create quest page
		renderTemplate(w, templates.CreateQuest, pageData)

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
		questID := fmt.Sprintf("quest%d", len(data.Quests)+1)

		// Structure new quest details
		newQuest := data.Quest{
			ID:       questID,
			Title:    title,
			PlayerID: playerID,
			XP:       xp,
			Gold:     gold,
			Complete: false,
		}

		// Add new quest to quests slice and save
		data.Quests = append(data.Quests, newQuest)
		data.SaveQuests()

		log.Printf("%v created for %v", newQuest.ID, newQuest.PlayerID)
		http.Redirect(w, r, "/quests", http.StatusSeeOther)
	}
}

// Edit quest
func editQuest(w http.ResponseWriter, r *http.Request) {
	// Extract quest ID from URL path
	id := r.URL.Path[len("/quest/edit/"):]

	// Get the quest by ID
	quest := data.GetQuest(id)
	if quest == nil {
		http.Error(w, "Quest not found", http.StatusNotFound)
		return
	}

	// Handle get requests
	if r.Method == http.MethodGet {
		data := struct {
			Quest *data.Quest
		}{
			Quest: quest,
		}

		// Render edit quest page
		renderTemplate(w, templates.EditQuest, data)

		// Handle post requests
	} else if r.Method == http.MethodPost {
		// Update the quest details
		r.ParseForm()

		title := r.FormValue("title")

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

		// Update quest details
		quest.Title = title
		quest.XP = xp
		quest.Gold = gold

		// Save changes to JSON file
		data.SaveQuests()

		log.Printf("%v - %v updated", quest.ID, quest.Title)
		http.Redirect(w, r, "/quests", http.StatusSeeOther)
	}
}

// Complete quest
func completeQuest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		// Get quest ID
		questID := r.FormValue("quest_id")
		quest := data.GetQuest(questID)
		if quest == nil {
			http.Error(w, "Quest not found", http.StatusNotFound)
			return
		}

		// Get player ID
		player := data.GetPlayer(quest.PlayerID)
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
		data.SaveQuests()
		data.SavePlayers()

		log.Printf("%s completed %v", player.ID, quest.ID)
		http.Redirect(w, r, "/quests", http.StatusSeeOther)
	}
}
