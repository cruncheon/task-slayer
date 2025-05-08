package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/cruncheon/task-slayer/data"
	"github.com/cruncheon/task-slayer/templates"
)

func listItems(w http.ResponseWriter, r *http.Request) {
	// Parse the template files
	tmpl, err := template.ParseFiles(templates.Base, templates.ListItems)
	if err != nil {
		log.Fatal(err)
	}

	// Define data structure for list of items
	data := struct {
		Items []data.Item
	}{
		Items: data.Items,
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
		tmpl, err := template.ParseFiles(templates.Base, templates.CreateItem)
		if err != nil {
		}

		// Define data structure for item creation form
		data := struct {
			Items []data.Item
		}{
			Items: data.Items,
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
		itemID := fmt.Sprintf("item%d", len(data.Items)+1)

		// Structure new item details
		newItem := data.Item{
			ID:    itemID,
			Name:  name,
			Price: price,
		}

		// Add new item to items slice
		data.Items = append(data.Items, newItem)

		// Save changes to JSON file
		data.SaveItems()

		log.Printf("%v - %v created", newItem.ID, newItem.Name)
		http.Redirect(w, r, "/items", http.StatusSeeOther)
	}
}

func editItem(w http.ResponseWriter, r *http.Request) {
	// Extract item ID from URL path
	id := r.URL.Path[len("/item/edit/"):]

	// Get the item by ID
	item := data.GetItem(id)
	if item == nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		// Render edit item page
		tmpl, err := template.ParseFiles(templates.Base, templates.EditItem)
		if err != nil {
			log.Fatal(err)
		}

		data := struct {
			Item *data.Item
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
		data.SaveItems()

		log.Printf("%v - %v updated", item.ID, item.Name)
		http.Redirect(w, r, "/items", http.StatusSeeOther)
	}
}
