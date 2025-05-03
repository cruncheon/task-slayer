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

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	XP   int64  `json:"xp"`
	Gold int64  `json:"gold"`
}

type Quest struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	UserID   string `json:"user_id"`
	XP       int64  `json:"xp"`
	Gold     int64  `json:"gold"`
	Complete bool   `json:"complete"`
}

var users []User

var quests []Quest

func main() {
	// Load users and quests
	loadUsers()
	loadQuests()

	// HTTP routes
	http.HandleFunc("/", serveGame)                 // Serve the game page
	http.HandleFunc("/create-quest", createQuest)   // Handle quest creation
	http.HandleFunc("/mark-complete", markComplete) // Handle quest completion
	http.HandleFunc("/quests", listQuests)          // List quests

	// HTTP Server
	srv := &http.Server{Addr: ":8080"}
	log.Println("Server started at http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// Load users
func loadUsers() {
	file, err := os.Open("users.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&users)
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

// Serve the game page
func serveGame(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("game.html")
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Users  []User
		Quests []Quest
	}{
		Users:  users,
		Quests: quests,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

// Create quest
func createQuest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		title := r.FormValue("title")
		userID := r.FormValue("user_id")

		xp, err := strconv.ParseInt(r.FormValue("xp"), 10, 64)
		if err != nil {
			log.Printf("Failed to parse XP: %v", err)
			http.Error(w, "Invalid XP value", http.StatusBadRequest)
			return
		}

		gold, err := strconv.ParseInt(r.FormValue("gold"), 10, 64)
		if err != nil {
			log.Printf("Failed to parse Gold: %v", err)
			http.Error(w, "Invalid Gold value", http.StatusBadRequest)
			return
		}

		questID := fmt.Sprintf("quest%d", len(quests)+1)
		newQuest := Quest{
			ID:       questID,
			Title:    title,
			UserID:   userID,
			XP:       xp,
			Gold:     gold,
			Complete: false,
		}

		quests = append(quests, newQuest)

		saveQuests()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Complete quest
func markComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		questID := r.FormValue("quest_id")
		quest := findQuestByID(questID)
		if quest == nil {
			http.Error(w, "Quest not found", http.StatusNotFound)
			return
		}

		user := findUserByID(quest.UserID)
		if user == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Mark quest as complete
		quest.Complete = true

		// Add quest XP and Gold to user's XP and Gold
		user.XP += quest.XP
		user.Gold += quest.Gold

		// Save changes
		saveQuests()
		saveUsers()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func findQuestByID(id string) *Quest {
	for i := range quests {
		if quests[i].ID == id {
			return &quests[i]
		}
	}
	return nil
}

func findUserByID(id string) *User {
	for i := range users {
		if users[i].ID == id {
			return &users[i]
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

func saveUsers() {
	userFile, err := os.Create("users.json")
	if err != nil {
		log.Fatal(err)
	}
	defer userFile.Close()

	encoder := json.NewEncoder(userFile)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(users)
	if err != nil {
		log.Fatal(err)
	}
}

// List quests (for testing)
func listQuests(w http.ResponseWriter, r *http.Request) {
	for _, quest := range quests {
		fmt.Fprintf(w, "Quest: %s, Assigned to: %s, Complete: %v\n", quest.Title, quest.UserID, quest.Complete)
	}
}
