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

//type Item struct {
//	Entity
//	Location
//	Type        ItemType
//	Size        string
//	Equipped    bool
//	Description string
//}

type Item interface {
	GetDescription() string
	GetName() string
	GetRune() rune
	GetEntity() Entity
	SetPos(Pos)
}

type EquippableItem interface {
	Item
	IsEquipped() bool
	Equip()
	UnEquip()
}

type ConsumableItem interface {
	Item
	GetSize() string
}

type Weapon struct {
	Entity
	Strength int
	Equipped bool
}

type Sword struct {
	Weapon
}

func (w *Weapon) GetDescription() string {
	return w.Description
}
func (w *Weapon) GetName() string {
	return w.Name
}
func (w *Weapon) GetRune() rune {
	return w.Rune
}
func (w *Weapon) IsEquipped() bool {
	return w.Equipped
}
func (w *Weapon) Equip() {
	w.Equipped = true
}
func (w *Weapon) UnEquip() {
	w.Equipped = false
}
func (w *Weapon) GetEntity() Entity {
	return w.Entity
}
func (w *Weapon) SetPos(pos Pos) {
	w.Pos = pos
}

type Armor struct {
	Entity
	Defense  int
	Equipped bool
}

type Helmet struct {
	Armor
}

func (a *Armor) GetDescription() string {
	return a.Description
}
func (a *Armor) GetName() string {
	return a.Name
}
func (a *Armor) GetRune() rune {
	return a.Rune
}
func (a *Armor) IsEquipped() bool {
	return a.Equipped
}
func (a *Armor) Equip() {
	a.Equipped = true
}
func (a *Armor) UnEquip() {
	a.Equipped = false
}
func (a *Armor) GetEntity() Entity {
	return a.Entity
}
func (a *Armor) SetPos(pos Pos) {
	a.Pos = pos
}

type Potion struct {
	Entity
	Size string
}

func (p *Potion) GetDescription() string {
	return p.Description
}
func (p *Potion) GetName() string {
	return p.Name
}
func (p *Potion) GetRune() rune {
	return p.Rune
}
func (p *Potion) GetEntity() Entity {
	return p.Entity
}
func (p *Potion) SetPos(pos Pos) {
	p.Pos = pos
}
func (p *Potion) GetSize() string {
	return p.Size
}

func NewSword(p Pos) *Sword {
	return &Sword{
		Weapon: Weapon{Entity: Entity{
			Pos:         p,
			Name:        "Sword",
			Rune:        's',
			Type:        Weapons,
			Description: "A common sword...",
			Location:    RightHand,
		}, Strength: 5},
	}
}

func NewHelmet(p Pos) *Helmet {
	return &Helmet{Armor: Armor{
		Entity: Entity{
			Pos:         p,
			Name:        "Helmet",
			Rune:        'h',
			Type:        Armors,
			Description: "A common Helmet...",
			Location:    Head,
		},
		Defense: 5,
	}}
}

func NewHealthPotion(p Pos, size string) *Potion {
	return &Potion{
		Entity: Entity{
			Pos:         p,
			Name:        "Potion",
			Rune:        'p',
			Location:    NoLoc,
			Type:        Potions,
			Description: "A small health potion...",
		},
		Size: "Small",
	}
}

//func NewSword(p Pos) *Item {
//	return &Item{Entity: Entity{
//		Pos:  p,
//		Name: "Sword",
//		Rune: 's',
//	}, Location: RightHand, Equipped: false, Type: Weapons, Description: "A common sword..."}
//}
//
//func NewHelmet(p Pos) *Item {
//	return &Item{Entity: Entity{
//		Pos:  p,
//		Name: "Helmet",
//		Rune: 'h',
//	}, Location: Head, Equipped: false, Type: Armors, Description: "A basic helmet..."}
//}

//func NewHealthPotion(p Pos, size string) *Item {
//	return &Item{Entity: Entity{
//		Pos:  p,
//		Name: "Potion",
//		Rune: 'p',
//	}, Location: NoLoc, Equipped: false, Type: Potions, Size: size, Description: "A " + size + " health potion"}
//}

func (game *Game) consumePotion(item ConsumableItem) {
	switch item.GetSize() {
	case "Small":
		game.heal(int(float64(game.CurrentLevel.Player.MaxHealth) * .25))
	case "Medium":
		game.heal(int(float64(game.CurrentLevel.Player.MaxHealth) * .50))
	case "Large":
		game.heal(int(float64(game.CurrentLevel.Player.MaxHealth) * .75))
	}
	game.removeInventoryItem(item, &game.CurrentLevel.Player.Character)
	game.CurrentLevel.AddEvent(game.CurrentLevel.Player.Character.Name + " consumed " + item.GetSize() + item.GetName())
	game.CurrentLevel.LastEvent = ConsumePotion
}

func (game *Game) equip(itemToEquip EquippableItem) {

	if game.slotFreeToEquip(itemToEquip) {
		itemToEquip.Equip()
		game.CurrentLevel.Player.EquippedItems = append(game.CurrentLevel.Player.EquippedItems, itemToEquip)
		for i, item := range game.CurrentLevel.Player.Items {
			if item == itemToEquip {
				game.CurrentLevel.Player.Items = append(game.CurrentLevel.Player.Items[:i], game.CurrentLevel.Player.Items[i+1:]...)
			}
		}
	}
}

func (game *Game) unEquip(itemToUnEquip EquippableItem) {
	itemToUnEquip.UnEquip()
	game.CurrentLevel.Player.Items = append(game.CurrentLevel.Player.Items, itemToUnEquip)
	for i, item := range game.CurrentLevel.Player.EquippedItems {
		if item == itemToUnEquip {
			game.CurrentLevel.Player.EquippedItems = append(game.CurrentLevel.Player.EquippedItems[:i], game.CurrentLevel.Player.EquippedItems[i+1:]...)
		}
	}
}

func (game *Game) slotFreeToEquip(itemToCheck Item) bool {
	for _, item := range game.CurrentLevel.Player.EquippedItems {
		if item.GetEntity().Location == itemToCheck.GetEntity().Location {
			return false
		}
	}
	return true
}
