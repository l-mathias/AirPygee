package game

import (
	"crypto/rand"
	"math/big"
)

type Monster struct {
	Character
}

func randomizeLoot(p Pos) []Item {
	numItems := 0

	number, err := rand.Int(rand.Reader, big.NewInt(100))
	CheckError(err)

	switch {
	case number.Int64() <= 1:
		numItems = 4
	case number.Int64() > 1 && number.Int64() <= 5:
		numItems = 3
	case number.Int64() > 5 && number.Int64() <= 15:
		numItems = 2
	case number.Int64() > 15 && number.Int64() <= 35:
		numItems = 1
	case number.Int64() > 35 && number.Int64() <= 100:
		numItems = 0
	}

	return randomLoot(p, numItems)
}

func NewBat(p Pos) *Monster {
	items := randomizeLoot(p)
	return &Monster{Character{
		Entity:       Entity{Pos: p, Name: "Bat", Rune: Bat},
		Health:       50,
		MaxHealth:    50,
		MinDamage:    2,
		MaxDamage:    3,
		Critical:     0,
		Armor:        0,
		Speed:        2.0,
		ActionPoints: 0.0,
		Items:        items,
	}}
}

func NewRat(p Pos) *Monster {
	items := randomizeLoot(p)
	return &Monster{Character{
		Entity:       Entity{Pos: p, Name: "Rat", Rune: Rat},
		Health:       50,
		MaxHealth:    50,
		MinDamage:    1,
		MaxDamage:    2,
		Critical:     0,
		Armor:        0,
		Speed:        2.0,
		ActionPoints: 0.0,
		Items:        items,
	}}
}

func NewSpider(p Pos) *Monster {
	items := randomizeLoot(p)
	return &Monster{Character{
		Entity:       Entity{Pos: p, Name: "Spider", Rune: Spider},
		Health:       10,
		MaxHealth:    10,
		MinDamage:    2,
		MaxDamage:    4,
		Critical:     0,
		Armor:        0,
		Speed:        1.0,
		ActionPoints: 0.0,
		Items:        items,
	}}
}

func (m *Monster) Kill(level *Level) {
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
