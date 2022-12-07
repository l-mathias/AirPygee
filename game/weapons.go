package game

type Weapon struct {
	Entity
	EquipableItemStats
	Equipped bool
	Location
	Rarity
}

type Sword struct {
	Weapon
}

type Bow struct {
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
func (w *Weapon) GetLocation() Location {
	return w.Location
}

func NewSword(p Pos) *Sword {
	rarity := randomizeRarity()
	stats := adaptStatsToRarity(rarity, &EquipableItemStats{
		MinDamage: 5,
		MaxDamage: 10,
		Armor:     0,
		Critical:  0,
	})
	return &Sword{
		Weapon: Weapon{Entity: Entity{
			Pos:         p,
			Name:        "Sword",
			Rune:        's',
			Type:        Weapons,
			Description: "A common sword...",
		},
			Location:           RightHand,
			Rarity:             rarity,
			EquipableItemStats: *stats,
		}}
}

func NewBow(p Pos) *Bow {
	rarity := randomizeRarity()
	stats := adaptStatsToRarity(rarity, &EquipableItemStats{
		MinDamage: 5,
		MaxDamage: 10,
		Armor:     0,
		Critical:  0,
	})
	return &Bow{
		Weapon: Weapon{Entity: Entity{
			Pos:         p,
			Name:        "Bow",
			Rune:        'B',
			Type:        Weapons,
			Description: "A common bow...",
		},
			Location:           RightHand,
			Rarity:             rarity,
			EquipableItemStats: *stats,
		}}
}
