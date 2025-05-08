package data

// Variables for data file paths
var (
	itemsFile   string = "data/items.json"
	playersFile string = "data/players.json"
	questsFile  string = "data/quests.json"
)

// Varialbes for data slices
var (
	Players []Player
	Quests  []Quest
	Items   []Item
)

// Player data struct
type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	XP   int64  `json:"xp"`
	Gold int64  `json:"gold"`
}

// Quest data struct
type Quest struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	PlayerID string `json:"player_id"`
	XP       int64  `json:"xp"`
	Gold     int64  `json:"gold"`
	Complete bool   `json:"complete"`
}

// Item data struct
type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

// Load data function
func LoadData() {
	loadItems()
	loadQuests()
	loadPlayers()
}
