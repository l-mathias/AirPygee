package game

import "fmt"

type Monster struct {
	Character
}

func NewRat(p Pos) *Monster {
	return &Monster{Character{
		Entity:       Entity{p, "Rat", 'R'},
		Hitpoints:    5,
		Strength:     5,
		Speed:        2.0,
		ActionPoints: 0.0,
	}}
}

func NewSpider(p Pos) *Monster {
	return &Monster{Character{
		Entity:       Entity{p, "Spider", 'S'},
		Hitpoints:    30,
		Strength:     10,
		Speed:        1.0,
		ActionPoints: 0.0,
	}}
}

func (m *Monster) Update(level *Level) {
	m.ActionPoints += m.Speed
	playerPos := level.Player.Pos
	apInt := int(m.ActionPoints)
	positions := level.astar(m.Pos, playerPos)
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
	} else {
		Attack(m, level.Player)
		if m.Hitpoints <= 0 {
			delete(level.Monsters, m.Pos)
			fmt.Printf("Monster %v is dead\n", m.Name)
		}
		if level.Player.Hitpoints <= 0 {
			fmt.Println("you are dead !")
			panic("You are dead")
		}
	}
}
