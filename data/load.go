package data

import (
	"encoding/json"
	"log"
	"os"
)

// Load players
func loadPlayers() {
	filePath := playersFile

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the file if it does not exist
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Initialize an empty slice of players and write to the file
		var initialPlayers []Player
		err = json.NewEncoder(file).Encode(initialPlayers)
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

	// Decode JSON from player file
	err = json.NewDecoder(file).Decode(&Players)
	if err != nil {
		log.Fatal(err)
	}
}

// Load items
func loadItems() {
	filePath := itemsFile

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
	err = json.NewDecoder(file).Decode(&Items)
	if err != nil {
		log.Fatal(err)
	}
}

// Load quests
func loadQuests() {
	filePath := questsFile

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the file if it does not exist
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Initialize an empty slice of quests and write to the file
		var initialQuests []Quest
		err = json.NewEncoder(file).Encode(initialQuests)
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
	err = json.NewDecoder(file).Decode(&Quests)
	if err != nil {
		log.Fatal(err)
	}
}
