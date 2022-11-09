package game

type Monster struct {
	Character
}

func NewRat(p Pos) *Monster {
	return &Monster{Character{
		Entity:       Entity{Pos: p, Name: "Rat", Rune: 'R'},
		Health:       5,
		MaxHealth:    5,
		Strength:     1,
		Speed:        2.0,
		ActionPoints: 0.0,
	}}
}

func NewSpider(p Pos) *Monster {
	return &Monster{Character{
		Entity:       Entity{Pos: p, Name: "Spider", Rune: 'S'},
		Health:       10,
		MaxHealth:    10,
		Strength:     2,
		Speed:        1.0,
		ActionPoints: 0.0,
	}}
}

func (m *Monster) kill(level *Level) {
	delete(level.Monsters, m.Pos)
	groundItems := level.Items[m.Pos]
	for _, item := range m.Items {
		item.SetPos(m.Pos)
		groundItems = append(groundItems, item)
	}
	level.Items[m.Pos] = groundItems
}

func (m *Monster) Update(game *Game) {
	m.ActionPoints += m.Speed
	playerPos := game.CurrentLevel.Player.Pos
	apInt := int(m.ActionPoints)
	positions := game.CurrentLevel.astar(m.Pos, playerPos)

	if len(positions) == 0 {
		m.Pass()
		return
	}

	moveIndex := 1
	for i := 0; i < apInt; i++ {
		if moveIndex < len(positions) {
			m.Move(positions[moveIndex], game)

			moveIndex++
			m.ActionPoints--
		}
	}
}

func (m *Monster) Move(to Pos, game *Game) {
	_, exists := game.CurrentLevel.Monsters[to]
	if !exists && to != game.CurrentLevel.Player.Pos {
		delete(game.CurrentLevel.Monsters, m.Pos)
		game.CurrentLevel.Monsters[to] = m
		m.Pos = to
		return
	}

	if to == game.CurrentLevel.Player.Pos {
		game.CurrentLevel.Attack(&m.Character, &game.CurrentLevel.Player.Character)

		if m.Health <= 0 {
			delete(game.CurrentLevel.Monsters, m.Pos)
		}
	}

}
