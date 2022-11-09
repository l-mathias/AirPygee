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
	Potions
)

type Item struct {
	Entity
	Location
	Type        ItemType
	Size        string
	Equipped    bool
	Description string
}

func NewSword(p Pos) *Item {
	return &Item{Entity: Entity{
		Pos:  p,
		Name: "Sword",
		Rune: 's',
	}, Location: RightHand, Equipped: false, Type: Weapons, Description: "A common sword..."}
}

func NewHelmet(p Pos) *Item {
	return &Item{Entity: Entity{
		Pos:  p,
		Name: "Helmet",
		Rune: 'h',
	}, Location: Head, Equipped: false, Type: Armors, Description: "A basic helmet..."}
}

func NewHealthPotion(p Pos, size string) *Item {
	return &Item{Entity: Entity{
		Pos:  p,
		Name: "Potion",
		Rune: 'p',
	}, Location: NoLoc, Equipped: false, Type: Potions, Size: size, Description: "A " + size + " health potion"}
}

func (game *Game) consumePotion(item *Item) {
	switch item.Size {
	case "Small":
		game.heal(int(float64(game.CurrentLevel.Player.MaxHealth) * .25))
	case "Medium":
		game.heal(int(float64(game.CurrentLevel.Player.MaxHealth) * .50))
	case "Large":
		game.heal(int(float64(game.CurrentLevel.Player.MaxHealth) * .75))
	}
	game.removeInventoryItem(interface{}(item).(*Item), &game.CurrentLevel.Player.Character)
	game.CurrentLevel.AddEvent(game.CurrentLevel.Player.Character.Name + " consumed " + item.Size + item.Name)
	game.CurrentLevel.LastEvent = ConsumePotion
}

func (game *Game) equip(itemToEquip *Item) {
	if game.slotFreeToEquip(itemToEquip) {
		itemToEquip.Equipped = true

		game.CurrentLevel.Player.EquippedItems = append(game.CurrentLevel.Player.EquippedItems, itemToEquip)
		for i, item := range game.CurrentLevel.Player.Items {
			if item == interface{}(itemToEquip).(*Item) {
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
