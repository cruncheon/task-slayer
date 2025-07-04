package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/cruncheon/task-slayer/templates"
)

func LoadRoutes() {
	// Index
	http.HandleFunc("/", indexPage)

	// Admin
	http.HandleFunc("/admin", adminPage)

	// Static directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	// Quests
	http.HandleFunc("/quests", listQuests)
	http.HandleFunc("/quest/create", createQuest)
	http.HandleFunc("/quest/edit/", editQuest)
	http.HandleFunc("/quest/complete", completeQuest)
	http.HandleFunc("/quest/delete/", deleteQuest)

	// Players
	http.HandleFunc("/players", listPlayers)
	http.HandleFunc("/player/create", createPlayer)
	http.HandleFunc("/player/edit/", editPlayer)
	http.HandleFunc("/player/delete/", deletePlayer)

	// Items
	http.HandleFunc("/items", listItems)
	http.HandleFunc("/item/create", createItem)
	http.HandleFunc("/item/edit/", editItem)
	http.HandleFunc("/item/delete/", deleteItem)
	http.HandleFunc("/item/buy/", buyItem)
}

// Index page
func indexPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/quests", http.StatusSeeOther)
}

// Render template helper
func renderTemplate(w http.ResponseWriter, tmplPath string, data any) {
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

// Parse form integers helper
func parseFormInt(r *http.Request, key string) (int64, error) {
	value := r.FormValue(key)
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Printf("Failed to parse %s: %v", key, err)
		return 0, fmt.Errorf("invalid %s value: %w", key, err)
	}
	return intValue, nil
}

func adminPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, templates.AdminPage, nil)
}
