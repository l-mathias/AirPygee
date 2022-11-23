package game

import (
	"crypto/rand"
	"math/big"
)

type Location int

const (
	Foots Location = iota
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
	TreasureChests
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
	GetLocation() Location
}

type OpenableItem interface {
	Item
	GetSize() int
	GetItems() []Item
	RemoveItems()
	GetPos() Pos
	GetState() bool
	Open()
	Close()
}

type ConsumableItem interface {
	Item
	GetSize() string
}

type EquipableItemStats struct {
	MinDamage int
	MaxDamage int
	Armor     int
	Critical  float64
}

func adaptStatsToRarity(rarity Rarity, stats *EquipableItemStats) *EquipableItemStats {
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

	stats.Critical *= multiplier
	stats.MinDamage = int(float64(stats.MinDamage) * multiplier)
	stats.MaxDamage = int(float64(stats.MaxDamage) * multiplier)
	stats.Armor = int(float64(stats.Armor) * multiplier)
	return stats
}

func randomizeRarity() Rarity {
	number, err := rand.Int(rand.Reader, big.NewInt(100))
	CheckError(err)

	switch {
	case number.Int64() <= 2:
		return Legendary
	case number.Int64() > 2 && number.Int64() <= 10:
		return Epic
	case number.Int64() > 10 && number.Int64() <= 20:
		return Rare
	case number.Int64() > 20 && number.Int64() <= 40:
		return Uncommon
	case number.Int64() > 40 && number.Int64() <= 100:
		return Common
	}

	return Common
}

func (game *Game) adaptPlayerStats(item EquipableItem, addOrRemove string) {
	if addOrRemove == "add" {
		game.CurrentLevel.Player.MinDamage += item.GetStats().MinDamage
		game.CurrentLevel.Player.MaxDamage += item.GetStats().MaxDamage
		game.CurrentLevel.Player.Critical += item.GetStats().Critical
		game.CurrentLevel.Player.Armor += item.GetStats().Armor
	} else if addOrRemove == "remove" {
		game.CurrentLevel.Player.MinDamage -= item.GetStats().MinDamage
		game.CurrentLevel.Player.MaxDamage -= item.GetStats().MaxDamage
		game.CurrentLevel.Player.Critical -= item.GetStats().Critical
		game.CurrentLevel.Player.Armor -= item.GetStats().Armor
	}
}

func (game *Game) equip(itemToEquip EquipableItem) {
	if game.slotFreeToEquip(itemToEquip) {
		itemToEquip.Equip()
		game.adaptPlayerStats(itemToEquip, "add")
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
	game.adaptPlayerStats(itemToUnEquip, "remove")
	game.CurrentLevel.Player.Items = append(game.CurrentLevel.Player.Items, itemToUnEquip)
	for i, item := range game.CurrentLevel.Player.EquippedItems {
		if item == itemToUnEquip {
			game.CurrentLevel.Player.EquippedItems = append(game.CurrentLevel.Player.EquippedItems[:i], game.CurrentLevel.Player.EquippedItems[i+1:]...)
		}
	}
}

func (game *Game) slotFreeToEquip(itemToCheck EquipableItem) bool {
	for _, item := range game.CurrentLevel.Player.EquippedItems {
		if item.(EquipableItem).GetLocation() == itemToCheck.GetLocation() {
			return false
		}
	}
	return true
}

func randomLoot(p Pos, numItems int) []Item {
	items := make([]Item, 0)

	for i := 0; i < numItems; i++ {
		number, err := rand.Int(rand.Reader, big.NewInt(4))
		CheckError(err)

		switch {
		case number.Int64() == 0:
			items = append(items, NewHelmet(p))
		case number.Int64() == 1:
			items = append(items, NewSword(p))
		case number.Int64() == 2:
			items = append(items, NewPlate(p))
		case number.Int64() == 3:
			items = append(items, NewHealthPotion(p, "Small"))
		case number.Int64() == 4:
			items = append(items, NewBoots(p))
		}
	}

	return items
}
