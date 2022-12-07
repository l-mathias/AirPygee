package game

//TODO - improve loadWorld loadLevels - one should call the other one
//TODO - Save / Load game state
//TODO - Hero classes + characters interfaces
//TODO - levels procedural generation
//TODO - random monsters placed randomly in a level
//TODO - Chests

import (
	"bufio"
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
)

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
	Equip
	Drop
	Restart
	SetDifficulty
)

type Game struct {
	LevelChans   []chan *Level
	InputChan    chan *Input
	Levels       map[string]*Level
	CurrentLevel *Level
	Difficulty   int
}

func NewGame(numWindows int) *Game {
	levelChans := make([]chan *Level, numWindows)
	for i := range levelChans {
		levelChans[i] = make(chan *Level)
	}
	inputChan := make(chan *Input, 10)

	game := &Game{LevelChans: levelChans, InputChan: inputChan, Levels: nil, CurrentLevel: nil, Difficulty: 1}

	return game
}

type InputType int

type Tile struct {
	Rune        rune
	OverlayRune rune
	AnimRune    rune
	Visible     bool
	Seen        bool
	Walkable    bool
	Actionable  bool
}

type Input struct {
	Typ          InputType
	Item         Item
	LevelChannel chan *Level
	Difficulty   int
}

// normal Tiles
const (
	StoneWall rune = '#'
	DirtFloor rune = '.'
	Blank     rune = 0
	Pending   rune = -1
)

// monsters tiles
const (
	Rat    rune = 'R'
	Spider rune = 'S'
	Bat    rune = 'B'
)

// Overlay tiles
const (
	ClosedDoor rune = '|'
	OpenDoor   rune = '/'
	DownStair  rune = 'd'
	UpStair    rune = 'u'
)

// animations
const (
	LeftAnim       rune = 'L'
	RightAnim      rune = 'R'
	DownAnim       rune = 'D'
	UpAnim         rune = 'U'
	AnimatedPortal rune = 'a'
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
	CameFrom    Pos
	WantedTo    Pos
	Name        string
	Rune        rune
	Type        ItemType
	Description string
}

type Character struct {
	Entity
	Health        int
	MaxHealth     int
	MinDamage     int
	MaxDamage     int
	Armor         int
	Critical      float64
	Speed         float64
	ActionPoints  float64
	SightRange    int
	EquippedItems []EquipableItem
	Items         []Item
	InventorySize int
}

type GameEvent int

const (
	Empty GameEvent = iota
	Move
	DoorOpen
	DoorClose
	Attack
	Pickup
	DropItem
	ConsumePotion
	OpenChest
)

type GameAttack struct {
	Damage     int
	IsCritical bool
	Who        *Character
}

type Level struct {
	Map        [][]Tile
	Player     *Player
	Monsters   map[Pos]*Monster
	Portals    map[Pos]*LevelPos
	Items      map[Pos][]Item
	Events     []string
	EventPos   int
	Debug      map[Pos]bool
	LastEvent  GameEvent
	LastAttack GameAttack
}

func (c *Character) Pass() {
	c.ActionPoints -= c.Speed
}

func (level *Level) MoveItem(itemToMove Item, character *Character) {
	pos := character.Pos
	for i, item := range level.Items[pos] {
		if item == itemToMove {
			if len(level.Player.Items) < level.Player.InventorySize {
				level.Items[pos] = append(level.Items[pos][:i], level.Items[pos][i+1:]...)
				character.Items = append(character.Items, item)
				level.AddEvent(character.Name + " picked up:" + item.GetName())
				level.LastEvent = Pickup
				return
			} else {
				level.AddEvent("Inventory full")
				return
			}
		}
	}
	panic("Tried to move an item we were not on top of")
}

func randomizeDamage(min, max int) int {
	number, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	CheckError(err)
	return int(number.Int64()) + min
}

func isCritical(crit float64) bool {
	number, err := rand.Int(rand.Reader, big.NewInt(100))
	CheckError(err)

	return float64(number.Int64()) <= crit
}

func (level *Level) Attack(c1, c2 *Character) {
	c1.ActionPoints--
	c1AttackPower := randomizeDamage(c1.MinDamage, c1.MaxDamage)
	damageDealt := c1AttackPower - c2.Armor
	if damageDealt < 0 {
		damageDealt = 0
	}

	level.LastAttack.Who = c2
	if isCritical(c1.Critical) {
		c1AttackPower *= 2
		level.LastAttack.IsCritical = true
	} else {
		level.LastAttack.IsCritical = false
	}
	c2.Health -= damageDealt
	level.LastAttack.Damage = damageDealt

	if c2.Health > 0 {
		level.AddEvent(c1.Name + " attacked " + c2.Name + " for " + strconv.Itoa(damageDealt))
	} else {
		level.AddEvent(c1.Name + " killed " + c2.Name)
	}
	level.LastEvent = Attack
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
		case !t.Walkable:
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
		level.Map[pos.Y][pos.X].Walkable = true
		level.LastEvent = DoorOpen
		level.lineOfSight()
	} else if level.Map[pos.Y][pos.X].OverlayRune == OpenDoor {
		level.Map[pos.Y][pos.X].OverlayRune = ClosedDoor
		level.Map[pos.Y][pos.X].Walkable = false
		level.LastEvent = DoorClose
		level.lineOfSight()
	}
}

func canSeeTrough(level *Level, pos Pos) bool {
	if inRange(level, pos) {
		t := level.Map[pos.Y][pos.X]
		switch {
		case !t.Walkable:
			return false
		default:
			return true
		}
	}
	return false
}

// pickup if nil, we'll take all the objects on the ground
func (game *Game) pickup(item Item) {
	if item != nil {
		game.CurrentLevel.MoveItem(item, &game.CurrentLevel.Player.Character)
	} else {
		pos := game.CurrentLevel.Player.Pos
		for i := len(game.CurrentLevel.Items[pos]) - 1; i >= 0; i-- {
			game.CurrentLevel.MoveItem(game.CurrentLevel.Items[pos][i], &game.CurrentLevel.Player.Character)
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

func (game *Game) Restart() {
	game.Levels = game.loadLevels()
	game.loadWorld()
	game.CurrentLevel.lineOfSight()
}

func (game *Game) Dead() {
	game.Restart()
}

func (game *Game) resolveMovement(pos Pos) {
	level := game.CurrentLevel
	monster, exists := game.CurrentLevel.Monsters[pos]
	game.CurrentLevel.Player.CameFrom = game.CurrentLevel.Player.Pos
	if exists {
		level.Player.WantedTo = pos
		game.CurrentLevel.Attack(&level.Player.Character, &monster.Character)
		if monster.Health <= 0 {
			monster.Kill(level)
		}
		if game.CurrentLevel.Player.Health <= 0 {
			game.Dead()
		}
	} else if canWalk(level, pos) {
		level.Player.WantedTo = pos
		game.Move(pos)
	} else {
		level.Player.WantedTo = pos
	}
}

func (game *Game) heal(c *Character, hp int) {
	c.Health += hp
	if c.Health > c.MaxHealth {
		c.Health = c.MaxHealth
	}
}

func (game *Game) action(pos Pos, item Item) {
	switch {
	case game.CurrentLevel.Map[pos.Y][pos.X].OverlayRune == ClosedDoor:
		checkDoor(game.CurrentLevel, pos)
	case game.CurrentLevel.Map[pos.Y][pos.X].OverlayRune == OpenDoor:
		checkDoor(game.CurrentLevel, pos)
	case item != nil:
		switch item.(type) {
		case ConsumableItem:
			game.consumePotion(item.(ConsumableItem))
		case OpenableItem:
			game.OpenItem(item.(OpenableItem))
		default:
		}
	}
}

func (level *Level) FrontOf() Pos {
	cameFrom := &level.Player.CameFrom
	currentPos := &level.Player.Pos
	switch {
	case cameFrom.X < currentPos.X:
		return Pos{currentPos.X + 1, currentPos.Y}
	case cameFrom.X > currentPos.X:
		return Pos{currentPos.X - 1, currentPos.Y}
	case cameFrom.Y < currentPos.Y:
		return Pos{currentPos.X, currentPos.Y + 1}
	case cameFrom.Y > currentPos.Y:
		return Pos{currentPos.X, currentPos.Y - 1}
	case *cameFrom == *currentPos:
		return level.Player.WantedTo
	default:
		return Pos{}
	}
}

func (game *Game) removeInventoryItem(itemToRemove Item, character *Character) {
	for i, item := range game.CurrentLevel.Player.Items {
		if item == itemToRemove {
			character.Items = append(game.CurrentLevel.Player.Items[:i], game.CurrentLevel.Player.Items[i+1:]...)
			return
		}
	}
	panic("Tried to drop bad item")
}

func (game *Game) dropItem(itemToDrop Item, character *Character) {
	for i, item := range game.CurrentLevel.Player.Items {
		if item == itemToDrop {
			character.Items = append(game.CurrentLevel.Player.Items[:i], game.CurrentLevel.Player.Items[i+1:]...)
			game.CurrentLevel.Items[character.Pos] = append(game.CurrentLevel.Items[character.Pos], itemToDrop)
			game.CurrentLevel.AddEvent(character.Name + " dropped " + itemToDrop.GetName())
			game.CurrentLevel.LastEvent = DropItem
			return
		}
	}
	panic("Tried to drop bad item")
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
		game.action(game.CurrentLevel.FrontOf(), input.Item)
	case TakeAll:
		game.pickup(nil)
	case TakeItem:
		game.pickup(input.Item)
	case Equip:
		if input.Item.(EquipableItem).IsEquipped() {
			game.unEquip(input.Item.(EquipableItem))
		} else {
			game.equip(input.Item.(EquipableItem))
		}
	case SetDifficulty:
		game.Difficulty = input.Difficulty
	case Drop:
		game.dropItem(input.Item, &game.CurrentLevel.Player.Character)
	case Restart:
		game.Restart()
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
	CheckError(err)

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true
	rows, err := csvReader.ReadAll()
	CheckError(err)

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
		CheckError(err)
		y, err := strconv.Atoi(row[2])
		CheckError(err)

		pos := Pos{x, y}

		levelToTeleport := row[3]

		x, err = strconv.Atoi(row[4])
		CheckError(err)
		y, err = strconv.Atoi(row[5])
		CheckError(err)

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
	CheckError(err)

	for _, fileName := range fileNames {
		levelName := filepath.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))])

		file, err := os.Open(fileName)
		CheckError(err)

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
		level.Items = make(map[Pos][]Item, 0)
		level.LastEvent = -1

		for i := range level.Map {
			level.Map[i] = make([]Tile, longestRow)
		}

		for y := range levelLines {
			line := levelLines[y]
			for x, c := range line {
				pos := Pos{X: x, Y: y}
				level.Map[y][x].Walkable = true
				level.Map[y][x].Actionable = false
				level.Map[y][x].AnimRune = Blank
				switch c {
				case ' ', '\n', '\t', '\r':
					level.Map[y][x].Rune = Blank
					level.Map[y][x].Walkable = false
				case '#':
					level.Map[y][x].Rune = StoneWall
					level.Map[y][x].Walkable = false
				case '.':
					level.Map[y][x].Rune = DirtFloor
				case '|':
					level.Map[y][x].OverlayRune = ClosedDoor
					level.Map[y][x].Rune = Pending
					level.Map[y][x].Walkable = false
					level.Map[y][x].Actionable = true
				case '/':
					level.Map[y][x].OverlayRune = OpenDoor
					level.Map[y][x].Rune = Pending
					level.Map[y][x].Actionable = true
				case 'd':
					level.Map[y][x].OverlayRune = DownStair
					level.Map[y][x].Rune = Pending
				case 'u':
					level.Map[y][x].OverlayRune = UpStair
					level.Map[y][x].Rune = Pending
				case 's':
					level.Items[pos] = append(level.Items[pos], NewSword(pos))
					level.Map[y][x].Rune = Pending
				//case 'B':
				//	level.Items[pos] = append(level.Items[pos], NewSword(pos))
				//	level.Map[y][x].Rune = Pending
				case 'h':
					level.Items[pos] = append(level.Items[pos], NewHelmet(pos))
					level.Map[y][x].Rune = Pending
				case 'b':
					level.Items[pos] = append(level.Items[pos], NewBoots(pos))
					level.Map[y][x].Rune = Pending
				case 't':
					level.Items[pos] = append(level.Items[pos], NewTreasureChest(pos, 3))
					level.Map[y][x].Rune = Pending
					level.Map[y][x].Walkable = false
					level.Map[y][x].Actionable = true
				case 'a':
					level.Items[pos] = append(level.Items[pos], NewPlate(pos))
					level.Map[y][x].Rune = Pending
				case 'p':
					level.Items[pos] = append(level.Items[pos], NewHealthPotion(pos, "Small"))
					level.Map[y][x].Rune = Pending
				case '@':
					level.Player.Pos = pos
					level.Map[y][x].Rune = Pending
				case 'B':
					level.Monsters[pos] = NewBat(pos)
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

		game.randomizeLevel(level)
		levels[levelName] = level
		err = file.Close()
		CheckError(err)
	}

	return levels
}

func countValidPositions(level *Level) int {
	count := 0
	for y := range level.Map {
		line := level.Map[y]
		for _, c := range line {
			switch c.Walkable {
			case true:
				count++
			}
		}
	}
	return count
}

func (game *Game) randomizeLevel(level *Level) {
	numChests := countValidPositions(level) * game.Difficulty / 100
	randomizeChests(numChests, level)
	numMonsters := countValidPositions(level) * game.Difficulty / 100
	randomizeMonsters(numMonsters, level)
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

func findValidPosition(level *Level) Pos {
	posList := make([]Pos, 0)
	for y := range level.Map {
		line := level.Map[y]
		for x, c := range line {
			switch c.Walkable {
			case true:
				posList = append(posList, Pos{X: x, Y: y})
			}
		}
	}
	randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(posList))))
	CheckError(err)
	return posList[randIndex.Int64()]
}

func randomChest() int {
	randIndex, err := rand.Int(rand.Reader, big.NewInt(100))
	CheckError(err)

	switch {
	case randIndex.Int64() < 2:
		return 8
	case randIndex.Int64() < 5:
		return 7
	case randIndex.Int64() < 7:
		return 6
	case randIndex.Int64() < 15:
		return 5
	case randIndex.Int64() < 20:
		return 4
	case randIndex.Int64() < 30:
		return 3
	case randIndex.Int64() < 40:
		return 2
	case randIndex.Int64() <= 100:
		return 1
	}
	return 0
}

func randomizeChests(numChests int, level *Level) {
	for i := 0; i < numChests; i++ {
		randPos := findValidPosition(level)
		randSize := randomChest()
		level.Items[randPos] = append(level.Items[randPos], NewTreasureChest(randPos, randSize))
		level.Map[randPos.Y][randPos.X].Walkable = false
		level.Map[randPos.Y][randPos.X].Actionable = true
	}
}

func randomMonster(p Pos) *Monster {
	number, err := rand.Int(rand.Reader, big.NewInt(3))
	CheckError(err)

	switch {
	case number.Int64() == 0:
		return NewBat(p)
	case number.Int64() == 1:
		return NewSpider(p)
	case number.Int64() == 2:
		return NewRat(p)
	}
	return NewRat(p)
}

func randomizeMonsters(numMonsters int, level *Level) {
	for i := 0; i < numMonsters; i++ {
		randPos := findValidPosition(level)
		level.Monsters[randPos] = randomMonster(randPos)
	}
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
			if game.CurrentLevel.Player.Health <= 0 {
				game.Dead()
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
