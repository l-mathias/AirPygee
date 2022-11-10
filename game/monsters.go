package game

import (
	"math/rand"
	"time"
)

type Monster struct {
	Character
}

func randomizeLoot() *[]Item {
	items := make([]Item, 0)
	numItems := 0

	rand.Seed(time.Now().UnixNano())
	number := rand.Intn(100)

	switch {
	case number <= 2:
		numItems = 4
	case number > 2 && number <= 10:
		numItems = 3
	case number > 10 && number <= 20:
		numItems = 2
	case number > 20 && number <= 40:
		numItems = 1
	case number > 40 && number <= 100:
		numItems = 0
	}

	for i := 0; i < numItems; i++ {
		rand.Seed(time.Now().UnixNano())
		switch rand.Intn(2) {
		case 0:
			items = append(items, NewSword(Pos{}))
		case 1:
			items = append(items, NewHelmet(Pos{}))
		case 2:
			items = append(items, NewHealthPotion(Pos{}, "Small"))
		}
	}
	return &items
}

func NewRat(p Pos) *Monster {
	items := randomizeLoot()
	return &Monster{Character{
		Entity:       Entity{Pos: p, Name: "Rat", Rune: 'R'},
		Health:       5,
		MaxHealth:    5,
		MinDamage:    1,
		MaxDamage:    2,
		Armor:        0,
		Speed:        2.0,
		ActionPoints: 0.0,
		Items:        *items,
	}}
}

func NewSpider(p Pos) *Monster {
	items := randomizeLoot()
	return &Monster{Character{
		Entity:       Entity{Pos: p, Name: "Spider", Rune: 'S'},
		Health:       10,
		MaxHealth:    10,
		MinDamage:    2,
		MaxDamage:    4,
		Armor:        0,
		Speed:        1.0,
		ActionPoints: 0.0,
		Items:        *items,
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
