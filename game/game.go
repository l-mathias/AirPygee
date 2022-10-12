package game

import (
	"bufio"
	"os"
)

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	Quit
)

type Input struct {
	Typ InputType
}

type GameUI interface {
	Draw(*Level)
	GetInput() *Input
}

type Tile rune

const (
	StoneWall Tile = '#'
	DirtFloor Tile = '.'
	Door      Tile = '|'
	Blank     Tile = ' '
	Pending   Tile = -1
)

type Entity struct {
	X, Y int
}

type Player struct {
	Entity
}
type Level struct {
	Map    [][]Tile
	Player Player
}

func LoadLevelFromFile(fileName string) *Level {
	file, err := os.Open(fileName)

	if err != nil {
		panic(err)
	}

	defer file.Close()

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
	level.Map = make([][]Tile, len(levelLines))

	for i := range level.Map {
		level.Map[i] = make([]Tile, longestRow)
	}

	for i := range levelLines {
		line := levelLines[i]
		for j, c := range line {
			switch c {
			case ' ', '\n', '\t', '\r':
				level.Map[i][j] = Blank
			case '#':
				level.Map[i][j] = StoneWall
			case '.':
				level.Map[i][j] = DirtFloor
			case '|':
				level.Map[i][j] = Door
			case 'P':
				level.Player.X = j
				level.Player.Y = i
				level.Map[i][j] = Pending
			default:
				panic("invalid character in map")
			}
		}
	}

	for y, row := range level.Map {
		for x, tile := range row {
			if tile == Pending {
			SearchLoop:
				for searchX := x - 1; searchX <= x+1; searchX++ {
					for searchY := y - 1; searchY <= y+1; searchY++ {
						searchTile := level.Map[searchY][searchY]
						switch searchTile {
						case DirtFloor:
							level.Map[y][x] = DirtFloor
							break SearchLoop
						}
					}
				}
			}
		}
	}

	return level
}

func Run(ui GameUI) {
	level := LoadLevelFromFile("game/maps/level1.map")

	for {
		ui.Draw(level)
		input := ui.GetInput()

		switch input.Typ {
		case Up:
			level.Player.Y--
		case Down:
			level.Player.Y++
		case Left:
			level.Player.X--
		case Right:
			level.Player.X++
		case Quit:
			return
		}
	}
}
