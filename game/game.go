package game

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

//TODO - improve loadWorld loadLevels - one should call the other one

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	QuitGame
	CloseWindow
	Action
	TakeAll
	TakeItem
)

type Game struct {
	LevelChans   []chan *Level
	InputChan    chan *Input
	Levels       map[string]*Level
	CurrentLevel *Level
}

func NewGame(numWindows int) *Game {
	levelChans := make([]chan *Level, numWindows)
	for i := range levelChans {
		levelChans[i] = make(chan *Level)
	}
	inputChan := make(chan *Input, 10)

	game := &Game{levelChans, inputChan, nil, nil}

	return game
}

type InputType int

type Tile struct {
	Rune        rune
	OverlayRune rune
	Visible     bool
	Seen        bool
}

type Input struct {
	Typ          InputType
	Item         *Item
	LevelChannel chan *Level
}

const (
	StoneWall  rune = '#'
	DirtFloor  rune = '.'
	ClosedDoor rune = '|'
	OpenDoor   rune = '/'
	DownStair  rune = 'd'
	UpStair    rune = 'u'
	Blank      rune = 0
	Pending    rune = -1
)

type Pos struct {
	X, Y int
}

type LevelPos struct {
	*Level
	Pos
}

type Entity struct {
	Pos
	Name string
	Rune rune
}

type Character struct {
	Entity
	Hitpoints     int
	MaxHitpoints  int
	Strength      int
	Speed         float64
	ActionPoints  float64
	SightRange    int
	Items         []*Item
	InventorySize int
}

type GameEvent int

const (
	Empty GameEvent = iota
	Move
	DoorOpen
	Attack
	Hit
	Portal
	Pickup
)

type Level struct {
	Map       [][]Tile
	Player    *Player
	Monsters  map[Pos]*Monster
	Portals   map[Pos]*LevelPos
	Items     map[Pos][]*Item
	Events    []string
	EventPos  int
	Debug     map[Pos]bool
	LastEvent GameEvent
}

func (c *Character) Pass() {
	c.ActionPoints -= c.Speed
}

func (level *Level) MoveItem(itemToMove *Item, character *Character) {
	pos := character.Pos
	items := level.Items[pos]
	for i, item := range items {
		if item == itemToMove {
			items = append(items[:i], items[i+1:]...)
			level.Items[pos] = items
			character.Items = append(character.Items, item)
			level.AddEvent(character.Name + " picked up " + itemToMove.Name)
			level.LastEvent = Pickup
			return
		}
	}
	panic("Tried to remove item we're not on top of")
}

func (level *Level) Attack(c1, c2 *Character) {
	c1.ActionPoints--
	c1AttackPower := c1.Strength
	c2.Hitpoints -= c1AttackPower

	if c2.Hitpoints > 0 {
		level.AddEvent(c1.Name + " attacked " + c2.Name + " for " + strconv.Itoa(c1AttackPower))
		level.LastEvent = Attack
	} else {
		level.AddEvent(c1.Name + " killed " + c2.Name)
	}
}

func (level *Level) lineOfSight() {
	pos := level.Player.Pos
	dist := level.Player.SightRange

	for y := pos.Y - dist; y <= pos.Y+dist; y++ {
		for x := pos.X - dist; x <= pos.X+dist; x++ {
			xDelta := pos.X - x
			yDelta := pos.Y - y
			d := math.Sqrt(float64(xDelta*xDelta + yDelta*yDelta))
			if d <= float64(dist) {
				level.bresenham(pos, Pos{x, y})
			}
		}
	}
}

func (level *Level) bresenham(start, end Pos) {
	steep := math.Abs(float64((end.Y)-start.Y)) > math.Abs(float64(end.X-start.X))
	if steep {
		start.X, start.Y = start.Y, start.X
		end.X, end.Y = end.Y, end.X
	}
	deltaY := int(math.Abs(float64(end.Y - start.Y)))
	err := 0
	y := start.Y
	yStep := 1
	if start.Y >= end.Y {
		yStep = -1
	}
	if start.X > end.X {
		deltaX := start.X - end.X

		for x := start.X; x > end.X; x-- {
			var pos Pos
			if steep {
				pos = Pos{y, x}
			} else {
				pos = Pos{x, y}
			}
			level.Map[pos.Y][pos.X].Seen = true
			level.Map[pos.Y][pos.X].Visible = true
			if !canSeeTrough(level, pos) {
				return
			}
			err += deltaY
			if 2*err >= deltaX {
				y += yStep
				err -= deltaX
			}
		}
	} else {
		deltaX := end.X - start.X

		for x := start.X; x < end.X; x++ {
			var pos Pos
			if steep {
				pos = Pos{y, x}
			} else {
				pos = Pos{x, y}
			}
			level.Map[pos.Y][pos.X].Seen = true
			level.Map[pos.Y][pos.X].Visible = true
			if !canSeeTrough(level, pos) {
				return
			}
			err += deltaY
			if 2*err >= deltaX {
				y += yStep
				err -= deltaX
			}
		}
	}
}

func (level *Level) AddEvent(event string) {
	level.Events[level.EventPos] = event
	level.EventPos++

	if level.EventPos == len(level.Events) {
		level.EventPos = 0
	}
}
func inRange(level *Level, pos Pos) bool {
	return pos.X < len(level.Map[0]) && pos.Y < len(level.Map) && pos.X >= 0 && pos.Y >= 0
}

func canWalk(level *Level, pos Pos) bool {
	if inRange(level, pos) {
		t := level.Map[pos.Y][pos.X]
		switch {
		case t.Rune == ClosedDoor, t.Rune == StoneWall, t.Rune == Blank, t.OverlayRune == ClosedDoor:
			return false
		}
		_, exists := level.Monsters[pos]
		if exists {
			return false
		}
		return true
	}
	return false
}

func checkDoor(level *Level, pos Pos) {
	if level.Map[pos.Y][pos.X].OverlayRune == ClosedDoor {
		level.Map[pos.Y][pos.X].OverlayRune = OpenDoor
		level.LastEvent = DoorOpen
		level.lineOfSight()
	}
}

func canSeeTrough(level *Level, pos Pos) bool {
	if inRange(level, pos) {
		t := level.Map[pos.Y][pos.X]
		switch {
		case t.Rune == ClosedDoor, t.Rune == StoneWall, t.Rune == Blank, t.OverlayRune == ClosedDoor:
			return false
		default:
			return true
		}
	}
	return false
}

func (game *Game) pickup(pos Pos, item *Item) {
	if item != nil {
		game.CurrentLevel.MoveItem(item, &game.CurrentLevel.Player.Character)
	} else {
		for _, item := range game.CurrentLevel.Items[pos] {
			game.CurrentLevel.MoveItem(item, &game.CurrentLevel.Player.Character)
		}
	}
}

func (game *Game) Move(to Pos) {
	level := game.CurrentLevel
	portal := level.Portals[to]
	if game.CurrentLevel.Portals[to] != nil {
		// transfer also events to new level
		events := game.CurrentLevel.Events
		eventPos := game.CurrentLevel.EventPos

		game.CurrentLevel = portal.Level
		game.CurrentLevel.Player.Pos = portal.Pos
		game.CurrentLevel.Events = events
		game.CurrentLevel.EventPos = eventPos
		game.CurrentLevel.lineOfSight()
	} else {
		game.CurrentLevel.Player.Pos = to
		level.LastEvent = Move
		for y, row := range game.CurrentLevel.Map {
			for x := range row {
				game.CurrentLevel.Map[y][x].Visible = false
			}
		}
		game.CurrentLevel.lineOfSight()
	}
}

func (game *Game) restart() {
	game.Levels = game.loadLevels()
	game.loadWorld()
	game.CurrentLevel.lineOfSight()
}

func (game *Game) dead() {
	game.restart()
}

func (game *Game) resolveMovement(pos Pos) {
	level := game.CurrentLevel
	monster, exists := game.CurrentLevel.Monsters[pos]
	if exists {
		game.CurrentLevel.Attack(&level.Player.Character, &monster.Character)
		if monster.Hitpoints <= 0 {
			monster.kill(level)
		}
		if game.CurrentLevel.Player.Hitpoints <= 0 {
			game.dead()
		}
	} else if canWalk(level, pos) {
		game.Move(pos)
	} else {
		checkDoor(level, pos)
	}
}

func (game *Game) handleInput(input *Input) {
	p := game.CurrentLevel.Player
	switch input.Typ {
	case Up:
		newPos := Pos{p.X, p.Y - 1}
		game.resolveMovement(newPos)
	case Down:
		newPos := Pos{p.X, p.Y + 1}
		game.resolveMovement(newPos)
	case Left:
		newPos := Pos{p.X - 1, p.Y}
		game.resolveMovement(newPos)
	case Right:
		newPos := Pos{p.X + 1, p.Y}
		game.resolveMovement(newPos)
	case Action:
		fmt.Println("Action launched")
	case TakeAll:
		game.pickup(game.CurrentLevel.Player.Pos, nil)
	case TakeItem:
		game.pickup(game.CurrentLevel.Player.Pos, input.Item)
	case CloseWindow:
		close(input.LevelChannel)
		chanIndex := 0
		for i, c := range game.LevelChans {
			if c == input.LevelChannel {
				chanIndex = i
				break
			}
		}
		game.LevelChans = append(game.LevelChans[:chanIndex], game.LevelChans[chanIndex+1:]...)
		if len(game.LevelChans) == 0 {
			os.Exit(1)
		}
	}
}

func (game *Game) loadWorld() {
	file, err := os.Open("game/maps/world.txt")
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	for rIndex, row := range rows {
		if rIndex == 0 {
			game.CurrentLevel = game.Levels[row[0]]
			if game.CurrentLevel == nil {
				fmt.Println("could'nt find level in world file")
				panic(nil)
			}
			continue
		}
		levelWithPortal := game.Levels[row[0]]
		if levelWithPortal == nil {
			fmt.Println("could'nt find level in world file")
			panic(nil)
		}
		x, err := strconv.Atoi(row[1])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(row[2])
		if err != nil {
			panic(err)
		}

		pos := Pos{x, y}

		levelToTeleport := row[3]

		x, err = strconv.Atoi(row[4])
		if err != nil {
			panic(err)
		}
		y, err = strconv.Atoi(row[5])
		if err != nil {
			panic(err)
		}

		posToTeleport := Pos{x, y}

		levelWithPortal.Portals[pos] = &LevelPos{
			Level: game.Levels[levelToTeleport],
			Pos:   posToTeleport,
		}
	}
}

func (game *Game) loadLevels() map[string]*Level {
	player := NewPlayer()

	levels := make(map[string]*Level, 0)

	fileNames, err := filepath.Glob("game/maps/*.map")
	if err != nil {
		panic(err)
	}

	for _, fileName := range fileNames {
		levelName := filepath.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))])

		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(file)
		levelLines := make([]string, 0)
		longestRow := 0
		index := 0
		for scanner.Scan() {
			levelLines = append(levelLines, scanner.Text())
			if len(levelLines[index]) > longestRow {
				longestRow = len(levelLines[index])
			}
			index++
		}

		level := &Level{}
		level.Events = make([]string, 15)
		level.Player = player
		level.Map = make([][]Tile, len(levelLines))
		level.Monsters = make(map[Pos]*Monster, 0)
		level.Portals = make(map[Pos]*LevelPos, 0)
		level.Items = make(map[Pos][]*Item, 0)
		level.LastEvent = -1

		for i := range level.Map {
			level.Map[i] = make([]Tile, longestRow)
		}

		for y := range levelLines {
			line := levelLines[y]
			for x, c := range line {
				pos := Pos{X: x, Y: y}
				switch c {
				case ' ', '\n', '\t', '\r':
					level.Map[y][x].Rune = Blank
				case '#':
					level.Map[y][x].Rune = StoneWall
				case '.':
					level.Map[y][x].Rune = DirtFloor
				case '|':
					level.Map[y][x].OverlayRune = ClosedDoor
					level.Map[y][x].Rune = Pending
				case '/':
					level.Map[y][x].OverlayRune = OpenDoor
					level.Map[y][x].Rune = Pending
				case 'd':
					level.Map[y][x].OverlayRune = DownStair
					level.Map[y][x].Rune = Pending
				case 'u':
					level.Map[y][x].OverlayRune = UpStair
					level.Map[y][x].Rune = Pending
				case 's':
					level.Items[pos] = append(level.Items[pos], NewSword(pos))
					level.Items[pos] = append(level.Items[pos], NewHelmet(pos))
					level.Map[y][x].Rune = Pending
				case 'h':
					level.Items[pos] = append(level.Items[pos], NewHelmet(pos))
					level.Map[y][x].Rune = Pending
				case '@':
					level.Player.X = x
					level.Player.Y = y
					level.Map[y][x].Rune = Pending
				case 'R':
					level.Monsters[pos] = NewRat(pos)
					level.Map[y][x].Rune = Pending
				case 'S':
					level.Monsters[pos] = NewSpider(pos)
					level.Map[y][x].Rune = Pending
				default:
					panic("invalid character in map")
				}
			}
		}

		for y, row := range level.Map {
			for x, tile := range row {
				if tile.Rune == Pending {
					level.Map[y][x].Rune = level.bfsFloor(Pos{x, y})
				}
			}
		}
		levels[levelName] = level
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}

	return levels
}

func getNeighbors(level *Level, pos Pos) []Pos {
	neighbors := make([]Pos, 0, 4)
	left := Pos{pos.X - 1, pos.Y}
	right := Pos{pos.X + 1, pos.Y}
	up := Pos{pos.X, pos.Y - 1}
	down := Pos{pos.X, pos.Y + 1}

	if canWalk(level, right) {
		neighbors = append(neighbors, right)
	}
	if canWalk(level, left) {
		neighbors = append(neighbors, left)
	}
	if canWalk(level, up) {
		neighbors = append(neighbors, up)
	}
	if canWalk(level, down) {
		neighbors = append(neighbors, down)
	}
	return neighbors
}

func (level *Level) bfsFloor(start Pos) rune {
	frontier := make([]Pos, 0, 8)
	frontier = append(frontier, start)
	visited := make(map[Pos]bool)
	visited[start] = true

	//level.Debug = visited

	for len(frontier) > 0 {
		current := frontier[0]

		currentTile := level.Map[current.Y][current.X]
		switch currentTile.Rune {
		case DirtFloor:
			return DirtFloor
		default:
		}

		frontier = frontier[1:]
		for _, next := range getNeighbors(level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
			}
		}
	}
	return DirtFloor
}

func (level *Level) astar(start Pos, goal Pos) []Pos {
	frontier := make(pqueue, 0, 8)
	frontier = frontier.push(start, 1)
	cameFrom := make(map[Pos]Pos)
	cameFrom[start] = start
	costSoFar := make(map[Pos]int)
	costSoFar[start] = 0

	//level.Debug = make(map[Pos]bool)

	var current Pos
	for len(frontier) > 0 {

		frontier, current = frontier.pop()

		if current == goal {
			path := make([]Pos, 0)
			p := current
			for p != start {
				path = append(path, p)
				p = cameFrom[p]
			}
			path = append(path, p)

			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}

			//level.Debug = make(map[Pos]bool)
			//for _, pos := range path {
			//	level.Debug[pos] = true
			//}
			return path
		}

		for _, next := range getNeighbors(level, current) {
			newCost := costSoFar[current] + 1
			_, exists := costSoFar[next]
			if !exists || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				xDist := int(math.Abs(float64(goal.X - next.X)))
				yDist := int(math.Abs(float64(goal.Y - next.Y)))
				priority := newCost + xDist + yDist
				frontier = frontier.push(next, priority)
				//level.Debug[next] = true
				cameFrom[next] = current
			}
		}
	}
	return nil
}

func (game *Game) Run() {
	game.Levels = game.loadLevels()
	game.loadWorld()
	game.CurrentLevel.lineOfSight()

	for _, lchan := range game.LevelChans {
		lchan <- game.CurrentLevel
	}

	for input := range game.InputChan {
		if input.Typ == QuitGame {
			return
		}
		game.handleInput(input)

		for _, monster := range game.CurrentLevel.Monsters {
			monster.Update(game)
			if game.CurrentLevel.Player.Hitpoints <= 0 {
				game.dead()
				break
			}
		}

		if len(game.LevelChans) == 0 {
			return
		}

		for _, lchan := range game.LevelChans {
			lchan <- game.CurrentLevel
		}

	}
}
