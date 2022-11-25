package game

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

func NewHealthPotion(p Pos, size string) *Potion {
	return &Potion{
		Entity: Entity{
			Pos:         p,
			Name:        "Potion",
			Rune:        'p',
			Type:        Potions,
			Description: "A small health potion...",
		},
		Size: size,
	}
}

func (game *Game) consumePotion(item ConsumableItem) {
	switch item.GetSize() {
	case "Small":
		game.heal(&game.CurrentLevel.Player.Character, int(float64(game.CurrentLevel.Player.MaxHealth)*.25))
	case "Medium":
		game.heal(&game.CurrentLevel.Player.Character, int(float64(game.CurrentLevel.Player.MaxHealth)*.50))
	case "Large":
		game.heal(&game.CurrentLevel.Player.Character, int(float64(game.CurrentLevel.Player.MaxHealth)*.75))
	}
	game.removeInventoryItem(item, &game.CurrentLevel.Player.Character)
	game.CurrentLevel.AddEvent(game.CurrentLevel.Player.Character.Name + " consumed " + item.GetSize() + item.GetName())
	game.CurrentLevel.LastEvent = ConsumePotion
}
