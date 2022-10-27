package game

type Location int

const (
	Foots Location = iota
	LeftHand
	RightHand
	Head
	Chest
	Legs
)

type Item struct {
	Entity
	Location
	Equipped bool
}

func NewSword(p Pos) *Item {
	return &Item{Entity{
		Pos:  p,
		Name: "Sword",
		Rune: 's',
	}, RightHand, false}
}

func NewHelmet(p Pos) *Item {
	return &Item{Entity{
		Pos:  p,
		Name: "Helmet",
		Rune: 'h',
	}, Head, false}
}
