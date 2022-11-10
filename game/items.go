package game

import (
	"math/rand"
	"time"
)

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
type Rarity int

const (
	Armors ItemType = iota
	Weapons
	Potions
)

const (
	Common Rarity = iota
	Uncommon
	Rare
	Epic
	Legendary
)

type Item interface {
	GetDescription() string
	GetName() string
	GetRune() rune
	GetEntity() *Entity
	SetPos(Pos)
}

type EquipableItem interface {
	Item
	IsEquipped() bool
	Equip()
	UnEquip()
	GetStats() *EquipableItemStats
	GetRarity() Rarity
	ToString(Rarity) string
}

type ConsumableItem interface {
	Item
	GetSize() string
}

type EquipableItemStats struct {
	MinStrength int
	MaxStrength int
	MinDefense  int
	MaxDefense  int
	MinCritical float64
	MaxCritical float64
}

type Weapon struct {
	Entity
	EquipableItemStats
	Equipped bool
	Rarity
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
func (w *Weapon) GetEntity() *Entity {
	return &w.Entity
}
func (w *Weapon) SetPos(pos Pos) {
	w.Pos = pos
}
func (w *Weapon) GetRarity() Rarity {
	return w.Rarity
}
func (w *Weapon) GetStats() *EquipableItemStats {
	return &w.EquipableItemStats
}
func (w *Weapon) ToString(rarity Rarity) string {
	switch rarity {
	case Common:
		return "Common"
	case Uncommon:
		return "Uncommon"
	case Rare:
		return "Rare"
	case Epic:
		return "Epic"
	case Legendary:
		return "Legendary"
	}
	return ""
}

type Armor struct {
	Entity
	EquipableItemStats
	Equipped bool
	Rarity
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
func (a *Armor) GetEntity() *Entity {
	return &a.Entity
}
func (a *Armor) SetPos(pos Pos) {
	a.Pos = pos
}
func (a *Armor) GetRarity() Rarity {
	return a.Rarity
}
func (a *Armor) GetStats() *EquipableItemStats {
	return &a.EquipableItemStats
}
func (a *Armor) ToString(rarity Rarity) string {
	switch rarity {
	case Common:
		return "Common"
	case Uncommon:
		return "Uncommon"
	case Rare:
		return "Rare"
	case Epic:
		return "Epic"
	case Legendary:
		return "Legendary"
	}
	return ""
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
func (p *Potion) GetEntity() *Entity {
	return &p.Entity
}
func (p *Potion) SetPos(pos Pos) {
	p.Pos = pos
}
func (p *Potion) GetSize() string {
	return p.Size
}

func randomizeStats(rarity Rarity, stats *EquipableItemStats) *EquipableItemStats {
	var multiplier float64

	switch rarity {
	case Common:
		multiplier = 1
	case Uncommon:
		multiplier = 1.5
	case Rare:
		multiplier = 1.75
	case Epic:
		multiplier = 2
	case Legendary:
		multiplier = 3
	}

	stats.MinCritical *= multiplier
	stats.MaxCritical *= multiplier
	stats.MinStrength = int(float64(stats.MinStrength) * multiplier)
	stats.MaxStrength = int(float64(stats.MaxStrength) * multiplier)
	stats.MinDefense = int(float64(stats.MinDefense) * multiplier)
	stats.MaxDefense = int(float64(stats.MaxDefense) * multiplier)
	return stats
}

func randomizeRarity() Rarity {
	rand.Seed(time.Now().UnixNano())
	number := rand.Intn(100)

	switch {
	case number <= 2:
		return Legendary
	case number > 2 && number <= 10:
		return Epic
	case number > 10 && number <= 20:
		return Rare
	case number > 20 && number <= 40:
		return Uncommon
	case number > 40 && number <= 100:
		return Common
	}

	return Common
}

func NewSword(p Pos) *Sword {
	rarity := randomizeRarity()
	stats := randomizeStats(rarity, &EquipableItemStats{
		MinStrength: 5,
		MaxStrength: 10,
		MinDefense:  0,
		MaxDefense:  0,
		MinCritical: 0,
		MaxCritical: 0,
	})
	return &Sword{
		Weapon: Weapon{Entity: Entity{
			Pos:         p,
			Name:        "Sword",
			Rune:        's',
			Type:        Weapons,
			Description: "A common sword...",
			Location:    RightHand,
		},
			Rarity:             rarity,
			EquipableItemStats: *stats,
		}}
}

func NewHelmet(p Pos) *Helmet {
	rarity := randomizeRarity()
	stats := randomizeStats(rarity, &EquipableItemStats{
		MinStrength: 0,
		MaxStrength: 0,
		MinDefense:  5,
		MaxDefense:  10,
		MinCritical: 0,
		MaxCritical: 0,
	})
	return &Helmet{Armor: Armor{
		Entity: Entity{
			Pos:         p,
			Name:        "Helmet",
			Rune:        'h',
			Type:        Armors,
			Description: "A common helmet...",
			Location:    Head,
		},
		Rarity:             rarity,
		EquipableItemStats: *stats,
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
		Size: size,
	}
}

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

func (game *Game) equip(itemToEquip EquipableItem) {

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

func (game *Game) unEquip(itemToUnEquip EquipableItem) {
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
