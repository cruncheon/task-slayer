package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/cruncheon/task-slayer/templates"
)

func LoadRoutes() {
	http.HandleFunc("/", indexPage)

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

// Index page
func indexPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, templates.Index, nil)
}

// Render template helper
func renderTemplate(w http.ResponseWriter, tmplPath string, data interface{}) {
	tmpl, err := template.ParseFiles(templates.Base, tmplPath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error parsing templates:", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
	}
}
