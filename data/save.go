package data

import (
	"encoding/json"
	"log"
	"os"
)

// Save players
func SavePlayers() {
	file, err := os.Create(playersFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(Players)
	if err != nil {
		log.Fatal(err)
	}
}

// Save quests
func SaveQuests() {
	file, err := os.Create(questsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(Quests)
	if err != nil {
		log.Fatal(err)
	}
}

// Save items
func SaveItems() {
	file, err := os.Create(itemsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(Items)
	if err != nil {
		log.Fatal(err)
	}
}
