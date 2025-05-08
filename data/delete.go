package data

func DeletePlayer(id string) error {
	var updatedPlayers []Player
	for _, player := range Players {
		if player.ID != id {
			updatedPlayers = append(updatedPlayers, player)
		}
	}

	Players = updatedPlayers

	// Save changes to JSON file
	SavePlayers()
	return nil
}

func DeleteQuest(id string) error {
	var updatedQuests []Quest
	for _, quest := range Quests {
		if quest.ID != id {
			updatedQuests = append(updatedQuests, quest)
		}
	}

	Quests = updatedQuests

	// Save changes to JSON file
	SaveQuests()
	return nil
}

func DeleteItem(id string) error {
	var updatedItems []Item
	for _, item := range Items {
		if item.ID != id {
			updatedItems = append(updatedItems, item)
		}
	}

	Items = updatedItems

	// Save changes to JSON file
	SaveItems()
	return nil
}
