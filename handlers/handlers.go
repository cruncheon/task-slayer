package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/cruncheon/task-slayer/templates"
)

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
