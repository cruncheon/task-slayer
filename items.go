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

func listItems(w http.ResponseWriter, r *http.Request) {
	// Parse the template files
	tmpl, err := template.ParseFiles("templates/base.html", "templates/list_items.html")
	if err != nil {
		log.Fatal(err)
	}

	// Define data structure for list of items
	data := struct {
		Items []Item
	}{
		Items: items,
	}

	// Execute the template with the data structure
	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Fatal(err)
	}
}

// Create Item
func createItem(w http.ResponseWriter, r *http.Request) {
	// If Get request, render create item page
	if r.Method == http.MethodGet {
		// Parse the template files
		tmpl, err := template.ParseFiles("templates/base.html", "templates/create_item.html")
		if err != nil {
		}

		// Define data structure for item creation form
		data := struct {
			Items []Item
		}{
			Items: items,
		}

		// Execute template and render create item page
		err = tmpl.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			log.Fatal(err)
		}

		// If Post request, create player and redirect back to list players page
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		name := r.FormValue("name")

		// Convert price form input from string to int
		price, err := strconv.ParseInt(r.FormValue("price"), 10, 64)
		if err != nil {
			log.Printf("Failed to parse Price: %v", err)
			http.Error(w, "Invalid Price value", http.StatusBadRequest)
			return
		}

		// Set player ID by getting current amount of quests and adding +1
		itemID := fmt.Sprintf("item%d", len(items)+1)

		// Structure new player details
		newItem := Item{
			ID:    itemID,
			Name:  name,
			Price: price,
		}

		// Add new player to players slice
		items = append(items, newItem)

		// Save changes to JSON file
		saveItems()

		log.Printf("%v - %v created", newItem.ID, newItem.Name)
		http.Redirect(w, r, "/items", http.StatusSeeOther)
	}
}

// Load items
func loadItems() {
	filePath := "data/items.json"

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the file if it does not exist
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Initialize an empty slice of items and write to the file
		var initialItems []Item
		err = json.NewEncoder(file).Encode(initialItems)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// Open the file if it exists
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decode JSON from quest file
	err = json.NewDecoder(file).Decode(&items)
	if err != nil {
		log.Fatal(err)
	}
}

// Save items
func saveItems() {
	itemFile, err := os.Create("data/items.json")
	if err != nil {
		log.Fatal(err)
	}
	defer itemFile.Close()

	encoder := json.NewEncoder(itemFile)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(items)
	if err != nil {
		log.Fatal(err)
	}
}

// Get item details by ID
func getItem(id string) *Item {
	for i := range items {
		if items[i].ID == id {
			return &items[i]
		}
	}
	return nil
}
