package game

type Monster struct {
	Character
}

func NewRat(p Pos) *Monster {
	return &Monster{Character{
		Entity:       Entity{p, "Rat", 'R'},
		Hitpoints:    50,
		Strength:     1,
		Speed:        2.0,
		ActionPoints: 0.0,
	}}
}

func NewSpider(p Pos) *Monster {
	return &Monster{Character{
		Entity:       Entity{p, "Spider", 'S'},
		Hitpoints:    10,
		Strength:     2,
		Speed:        1.0,
		ActionPoints: 0.0,
	}}
}

func (m *Monster) Update(level *Level) {
	m.ActionPoints += m.Speed
	playerPos := level.Player.Pos
	apInt := int(m.ActionPoints)
	positions := level.astar(m.Pos, playerPos)

	if len(positions) == 0 {
		m.Pass()
		return
	}

	moveIndex := 1
	for i := 0; i < apInt; i++ {
		if moveIndex < len(positions) {
			m.Move(positions[moveIndex], level)

			moveIndex++
			m.ActionPoints--
		}
	}
}

func (m *Monster) Move(to Pos, level *Level) {
	_, exists := level.Monsters[to]
	if !exists && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
		return
	}

	if to == level.Player.Pos {
		level.Attack(&m.Character, &level.Player.Character)

		if m.Hitpoints <= 0 {
			delete(level.Monsters, m.Pos)
		}
		if level.Player.Hitpoints <= 0 {
			panic("You are dead")
		}
	}

}
