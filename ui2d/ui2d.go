package ui2d

//TODO - add damage on top of character when combat
//TODO - add Player character selection
//TODO - improve life gauge using shading
//TODO - Improve fog of war effect using transparent texture or special tiles
//TODO - Fix UI bug when dead, should recenter camera

import (
	"AirPygee/game"
	"bufio"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"
)

const (
	FontSmall FontSize = iota
	FontMedium
	FontLarge
	musicVolume  int = 32
	soundsVolume int = 10
)

type uiState int

const (
	UIMain uiState = iota
	UIInventory
)

type ui struct {
	state                                     uiState
	sounds                                    sounds
	winWidth, winHeight                       int
	renderer                                  *sdl.Renderer
	window                                    *sdl.Window
	textureAtlas                              *sdl.Texture
	tileMap                                   *sdl.Texture
	textureIndex                              map[rune][]sdl.Rect
	centerX, centerY                          int
	r                                         *rand.Rand
	levelChan                                 chan *game.Level
	inputChan                                 chan *game.Input
	offsetX, offsetY                          int32
	eventBackground                           *sdl.Texture
	itemsInventoryBackground                  *sdl.Texture
	fontSmall, fontMedium, fontLarge          *ttf.Font
	str2TexSmall, str2TexMedium, str2TexLarge map[string]*sdl.Texture
}

func NewUI(inputChan chan *game.Input, levelChan chan *game.Level) *ui {
	ui := &ui{}
	ui.state = UIMain
	ui.inputChan = inputChan
	ui.levelChan = levelChan
	ui.str2TexSmall = make(map[string]*sdl.Texture)
	ui.str2TexMedium = make(map[string]*sdl.Texture)
	ui.str2TexLarge = make(map[string]*sdl.Texture)
	ui.r = rand.New(rand.NewSource(1))
	ui.winWidth = 1280
	ui.winHeight = 720
	window, err := sdl.CreateWindow("AirPygee", 100, 100, int32(ui.winWidth), int32(ui.winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	ui.window = window

	ui.renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")

	ui.textureAtlas = ui.imgFileToTexture("ui2d/assets/tiles.png")
	ui.tileMap = ui.imgFileToTexture("ui2d/assets/tilemap.png")

	ui.loadTextureIndex()

	ui.centerX = -1
	ui.centerY = -1

	ui.fontSmall, err = ttf.OpenFont("ui2d/assets/Kingthings_Foundation.ttf", int(float64(ui.winWidth)*0.015))
	if err != nil {
		panic(err)
	}
	ui.fontMedium, err = ttf.OpenFont("ui2d/assets/Kingthings_Foundation.ttf", 32)
	if err != nil {
		panic(err)
	}
	ui.fontLarge, err = ttf.OpenFont("ui2d/assets/Kingthings_Foundation.ttf", 64)
	if err != nil {
		panic(err)
	}

	ui.itemsInventoryBackground = ui.getSinglePixel(sdl.Color{R: 45, G: 45, B: 45, A: 128})
	if err = ui.itemsInventoryBackground.SetBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		panic(err)
	}

	ui.eventBackground = ui.getSinglePixel(sdl.Color{R: 45, G: 45, B: 45, A: 128})
	if err = ui.eventBackground.SetBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		panic(err)
	}

	ui.loadSounds()

	return ui
}

func (ui *ui) loadSounds() {
	if err := mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		panic(err)
	}

	mus, err := mix.LoadMUS("ui2d/assets/audio/music/the_field_of_dreams.mp3")
	mix.VolumeMusic(musicVolume)

	if err != nil {
		panic(err)
	}

	err = mus.Play(-1)
	if err != nil {
		panic(err)
	}

	ui.sounds.footstep = buildSoundsVariations("ui2d/assets/audio/sounds/Kenney/footstep*.ogg")
	ui.sounds.openDoor = buildSoundsVariations("ui2d/assets/audio/sounds/Kenney/doorOpen*.ogg")
	ui.sounds.swing = buildSoundsVariations("ui2d/assets/audio/sounds/battle/swing*.wav")
	ui.sounds.pickup = buildSoundsVariations("ui2d/assets/audio/sounds/Kenney/cloth*.ogg")
}

func buildSoundsVariations(pattern string) []*mix.Chunk {
	fileNames, err := filepath.Glob(pattern)
	result := make([]*mix.Chunk, 0)
	if err != nil {
		panic(err)
	}

	for _, fileName := range fileNames {
		sound, err := mix.LoadWAV(fileName)
		if err != nil {
			panic(err)
		}
		result = append(result, sound)
	}

	return result
}

func playRandomSound(chunks []*mix.Chunk, volume int) {
	chunkIndex := rand.Intn(len(chunks))
	chunks[chunkIndex].Volume(volume)
	chunks[chunkIndex].Play(-1, 0)
}

type sounds struct {
	openDoor []*mix.Chunk
	footstep []*mix.Chunk
	swing    []*mix.Chunk
	pickup   []*mix.Chunk
}

type FontSize int

func (ui *ui) stringToTexture(s string, color sdl.Color, size FontSize) *sdl.Texture {

	var font *ttf.Font
	switch size {
	case FontSmall:
		font = ui.fontSmall
		if tex, exists := ui.str2TexSmall[s]; exists {
			return tex
		}
	case FontMedium:
		font = ui.fontMedium
		if tex, exists := ui.str2TexMedium[s]; exists {
			return tex
		}
	case FontLarge:
		font = ui.fontLarge
		if tex, exists := ui.str2TexLarge[s]; exists {
			return tex
		}
	}
	fontSurface, err := font.RenderUTF8Blended(s, color)
	if err != nil {
		panic(err)
	}
	defer fontSurface.Free()

	fontTexture, err := ui.renderer.CreateTextureFromSurface(fontSurface)
	if err != nil {
		panic(err)
	}

	switch size {
	case FontSmall:
		ui.str2TexSmall[s] = fontTexture
	case FontMedium:
		ui.str2TexMedium[s] = fontTexture
	case FontLarge:
		ui.str2TexLarge[s] = fontTexture
	}

	return fontTexture
}

func (ui *ui) loadTextureIndex() {
	ui.textureIndex = make(map[rune][]sdl.Rect)
	infile, err := os.Open("ui2d/assets/atlas-index.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		tileRune := rune(line[0])
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

		var rects []sdl.Rect
		for i := int64(0); i < variationCount; i++ {
			rects = append(rects, sdl.Rect{X: int32(x * 32), Y: int32(y * 32), W: 32, H: 32})
			x++
			if x > 62 {
				x = 0
				y++
			}
		}
		ui.textureIndex[tileRune] = rects
	}

}

func (ui *ui) imgFileToTexture(filename string) *sdl.Texture {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	image, err := img.Load(filename)
	if err != nil {
		panic(err)
	}

	tex, err := ui.renderer.CreateTextureFromSurface(image)
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

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
		return
	}

	if err = ttf.Init(); err != nil {
		return
	}

	if err = mix.Init(mix.INIT_MP3); err != nil {
		panic(err)
	}

	if err = mix.Init(mix.INIT_OGG); err != nil {
		panic(err)
	}

}

func (ui *ui) LoadPlayer(level *game.Level, r *sdl.Renderer) {
	image, err := img.Load(level.Player.Filename)
	if err != nil {
		panic(err)
	}
	defer image.Free()

	level.Player.Texture, err = r.CreateTextureFromSurface(image)
	if err != nil {
		panic(err)
	}

	_, _, imageWidth, imageHeight, _ := level.Player.Texture.Query()
	level.Player.Width = imageWidth / level.Player.FramesX
	level.Player.Height = imageHeight / level.Player.FramesY
}

func (ui *ui) drawPlayer(level *game.Level) {
	p := level.Player
	p.FromY = p.CurrentFrame * p.Width

	p.Src = sdl.Rect{X: p.FromX, Y: p.FromY, W: p.Width, H: p.Height}
	p.Dest = sdl.Rect{X: int32(p.X*32) + ui.offsetX - ((p.Width - 32) / 2), Y: int32(p.Y*32) + ui.offsetY - ((p.Height - 32) / 2), W: p.Width, H: p.Height}

	err := ui.renderer.Copy(p.Texture, &p.Src, &p.Dest)
	if err != nil {
		panic(err)
	}
}

func (ui *ui) DrawInventory(level *game.Level) {
	ui.renderer.Copy(ui.itemsInventoryBackground, nil, &sdl.Rect{
		X: 100,
		Y: 100,
		W: 500,
		H: 500,
	})
}

func (ui *ui) Draw(level *game.Level) {
	if ui.centerX == -1 && ui.centerY == -1 {
		ui.LoadPlayer(level, ui.renderer)
		ui.centerX = level.Player.X
		ui.centerY = level.Player.Y
	}

	limit := 5

	if level.Player.X > ui.centerX+limit {
		diff := level.Player.X - (ui.centerX + limit)
		ui.centerX += diff
	} else if level.Player.X < ui.centerX-limit {
		diff := (ui.centerX - limit) - level.Player.X
		ui.centerX -= diff
	} else if level.Player.Y > ui.centerY+limit {
		diff := level.Player.Y - (ui.centerY + limit)
		ui.centerY += diff
	} else if level.Player.Y < ui.centerY-limit {
		diff := (ui.centerY - limit) - level.Player.Y
		ui.centerY -= diff
	}
	ui.offsetX = int32(ui.winWidth/2 - ui.centerX*32)
	ui.offsetY = int32(ui.winHeight/2 - ui.centerY*32)

	err := ui.renderer.Clear()
	if err != nil {
		panic(err)
	}
	ui.r.Seed(1)
	for y, row := range level.Map {
		for x, tile := range row {
			if tile.Rune != game.Blank {
				srcRects := ui.textureIndex[tile.Rune]
				srcRect := srcRects[ui.r.Intn(len(srcRects))]
				if tile.Visible || tile.Seen {
					dstRect := sdl.Rect{X: int32(x*32) + ui.offsetX, Y: int32(y*32) + ui.offsetY, W: 32, H: 32}
					pos := game.Pos{X: x, Y: y}
					if level.Debug[pos] {
						if err := ui.textureAtlas.SetColorMod(128, 0, 0); err != nil {
							panic(err)
						}
					} else if tile.Seen && !tile.Visible {
						if err := ui.textureAtlas.SetColorMod(128, 128, 128); err != nil {
							panic(err)
						}
					} else {
						if err := ui.textureAtlas.SetColorMod(255, 255, 255); err != nil {
							panic(err)
						}
					}
					if err := ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect); err != nil {
						panic(err)
					}

					if tile.OverlayRune != game.Blank {
						// TODO - if multiple variants of a tile, adapt srcRects
						srcRect := ui.textureIndex[tile.OverlayRune][0]
						if err := ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect); err != nil {
							panic(err)
						}
					}
				}
			}
		}
	}

	ui.displayMonsters(level)
	ui.displayItems(level)
	ui.drawPlayer(level)
	ui.displayStats(level)
	ui.displayEvents(level)

	// display item we are on top of
	if len(level.Items[level.Player.Pos]) > 0 {
		groundItems := level.Items[level.Player.Pos]
		for i, item := range groundItems {
			itemSrcRect := ui.textureIndex[item.Rune][0]
			err := ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{X: int32(ui.winWidth - 32 - i*32), Y: 0, W: 32, H: 32})
			if err != nil {
				panic(err)
			}
		}
		// drawing help letter T
		if err := ui.renderer.Copy(ui.tileMap, &sdl.Rect{X: 358, Y: 34, W: 16, H: 16}, &sdl.Rect{X: int32(ui.winWidth - 32 - len(groundItems)*32), Y: 0, W: 32, H: 32}); err != nil {
			panic(err)
		}
	}
	ui.renderer.Present()
}

func (ui *ui) displayEvents(level *game.Level) {
	textStart := int32(float64(ui.winHeight) * .68)
	textWidth := int32(float64(ui.winWidth) * .25)
	err := ui.renderer.Copy(ui.eventBackground, nil, &sdl.Rect{X: 0, Y: textStart, W: textWidth, H: int32(ui.winHeight) - textStart})
	if err != nil {
		panic(err)
	}

	i := level.EventPos
	count := 0
	_, fontSizeY, _ := ui.fontSmall.SizeUTF8("A")
	for {
		event := level.Events[i]
		if event != "" {
			tex := ui.stringToTexture(event, sdl.Color{R: 255}, FontSmall)
			_, _, w, h, _ := tex.Query()
			err := ui.renderer.Copy(tex, nil, &sdl.Rect{Y: int32(count*fontSizeY) + textStart, W: w, H: h})
			if err != nil {
				panic(err)
			}

		}
		i = (i + 1) % (len(level.Events))
		count++
		if i == level.EventPos {
			break
		}
	}

	// draw inventory
	inventoryStart := int32(float64(ui.winWidth) * 0.7)
	inventoryWidth := int32(ui.winWidth) - inventoryStart

	err = ui.renderer.Copy(ui.itemsInventoryBackground, nil, &sdl.Rect{X: inventoryStart, Y: int32(ui.winHeight - 32), W: inventoryWidth, H: 32})
	if err != nil {
		panic(err)
	}

	items := level.Player.Items
	for i, item := range items {
		itemSrcRect := ui.textureIndex[item.Rune][0]
		err := ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{X: int32(ui.winWidth - 32 - i*32), Y: int32(ui.winHeight - 32), W: 32, H: 32})
		if err != nil {
			panic(err)
		}
	}

}

func (ui *ui) displayItems(level *game.Level) {
	// Display Items
	for pos, items := range level.Items {
		if level.Map[pos.Y][pos.X].Visible {
			for _, item := range items {
				itemSrcRect := ui.textureIndex[item.Rune][0]
				err := ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{X: int32(pos.X*32) + ui.offsetX, Y: int32(pos.Y*32) + ui.offsetY, W: 32, H: 32})
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func (ui *ui) displayMonsters(level *game.Level) {
	// Display Monsters
	if err := ui.textureAtlas.SetColorMod(255, 255, 255); err != nil {
		panic(err)
	}
	for pos, monster := range level.Monsters {
		if level.Map[pos.Y][pos.X].Visible {

			if err := ui.renderer.Copy(ui.textureAtlas, &sdl.Rect{X: 928, Y: 1600, W: 32, H: 32}, &sdl.Rect{X: int32(level.Monsters[pos].X*32) + ui.offsetX, Y: int32((level.Monsters[pos].Y-1)*32) + ui.offsetY + 20, W: 32, H: 5}); err != nil {
				panic(err)
			}
			var gauge float64
			gauge = float64(level.Monsters[pos].Hitpoints) / float64(level.Monsters[pos].MaxHitpoints)

			if err := ui.renderer.Copy(ui.textureAtlas, &sdl.Rect{X: 1024, Y: 1600, W: 32, H: 32}, &sdl.Rect{X: int32(level.Monsters[pos].X*32) + ui.offsetX, Y: int32((level.Monsters[pos].Y-1)*32) + ui.offsetY + 20, W: int32(32 * gauge), H: 5}); err != nil {
				panic(err)
			}

			monsterSrcRect := ui.textureIndex[monster.Rune][0]
			err := ui.renderer.Copy(ui.textureAtlas, &monsterSrcRect, &sdl.Rect{X: int32(pos.X*32) + ui.offsetX, Y: int32(pos.Y*32) + ui.offsetY, W: 32, H: 32})
			if err != nil {
				panic(err)
			}
		}
	}
}

func (ui *ui) displayStats(level *game.Level) {
	firstFrameX := 512
	for i := 0; i < 4; i++ {
		if err := ui.renderer.Copy(ui.tileMap, &sdl.Rect{X: int32(firstFrameX + i*16), Y: 68, W: 16, H: 16}, &sdl.Rect{X: int32(i * 32), Y: 0, W: 32, H: 32}); err != nil {
			panic(err)
		}
	}

	// Move instruction after arrows
	tex := ui.stringToTexture("Move", sdl.Color{R: 255}, FontSmall)
	_, _, w, h, _ := tex.Query()
	err := ui.renderer.Copy(tex, nil, &sdl.Rect{X: 144, Y: 8, W: w, H: h})
	if err != nil {
		panic(err)
	}

	// Life symbol
	if err := ui.renderer.Copy(ui.textureAtlas, &sdl.Rect{X: 32, Y: 0, W: 32, H: 32}, &sdl.Rect{Y: int32(ui.winHeight / 2), W: 32, H: 32}); err != nil {
		panic(err)
	}

	// Drawing Hitpoints count
	tex = ui.stringToTexture("Life "+strconv.Itoa(level.Player.Hitpoints), sdl.Color{R: 255}, FontSmall)
	_, _, w, h, _ = tex.Query()
	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: 32, Y: int32(ui.winHeight / 2), W: w, H: h})
	if err != nil {
		panic(err)
	}

	// Life gauge using red rect on black rect
	if err := ui.renderer.Copy(ui.textureAtlas, &sdl.Rect{X: 928, Y: 1600, W: 32, H: 32}, &sdl.Rect{X: int32(level.Player.Pos.X*32) + ui.offsetX, Y: int32((level.Player.Pos.Y-1)*32) + ui.offsetY + 20, W: 32, H: 5}); err != nil {
		panic(err)
	}

	var gauge float64
	gauge = float64(level.Player.Hitpoints) / float64(level.Player.MaxHitpoints)

	if err := ui.renderer.Copy(ui.textureAtlas, &sdl.Rect{X: 1024, Y: 1600, W: 32, H: 32}, &sdl.Rect{X: int32(level.Player.Pos.X*32) + ui.offsetX, Y: int32((level.Player.Pos.Y-1)*32) + ui.offsetY + 20, W: int32(32 * gauge), H: 5}); err != nil {
		panic(err)
	}

}

func (ui *ui) getSinglePixel(color sdl.Color) *sdl.Texture {
	tex, err := ui.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, 1, 1)
	if err != nil {
		panic(err)
	}

	pixels := []byte{color.R, color.G, color.B, color.A}
	err = tex.Update(nil, unsafe.Pointer(&pixels[0]), 4)
	if err != nil {
		panic(err)
	}
	return tex
}

func (ui *ui) getGroundItemRect(i int) *sdl.Rect {
	return &sdl.Rect{X: int32(ui.winWidth - 32 - i*32), Y: 0, W: 32, H: 32}
}

func (ui *ui) pickupItem(level *game.Level, mouseX, mouseY int32) *game.Item {
	items := level.Items[level.Player.Pos]
	for i, item := range items {
		itemRect := ui.getGroundItemRect(i)
		if itemRect.HasIntersection(&sdl.Rect{X: mouseX, Y: mouseY, W: 1, H: 1}) {
			return item
		}
	}
	return nil
}

func (ui *ui) Run() {
	var newLevel *game.Level
	var ok bool
	for {
		input := game.Input{}
		select {
		case newLevel, ok = <-ui.levelChan:
			if ok {
				switch newLevel.LastEvent {
				case game.Move:
					playRandomSound(ui.sounds.footstep, soundsVolume)
					//TODO - improve animations
				case game.DoorOpen:
					playRandomSound(ui.sounds.openDoor, soundsVolume)
				case game.Attack:
					playRandomSound(ui.sounds.swing, soundsVolume)
				case game.Pickup:
					playRandomSound(ui.sounds.pickup, soundsVolume)
				default:
				}
				newLevel.LastEvent = game.Empty
				if ui.state == UIMain {
					ui.Draw(newLevel)
				} else {
					ui.Draw(newLevel)
					ui.DrawInventory(newLevel)
				}
			}
		default:
		}
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				ui.inputChan <- &game.Input{Typ: game.QuitGame}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					ui.inputChan <- &game.Input{Typ: game.CloseWindow, LevelChannel: ui.levelChan}
				}
			case *sdl.MouseButtonEvent:
				if e.State == sdl.PRESSED && e.Button == sdl.BUTTON_LEFT {
					item := ui.pickupItem(newLevel, e.X, e.Y)
					if item != nil {
						ui.inputChan <- &game.Input{Typ: game.TakeItem, Item: item}
					}
				}
			case *sdl.KeyboardEvent:
				// Seems not needed
				//if sdl.GetKeyboardFocus() == ui.window || sdl.GetMouseFocus() == ui.window {
				if e.State != sdl.PRESSED {
					break
				}
				switch e.Keysym.Sym {
				case sdl.K_UP:
					input = game.Input{Typ: game.Up}
				case sdl.K_DOWN:
					input = game.Input{Typ: game.Down}
				case sdl.K_LEFT:
					input = game.Input{Typ: game.Left}
				case sdl.K_RIGHT:
					input = game.Input{Typ: game.Right}
				case sdl.K_e:
					input = game.Input{Typ: game.Action}
				case sdl.K_t:
					input = game.Input{Typ: game.TakeAll}
				case sdl.K_i:
					if ui.state == UIMain {
						ui.state = UIInventory
					} else {
						ui.state = UIMain
					}
				default:
					input = game.Input{Typ: game.None}
				}
				if input.Typ != game.None {
					ui.inputChan <- &input
				}
				//}
			}
		}
		sdl.Delay(10)
	}
}
