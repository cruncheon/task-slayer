package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/cruncheon/task-slayer/templates"
)

var (
	itemsFile   string = "data/items.json"
	playersFile string = "data/players.json"
	questsFile  string = "data/quests.json"
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

type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

var (
	players []Player
	quests  []Quest
	items   []Item
)

func LoadData() {
	loadPlayers()
	loadQuests()
	loadItems()
}

func LoadRoutes() {
	http.HandleFunc("/", homePage)

	http.HandleFunc("/quests", listQuests)
	http.HandleFunc("/quest/create", createQuest)
	http.HandleFunc("/quest/edit/", editQuest)
	http.HandleFunc("/quest/complete", completeQuest)

	http.HandleFunc("/players", listPlayers)
	http.HandleFunc("/player/create", createPlayer)
	http.HandleFunc("/player/edit/", editPlayer)

	http.HandleFunc("/items", listItems)
	http.HandleFunc("/item/create", createItem)
	http.HandleFunc("/item/edit/", editItem)
}

// Home page
func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(templates.Base, "templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, "base.html")
	if err != nil {
		log.Fatal(err)
	}
}
