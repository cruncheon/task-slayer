package data

// Get item details by ID
func GetItem(id string) *Item {
	for i := range Items {
		if Items[i].ID == id {
			return &Items[i]
		}
	}
	return nil
}

// Get player details by ID
func GetPlayer(id string) *Player {
	for i := range Players {
		if Players[i].ID == id {
			return &Players[i]
		}
	}
	return nil
}

// Get quest details by ID
func GetQuest(id string) *Quest {
	for i := range Quests {
		if Quests[i].ID == id {
			return &Quests[i]
		}
	}
	return nil
}
