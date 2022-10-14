package ui2d

import (
	"AirPygee/game"
	"bufio"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const winWidth, winHeight = 1280, 720

var (
	renderer         *sdl.Renderer
	textureAtlas     *sdl.Texture
	textureIndex     map[game.Tile]sdl.Rect
	centerX, centerY int
)

type UI2d struct {
}

func loadTextureIndex() {
	textureIndex = make(map[game.Tile]sdl.Rect)
	infile, err := os.Open("ui2d/assets/atlas-index.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(infile)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		tileRune := game.Tile(line[0])
		xy := line[1:]
		splitXyC := strings.Split(xy, ",")
		x, err := strconv.ParseInt(strings.TrimSpace(splitXyC[0]), 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(strings.TrimSpace(splitXyC[1]), 10, 64)
		if err != nil {
			panic(err)
		}
		variationCount, err := strconv.ParseInt(strings.TrimSpace(splitXyC[2]), 10, 64)
		if err != nil {
			panic(err)
		}

		rand.Seed(time.Now().UnixNano())
		variation := rand.Int63n(variationCount)

		x += variation

		if x > 62 {
			x -= 62
			y++
		}

		rect := sdl.Rect{int32(x * 32), int32(y * 32), 32, 32}
		textureIndex[tileRune] = rect
	}

}

func imgFileToTexture(filename string) *sdl.Texture {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	image, err := img.Load(filename)
	if err != nil {
		panic(err)
	}

	tex, err := renderer.CreateTextureFromSurface(image)
	if err != nil {
		panic(err)
	}
	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}
	return tex
}

func init() {
	sdl.LogSetAllPriority(sdl.LOG_PRIORITY_VERBOSE)
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		log.Println(err)
		return
	}

	window, err := sdl.CreateWindow("AirPygee", 100, 100, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		log.Println(err)
		return
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Println(err)
		return
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	textureAtlas = imgFileToTexture("ui2d/assets/tiles.png")
	loadTextureIndex()

	centerX = -1
	centerY = -1
}

func (ui *UI2d) Draw(level *game.Level) {
	if centerX == -1 && centerY == -1 {
		centerX = level.Player.X
		centerY = level.Player.Y
	}

	limit := 5
	if level.Player.X > centerX+limit {
		centerX++
	} else if level.Player.X < centerX-limit {
		centerX--
	} else if level.Player.Y > centerY+limit {
		centerY++
	} else if level.Player.Y < centerY-limit {
		centerY--
	}
	offsetX := int32(winWidth/2 - centerX*32)
	offsetY := int32(winHeight/2 - centerY*32)

	err := renderer.Clear()
	if err != nil {
		panic(err)
	}
	for y, row := range level.Map {
		for x, tile := range row {
			srcRect := textureIndex[tile]
			dstRect := sdl.Rect{int32(x*32) + offsetX, int32(y*32) + offsetY, 32, 32}

			pos := game.Pos{X: x, Y: y}
			if level.Debug[pos] {
				err := textureAtlas.SetColorMod(128, 0, 0)
				if err != nil {
					panic(err)
				}
			} else {
				err := textureAtlas.SetColorMod(255, 255, 255)
				if err != nil {
					panic(err)
				}
			}

			if err := renderer.Copy(textureAtlas, &srcRect, &dstRect); err != nil {
				panic(err)
			}

		}
	}

	if err := renderer.Copy(textureAtlas, &sdl.Rect{21 * 32, 59 * 32, 32, 32}, &sdl.Rect{int32(level.Player.X*32) + offsetX, int32(level.Player.Y*32) + offsetY, 32, 32}); err != nil {
		panic(err)
	}
	renderer.Present()
}

func (ui *UI2d) GetInput() *game.Input {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			return &game.Input{Typ: game.Quit}
		case *sdl.KeyboardEvent:
			if e.Type != sdl.KEYDOWN {
				break
			}
			switch e.Keysym.Sym {
			case sdl.K_UP:
				return &game.Input{Typ: game.Up}
			case sdl.K_DOWN:
				return &game.Input{Typ: game.Down}
			case sdl.K_LEFT:
				return &game.Input{Typ: game.Left}
			case sdl.K_RIGHT:
				return &game.Input{Typ: game.Right}
			case sdl.K_s:
				return &game.Input{Typ: game.Search}
			default:
				return &game.Input{Typ: game.None}
			}
		}
		sdl.Delay(10)
	}
	return &game.Input{Typ: game.None}
}
