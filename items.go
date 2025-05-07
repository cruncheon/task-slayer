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
	tmpl, err := template.ParseFiles("templates/base.html", "templates/items/list_items.html")
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
		tmpl, err := template.ParseFiles("templates/base.html", "templates/items/create_item.html")
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

		// If Post request, create item and redirect back to list items page
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

		// Set item ID by getting current amount of items and adding +1
		itemID := fmt.Sprintf("item%d", len(items)+1)

		// Structure new item details
		newItem := Item{
			ID:    itemID,
			Name:  name,
			Price: price,
		}

		// Add new item to items slice
		items = append(items, newItem)

		// Save changes to JSON file
		saveItems()

		log.Printf("%v - %v created", newItem.ID, newItem.Name)
		http.Redirect(w, r, "/items", http.StatusSeeOther)
	}
}

func editItem(w http.ResponseWriter, r *http.Request) {
	// Extract item ID from URL path
	id := r.URL.Path[len("/item/edit/"):]

	// Get the item by ID
	item := getItem(id)
	if item == nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		// Render edit item page
		tmpl, err := template.ParseFiles("templates/base.html", "templates/items/edit_item.html")
		if err != nil {
			log.Fatal(err)
		}

		data := struct {
			Item *Item
		}{
			Item: item,
		}

		err = tmpl.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			log.Fatal(err)
		}
	} else if r.Method == http.MethodPost {
		// Update the item details
		r.ParseForm()

		name := r.FormValue("name")

		// Convert price form input from string to int
		price, err := strconv.ParseInt(r.FormValue("price"), 10, 64)
		if err != nil {
			log.Printf("Failed to parse Price: %v", err)
			http.Error(w, "Invalid Price value", http.StatusBadRequest)
			return
		}

		// Update item details
		item.Name = name
		item.Price = price

		// Save changes to JSON file
		saveItems()

		log.Printf("%v - %v updated", item.ID, item.Name)
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
