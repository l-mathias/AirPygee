package game

type Player struct {
	Character
}

func NewPlayer() *Player {
	player := &Player{Character: Character{
		Entity:        Entity{Name: "Wizard", Rune: '@'},
		Health:        20,
		MaxHealth:     20,
		MinDamage:     10,
		MaxDamage:     20,
		Armor:         0,
		Critical:      0,
		Speed:         1.0,
		ActionPoints:  0,
		SightRange:    10,
		InventorySize: 20,
	}}
	return player
}
