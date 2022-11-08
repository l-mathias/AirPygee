package game

type Player struct {
	Character
}

func NewPlayer() *Player {
	player := &Player{Character: Character{
		Entity:        Entity{Name: "Wizard", Rune: '@'},
		Health:        20,
		MaxHealth:     20,
		Strength:      20,
		Speed:         1.0,
		ActionPoints:  0,
		SightRange:    10,
		InventorySize: 20,
	}}
	return player
}
