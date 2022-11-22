package game

type TreasureChest struct {
	Entity
	Size   int
	Items  []Item
	Opened bool
}

func (t *TreasureChest) GetDescription() string {
	return t.Description
}
func (t *TreasureChest) GetName() string {
	return t.Name
}
func (t *TreasureChest) GetRune() rune {
	return t.Rune
}
func (t *TreasureChest) GetEntity() *Entity {
	return &t.Entity
}
func (t *TreasureChest) SetPos(pos Pos) {
	t.Pos = pos
}
func (t *TreasureChest) GetSize() int {
	return t.Size
}
func (t *TreasureChest) GetItems() *[]Item {
	return &t.Items
}
func (t *TreasureChest) RemoveItem(item *Item) {
	for i, v := range t.Items {
		if &v == item {
			t.Items = append(t.Items[:i], t.Items[i+1:]...)
			break
		}
	}
}
func (t *TreasureChest) GetPos() Pos {
	return t.Pos
}
func (t *TreasureChest) GetState() bool {
	return t.Opened
}
func (t *TreasureChest) Open() {
	t.Opened = true
}
func (t *TreasureChest) Close() {
	t.Opened = false
}

func NewTreasureChest(p Pos, size int) *TreasureChest {
	items := randomLoot(p, size)
	return &TreasureChest{
		Entity: Entity{
			Pos:         p,
			Name:        "Treasure Chest",
			Rune:        't',
			Type:        TreasureChests,
			Description: "A treasure chest...",
		},
		Items: *items,
		Size:  size,
	}
}

func (game *Game) OpenItem(chest OpenableItem) {
	items := chest.GetItems()
	for _, item := range *items {
		game.CurrentLevel.Items[game.CurrentLevel.Player.Pos] = append(game.CurrentLevel.Items[game.CurrentLevel.Player.Pos], item)
		chest.RemoveItem(&item)
	}
	chest.Open()
	game.CurrentLevel.AddEvent(game.CurrentLevel.Player.Name + " Opened chest")
	game.CurrentLevel.LastEvent = OpenChest
	game.CurrentLevel.Map[chest.GetPos().Y][chest.GetPos().X].Actionable = false
}
