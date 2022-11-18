package game

type Armor struct {
	Entity
	EquipableItemStats
	Equipped bool
	Rarity
	Location
}

type Helmet struct {
	Armor
}

type Boots struct {
	Armor
}

type Plate struct {
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
func (a *Armor) GetLocation() Location {
	return a.Location
}

func NewPlate(p Pos) *Boots {
	rarity := randomizeRarity()
	stats := adaptStatsToRarity(rarity, &EquipableItemStats{
		MinDamage: 0,
		MaxDamage: 0,
		Armor:     10,
		Critical:  0,
	})
	return &Boots{Armor: Armor{
		Entity: Entity{
			Pos:         p,
			Name:        "Plate",
			Rune:        'a',
			Type:        Armors,
			Description: "Common plate...",
		},
		Location:           Chest,
		Rarity:             rarity,
		EquipableItemStats: *stats,
	}}
}

func NewBoots(p Pos) *Boots {
	rarity := randomizeRarity()
	stats := adaptStatsToRarity(rarity, &EquipableItemStats{
		MinDamage: 0,
		MaxDamage: 0,
		Armor:     5,
		Critical:  0,
	})
	return &Boots{Armor: Armor{
		Entity: Entity{
			Pos:         p,
			Name:        "Boots",
			Rune:        'b',
			Type:        Armors,
			Description: "Common boots...",
		},
		Location:           Foots,
		Rarity:             rarity,
		EquipableItemStats: *stats,
	}}
}

func NewHelmet(p Pos) *Helmet {
	rarity := randomizeRarity()
	stats := adaptStatsToRarity(rarity, &EquipableItemStats{
		MinDamage: 0,
		MaxDamage: 0,
		Armor:     5,
		Critical:  0,
	})
	return &Helmet{Armor: Armor{
		Entity: Entity{
			Pos:         p,
			Name:        "Helmet",
			Rune:        'h',
			Type:        Armors,
			Description: "A common helmet...",
		},
		Location:           Head,
		Rarity:             rarity,
		EquipableItemStats: *stats,
	}}
}
