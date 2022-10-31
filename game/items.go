package game

type Location int

const (
	NoLoc Location = iota
	Foots
	LeftHand
	RightHand
	Head
	Chest
	Legs
)

type ItemType int

const (
	Armors ItemType = iota
	Weapons
	Potion
)

type Item struct {
	Entity
	Location
	Type     ItemType
	Equipped bool
	Size     rune
}

func NewSword(p Pos) *Item {
	return &Item{Entity: Entity{
		Pos:  p,
		Name: "Sword",
		Rune: 's',
	}, Location: RightHand, Equipped: false, Type: Weapons}
}

func NewHelmet(p Pos) *Item {
	return &Item{Entity: Entity{
		Pos:  p,
		Name: "Helmet",
		Rune: 'h',
	}, Location: Head, Equipped: false, Type: Armors}
}

func NewHealthPotion(p Pos) *Item {
	return &Item{Entity: Entity{
		Pos:  p,
		Name: "Potion",
		Rune: 'p',
	}, Location: NoLoc, Equipped: false, Type: Potion, Size: 's'}
}

func (game *Game) consumePotion(item *Item) {
	switch item.Size {
	case 's':
		game.CurrentLevel.Player.Hitpoints += int(float64(game.CurrentLevel.Player.MaxHitpoints) * .25)
		game.removeInventoryItem(item, &game.CurrentLevel.Player.Character)
		game.CurrentLevel.AddEvent(game.CurrentLevel.Player.Character.Name + " consumed " + string(item.Size) + " " + item.Name)
		game.CurrentLevel.LastEvent = ConsumePotion
	case 'm':
	case 'l':
	}
}

func (game *Game) equip(itemToEquip *Item) {
	if game.slotFreeToEquip(itemToEquip) {
		itemToEquip.Equipped = true
		game.CurrentLevel.Player.EquippedItems = append(game.CurrentLevel.Player.EquippedItems, itemToEquip)
		for i, item := range game.CurrentLevel.Player.Items {
			if item == itemToEquip {
				game.CurrentLevel.Player.Items = append(game.CurrentLevel.Player.Items[:i], game.CurrentLevel.Player.Items[i+1:]...)
			}
		}
	}
}

func (game *Game) unEquip(itemToUnEquip *Item) {
	itemToUnEquip.Equipped = false
	game.CurrentLevel.Player.Items = append(game.CurrentLevel.Player.Items, itemToUnEquip)
	for i, item := range game.CurrentLevel.Player.EquippedItems {
		if item == itemToUnEquip {
			game.CurrentLevel.Player.EquippedItems = append(game.CurrentLevel.Player.EquippedItems[:i], game.CurrentLevel.Player.EquippedItems[i+1:]...)
		}
	}
}

func (game *Game) slotFreeToEquip(itemToCheck *Item) bool {
	for _, item := range game.CurrentLevel.Player.EquippedItems {
		if item.Location == itemToCheck.Location {
			return false
		}
	}
	return true
}
