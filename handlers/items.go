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
	switch r.Method {
	case "GET":
		pageData := struct {
			Items []data.Item
		}{
			data.Items,
		}

		renderTemplate(w, templates.CreateItem, pageData)

	case "POST":
		r.ParseForm()

		name := r.FormValue("name")

		price, err := parseFormInt(r, "price")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		itemID := fmt.Sprintf("item%d", len(data.Items)+1)

		newItem := data.Item{
			ID:    itemID,
			Name:  name,
			Price: price,
		}

		data.Items = append(data.Items, newItem)

		data.SaveItems()

		log.Printf("%v - %v created", newItem.ID, newItem.Name)
		http.Redirect(w, r, "/items", http.StatusSeeOther)

	default:
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
	}
}

func editItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/item/edit/"):]

	item := data.GetItem(id)
	if item == nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		pageData := struct {
			Item *data.Item
		}{
			item,
		}

		renderTemplate(w, templates.EditItem, pageData)

	case "POST":
		r.ParseForm()

		name := r.FormValue("name")

		price, err := parseFormInt(r, "price")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		item.Name = name
		item.Price = price

		data.SaveItems()

		log.Printf("%v - %v updated", item.ID, item.Name)
		http.Redirect(w, r, "/items", http.StatusSeeOther)

	default:
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
	}
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/item/delete/"):]
	if err := data.DeleteItem(id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete item: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("%v deleted", id)
	http.Redirect(w, r, "/items", http.StatusSeeOther)
}

func buyItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/item/buy/"):]
	pid := r.FormValue("player_id")

	item := data.GetItem(id)
	if item == nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	player := data.GetPlayer(pid)
	if player == nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	if player.Gold < item.Price {
		http.Error(w, "Not enough gold", http.StatusBadRequest)
		return
	}

	player.Gold -= item.Price

	player.Items = append(player.Items, item.Name)

	data.SavePlayers()

	log.Printf("%v bought %v", player.Name, item.Name)
	http.Redirect(w, r, "/items", http.StatusSeeOther)
}
