package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cruncheon/task-slayer/data"
	"github.com/cruncheon/task-slayer/templates"
)

func listItems(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Items   []data.Item
		Players []data.Player
	}{
		data.Items,
		data.Players,
	}

	renderTemplate(w, templates.ListItems, data)
}

// Create Item
func createItem(w http.ResponseWriter, r *http.Request) {
	// If Get request, render create item page
	if r.Method == http.MethodGet {
		pageData := struct {
			Items []data.Item
		}{
			data.Items,
		}

		renderTemplate(w, templates.CreateItem, pageData)

		// If Post request, create item and redirect back to list items page
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		name := r.FormValue("name")

		// Convert price form input from string to int
		price, err := parseFormInt(r, "price")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
		pageData := struct {
			Item *data.Item
		}{
			item,
		}

		renderTemplate(w, templates.EditItem, pageData)

	} else if r.Method == http.MethodPost {
		// Update the item details
		r.ParseForm()

		name := r.FormValue("name")

		// Convert price form input from string to int
		price, err := parseFormInt(r, "price")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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

func deleteItem(w http.ResponseWriter, r *http.Request) {
	// Extract item ID from URL path
	id := r.URL.Path[len("/item/delete/"):]
	if err := data.DeleteItem(id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete item: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("%v deleted", id)
	http.Redirect(w, r, "/items", http.StatusSeeOther)
}

func buyItem(w http.ResponseWriter, r *http.Request) {
	// Extract item ID an player ID
	id := r.URL.Path[len("/item/buy/"):]
	pid := r.FormValue("player_id")

	// Get the item by ID
	item := data.GetItem(id)
	if item == nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// Get the player by ID
	player := data.GetPlayer(pid)
	if player == nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	// Check if player has enough gold to buy item
	if player.Gold < item.Price {
		http.Error(w, "Not enough gold", http.StatusBadRequest)
		return
	}

	// Deduct gold from player
	player.Gold -= item.Price

	// Add item to player's inventory
	player.Items = append(player.Items, item.Name)

	// Save changes
	data.SavePlayers()

	log.Printf("%v bought %v", player.Name, item.Name)
	http.Redirect(w, r, "/items", http.StatusSeeOther)
}
