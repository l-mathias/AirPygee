package ui2d

// TODO - add damage on top of character when combat
// TODO - add Player character selection
// TODO - Improve fog of war effect using transparent texture or special tiles

import (
	"AirPygee/game"
	"bufio"
	"encoding/xml"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

const (
	FontSmall FontSize = iota
	FontMedium
	FontLarge
)

type uiState int

const (
	UIMain uiState = iota
	UIInventory
	UIMenu
	itemSizeRatio float64 = 0.15
	tileSize      int32   = 32
)

type dragMode int

const (
	none dragMode = iota
	fromInventory
	fromEquippedItems
)

type mouseState struct {
	leftButton  bool
	rightButton bool
	pos         game.Pos
	xrel, yrel  int32
}

func getMouseState() *mouseState {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()
	var result mouseState
	result.pos = game.Pos{int(mouseX), int(mouseY)}
	result.leftButton = !(leftButton == 0)
	result.rightButton = !(rightButton == 0)

	return &result
}

type menuButton struct {
	name           string
	buttonRect     *sdl.Rect
	buttonTexture  *sdl.Texture
	buttonTextRect *sdl.Rect
	highlighted    bool
}

type coloredFont struct {
	text  string
	color sdl.Color
}

type TextureIndex struct {
	mu    sync.RWMutex
	rects map[rune][]sdl.Rect
}

type ui struct {
	state               uiState
	sounds              sounds
	winWidth, winHeight int
	renderer            *sdl.Renderer
	window              *sdl.Window
	textureAtlas        *sdl.Texture
	tileMap             *sdl.Texture

	// Sounds & Music
	musicVolume  int
	soundsVolume int

	//player
	pTexture                                          *sdl.Texture
	pWidthTex, pHeightTex                             int32
	pFromX, pFromY, pFramesX, pFramesY, pCurrentFrame int32
	pSrc                                              sdl.Rect
	pDest                                             sdl.Rect

	// UI Theme
	uipack       *sdl.Texture
	texturesList SubTextures

	textureIndexTiles, textureIndexMonsters, textureIndexItems TextureIndex
	animations                                                 map[rune][]*sdl.Rect

	centerX, centerY int
	r                *rand.Rand
	levelChan        chan *game.Level
	inputChan        chan *game.Input
	offsetX, offsetY int32

	// Inventory
	invOffsetX, invOffsetY, invWidth, invHeight, itemW, itemH      int32
	invHeadX, invHeadY, invLHandX, invLHandY, invRHandX, invRHandY int32
	invLegsX, invLegsY, invChestX, invChestY, invFootsX, invFootsY int32

	// drag&drop
	draggedItem       game.Item
	dragMode          dragMode
	currentMouseState *mouseState
	prevMouseState    *mouseState

	// Fonts
	fontSmall, fontMedium, fontLarge          *ttf.Font
	str2TexSmall, str2TexMedium, str2TexLarge map[coloredFont]*sdl.Texture
	//Main Menu
	menuButtons []*menuButton
}

func NewUI(inputChan chan *game.Input, levelChan chan *game.Level) *ui {
	ui := &ui{}
	ui.state = UIMain
	ui.inputChan = inputChan
	ui.levelChan = levelChan
	ui.str2TexSmall = make(map[coloredFont]*sdl.Texture)
	ui.str2TexMedium = make(map[coloredFont]*sdl.Texture)
	ui.str2TexLarge = make(map[coloredFont]*sdl.Texture)
	ui.animations = make(map[rune][]*sdl.Rect)
	ui.r = rand.New(rand.NewSource(1))
	ui.winWidth = 1280
	ui.winHeight = 720
	window, err := sdl.CreateWindow("AirPygee", 100, 100, int32(ui.winWidth), int32(ui.winHeight), sdl.WINDOW_SHOWN)
	game.CheckError(err)

	ui.window = window

	ui.renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	game.CheckError(err)

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")

	ui.textureAtlas = ui.imgFileToTexture("ui2d/assets/tiles.png")
	ui.tileMap = ui.imgFileToTexture("ui2d/assets/tilemap.png")
	ui.uipack = ui.imgFileToTexture("ui2d/assets/uipack_rpg_sheet.png")

	ui.loadTextureIndex(&ui.textureIndexTiles, "ui2d/assets/atlas-index.txt")
	ui.loadTextureIndex(&ui.textureIndexMonsters, "ui2d/assets/atlas-index-monsters.txt")
	ui.loadTextureIndex(&ui.textureIndexItems, "ui2d/assets/atlas-index-items.txt")
	ui.LoadPlayer()
	ui.loadSpritesheetFromXml()

	ui.centerX = -1
	ui.centerY = -1

	ui.fontSmall, err = ttf.OpenFont("ui2d/assets/Kingthings_Foundation.ttf", int(float64(ui.winWidth)*0.010))
	game.CheckError(err)

	ui.fontMedium, err = ttf.OpenFont("ui2d/assets/Kingthings_Foundation.ttf", 24)
	game.CheckError(err)

	ui.fontLarge, err = ttf.OpenFont("ui2d/assets/Kingthings_Foundation.ttf", 32)
	game.CheckError(err)

	ui.musicVolume = 32
	ui.soundsVolume = 10
	ui.loadSounds()

	ui.invWidth = int32(float64(ui.winWidth) * 0.40)
	ui.invHeight = int32(float64(ui.winHeight) * 0.75)
	ui.invOffsetX = (int32(ui.winWidth) - ui.invWidth) / 2
	ui.invOffsetY = (int32(ui.winHeight) - ui.invHeight) / 2
	ui.itemW = int32(float64(ui.invWidth) * itemSizeRatio)
	ui.itemH = int32(float64(ui.invHeight) * itemSizeRatio)

	ui.invHeadX = ui.invOffsetX + ui.invWidth/2 - (ui.itemW / 2)
	ui.invHeadY = ui.invOffsetY + ui.itemH

	ui.invRHandX = ui.invOffsetX + ui.itemW
	ui.invRHandY = ui.invOffsetY + ui.itemH*4

	ui.invLHandX = ui.invOffsetX + ui.invWidth - ui.itemW*2
	ui.invLHandY = ui.invOffsetY + ui.itemH*4

	ui.invFootsX = ui.invOffsetX + ui.invWidth/2 - ui.itemW
	ui.invFootsY = ui.invOffsetY + ui.invHeight - ui.itemH*2

	ui.invChestX = ui.invOffsetX + ui.itemW
	ui.invChestY = ui.invOffsetY + ui.itemH*3

	ui.invLegsX = ui.invOffsetX + ui.invWidth/2
	ui.invLegsY = ui.invOffsetY + ui.invHeight - ui.itemH*2

	ui.currentMouseState = &mouseState{
		leftButton:  false,
		rightButton: false,
		pos:         game.Pos{},
	}

	ui.prevMouseState = &mouseState{
		leftButton:  false,
		rightButton: false,
		pos:         game.Pos{},
	}

	ui.dragMode = none

	ui.menuButtons = make([]*menuButton, 0)
	ui.buildMenuButtons()
	ui.buildAnimations()

	return ui
}

func (ui *ui) buildAnimations() {
	ui.buildAnimation(game.AnimatedPortal, &sdl.Rect{X: 1376, Y: 320, W: tileSize, H: tileSize}, &sdl.Rect{X: 1408, Y: 320, W: tileSize, H: tileSize}, &sdl.Rect{X: 1440, Y: 320, W: tileSize, H: tileSize})

	ui.buildAnimation(game.UpAnim, &sdl.Rect{X: 352, Y: 768, W: tileSize, H: tileSize})
	ui.buildAnimation(game.DownAnim, &sdl.Rect{X: 480, Y: 768, W: tileSize, H: tileSize})
	ui.buildAnimation(game.LeftAnim, &sdl.Rect{X: 544, Y: 768, W: tileSize, H: tileSize})
	ui.buildAnimation(game.RightAnim, &sdl.Rect{X: 416, Y: 768, W: tileSize, H: tileSize})
}

func (ui *ui) loadSounds() {
	err := mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096)
	game.CheckError(err)

	mus, err := mix.LoadMUS("ui2d/assets/audio/music/the_field_of_dreams.mp3")
	game.CheckError(err)
	mix.VolumeMusic(ui.musicVolume)

	err = mus.Play(-1)
	game.CheckError(err)

	ui.sounds.footstep = buildSoundsVariations("ui2d/assets/audio/sounds/Kenney/footstep*.ogg")
	ui.sounds.openDoor = buildSoundsVariations("ui2d/assets/audio/sounds/Kenney/doorOpen*.ogg")
	ui.sounds.closeDoor = buildSoundsVariations("ui2d/assets/audio/sounds/Kenney/doorClose*.ogg")
	ui.sounds.swing = buildSoundsVariations("ui2d/assets/audio/sounds/battle/swing*.wav")
	ui.sounds.pickup = buildSoundsVariations("ui2d/assets/audio/sounds/Kenney/cloth*.ogg")
	ui.sounds.potion = buildSoundsVariations("ui2d/assets/audio/sounds/Kenney/bubble*.wav")

}

func buildSoundsVariations(pattern string) []*mix.Chunk {
	fileNames, err := filepath.Glob(pattern)
	game.CheckError(err)
	result := make([]*mix.Chunk, 0)

	for _, fileName := range fileNames {
		sound, err := mix.LoadWAV(fileName)
		game.CheckError(err)

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
	openDoor  []*mix.Chunk
	closeDoor []*mix.Chunk
	footstep  []*mix.Chunk
	swing     []*mix.Chunk
	pickup    []*mix.Chunk
	potion    []*mix.Chunk
}

type FontSize int

func (ui *ui) stringToTexture(s string, color sdl.Color, size FontSize) *sdl.Texture {

	coloredFont := coloredFont{s, color}
	var font *ttf.Font
	switch size {
	case FontSmall:
		font = ui.fontSmall
		if tex, exists := ui.str2TexSmall[coloredFont]; exists {
			return tex
		}
	case FontMedium:
		font = ui.fontMedium
		if tex, exists := ui.str2TexMedium[coloredFont]; exists {
			return tex
		}
	case FontLarge:
		font = ui.fontLarge
		if tex, exists := ui.str2TexLarge[coloredFont]; exists {
			return tex
		}
	}
	fontSurface, err := font.RenderUTF8Blended(s, color)
	if err != nil {
		panic(err)
	}
	defer fontSurface.Free()

	tex, err := ui.renderer.CreateTextureFromSurface(fontSurface)
	game.CheckError(err)

	switch size {
	case FontSmall:
		ui.str2TexSmall[coloredFont] = tex
	case FontMedium:
		ui.str2TexMedium[coloredFont] = tex
	case FontLarge:
		ui.str2TexLarge[coloredFont] = tex
	}

	return tex
}

func (ui *ui) loadTextureIndex(textureIndex *TextureIndex, fileName string) {
	//ui.textureIndexTiles.rects = make(map[rune][]sdl.Rect)
	textureIndex.rects = make(map[rune][]sdl.Rect)

	//infile, err := os.Open("ui2d/assets/atlas-index.txt")
	infile, err := os.Open(fileName)
	game.CheckError(err)

	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		tileRune := rune(line[0])
		xy := line[1:]
		splitXyC := strings.Split(xy, ",")
		x, err := strconv.ParseInt(strings.TrimSpace(splitXyC[0]), 10, 64)
		game.CheckError(err)

		y, err := strconv.ParseInt(strings.TrimSpace(splitXyC[1]), 10, 64)
		game.CheckError(err)

		variationCount, err := strconv.ParseInt(strings.TrimSpace(splitXyC[2]), 10, 64)
		game.CheckError(err)

		var rects []sdl.Rect
		for i := int64(0); i < variationCount; i++ {
			rects = append(rects, sdl.Rect{X: int32(x) * tileSize, Y: int32(y) * tileSize, W: tileSize, H: tileSize})
			x++
			if x > 62 {
				x = 0
				y++
			}
		}
		//ui.textureIndexTiles.rects[tileRune] = rects
		textureIndex.rects[tileRune] = rects
	}

}

func (ui *ui) imgFileToTexture(filename string) *sdl.Texture {
	infile, err := os.Open(filename)
	game.CheckError(err)
	defer infile.Close()

	image, err := img.Load(filename)
	if err != nil {
		panic(err)
	}

	tex, err := ui.renderer.CreateTextureFromSurface(image)
	game.CheckError(err)
	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	game.CheckError(err)
	return tex
}

func init() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	game.CheckError(err)

	err = ttf.Init()
	game.CheckError(err)

	err = mix.Init(mix.INIT_MP3)
	game.CheckError(err)

	err = mix.Init(mix.INIT_OGG)
	game.CheckError(err)

}

type SubTextures struct {
	XMLName    xml.Name     `xml:"TextureAtlas"`
	SubTexture []SubTexture `xml:"SubTexture"`
}

type SubTexture struct {
	XMLName xml.Name `xml:"SubTexture"`
	Name    string   `xml:"name,attr"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
}

func (ui *ui) loadSpritesheetFromXml() {
	xmlFile, err := os.Open("ui2d/assets/uipack_rpg_sheet.xml")
	game.CheckError(err)

	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	xml.Unmarshal(byteValue, &ui.texturesList)
}

func (ui *ui) getRectFromTextureName(name string) *sdl.Rect {
	for i := 0; i < len(ui.texturesList.SubTexture); i++ {
		if ui.texturesList.SubTexture[i].Name == name {
			x, _ := strconv.Atoi(ui.texturesList.SubTexture[i].X)
			y, _ := strconv.Atoi(ui.texturesList.SubTexture[i].Y)
			w, _ := strconv.Atoi(ui.texturesList.SubTexture[i].Width)
			h, _ := strconv.Atoi(ui.texturesList.SubTexture[i].Height)
			return &sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}
		}
	}
	return &sdl.Rect{}
}

func (ui *ui) draw(level *game.Level) {
	if ui.centerX == -1 && ui.centerY == -1 {
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
	ui.offsetX = int32(ui.winWidth/2) - int32(ui.centerX)*tileSize
	ui.offsetY = int32(ui.winHeight/2) - int32(ui.centerY)*tileSize

	err := ui.renderer.Clear()
	game.CheckError(err)
	ui.r.Seed(1)
	for y, row := range level.Map {
		for x, tile := range row {
			if tile.Rune != game.Blank {
				ui.textureIndexTiles.mu.RLock()
				srcRects := ui.textureIndexTiles.rects[tile.Rune]
				ui.textureIndexTiles.mu.RUnlock()
				srcRect := srcRects[ui.r.Intn(len(srcRects))]
				if tile.Visible || tile.Seen {
					dstRect := sdl.Rect{X: int32(x)*tileSize + ui.offsetX, Y: int32(y)*tileSize + ui.offsetY, W: tileSize, H: tileSize}
					pos := game.Pos{X: x, Y: y}
					if level.Debug[pos] {
						err = ui.textureAtlas.SetColorMod(128, 0, 0)
						game.CheckError(err)
					} else if tile.Seen && !tile.Visible {
						err = ui.textureAtlas.SetColorMod(128, 128, 128)
						game.CheckError(err)
					} else {
						err = ui.textureAtlas.SetColorMod(255, 255, 255)
						game.CheckError(err)
					}

					err = ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)
					game.CheckError(err)

					if tile.OverlayRune != game.Blank {
						// TODO - if multiple variants of a tile, adapt srcRects
						ui.textureIndexTiles.mu.RLock()
						srcRects = ui.textureIndexTiles.rects[tile.OverlayRune]
						ui.textureIndexTiles.mu.RUnlock()
						srcRect = srcRects[0]

						err = ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)
						game.CheckError(err)
					}
				}
			}
		}
	}

	ui.displayMonsters(level)
	ui.displayItems(level)
	ui.drawPlayer(level)
	ui.displayHUD(level)
	ui.displayStats(level)
	ui.displayEvents(level)

	// display item we are on top of
	if len(level.Items[level.Player.Pos]) > 0 {
		groundItems := level.Items[level.Player.Pos]
		for i, item := range groundItems {
			ui.textureIndexItems.mu.RLock()
			itemSrcRect := ui.textureIndexItems.rects[item.GetRune()][0]
			ui.textureIndexItems.mu.RUnlock()

			err = ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{X: int32(ui.winWidth) - tileSize - int32(i)*tileSize, Y: 0, W: tileSize, H: tileSize})
			game.CheckError(err)

		}
		// drawing help letter T
		err = ui.renderer.Copy(ui.tileMap, &sdl.Rect{X: 358, Y: 34, W: 16, H: 16}, &sdl.Rect{X: int32(ui.winWidth) - tileSize - int32(len(groundItems))*tileSize, Y: 0, W: tileSize, H: tileSize})
		game.CheckError(err)

	}

	if level.Map[level.FrontOf().Y][level.FrontOf().X].Actionable {
		// drawing help letter E
		err = ui.renderer.Copy(ui.tileMap, &sdl.Rect{X: 324, Y: 34, W: 16, H: 16}, &sdl.Rect{X: int32(ui.winWidth) - tileSize - tileSize, Y: 0, W: tileSize, H: tileSize})
		game.CheckError(err)

	}
	//ui.renderer.Present()
}

func (ui *ui) getSinglePixel(color sdl.Color) *sdl.Texture {
	tex, err := ui.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, 1, 1)
	game.CheckError(err)

	pixels := []byte{color.R, color.G, color.B, color.A}
	err = tex.Update(nil, unsafe.Pointer(&pixels[0]), 4)
	game.CheckError(err)
	return tex
}

func (ui *ui) getGroundItemRect(i int) *sdl.Rect {
	return &sdl.Rect{X: int32(ui.winWidth) - tileSize - int32(i)*tileSize, Y: 0, W: tileSize, H: tileSize}
}

// pickupGroundItem will check if clicked on a ground item, then take it
func (ui *ui) pickupGroundItem(level *game.Level, mouseX, mouseY int32) game.Item {
	items := level.Items[level.Player.Pos]
	for i, item := range items {
		itemRect := ui.getGroundItemRect(i)
		if itemRect.HasIntersection(&sdl.Rect{X: mouseX, Y: mouseY, W: 1, H: 1}) {
			return item
		}
	}
	return nil
}

func (ui *ui) fire(level *game.Level, attackRange int) {
	var direction rune
	var deltaX, deltaY int
	firstPos := level.FrontOf()
	positions := make([]game.Pos, 0)

	switch {
	case firstPos.X > level.Player.X:
		direction = game.RightAnim
		deltaX = 1
	case firstPos.X < level.Player.X:
		direction = game.LeftAnim
		deltaX = -1
	case firstPos.Y > level.Player.Y:
		direction = game.DownAnim
		deltaY = 1
	case firstPos.Y < level.Player.Y:
		direction = game.UpAnim
		deltaY = -1
	}

	for i := 1; i <= attackRange; i++ {
		x := level.Player.X + deltaX*i
		y := level.Player.Y + deltaY*i

		if x < 0 || !level.Map[y][x].Walkable {
			break
		}
		if y < 0 || !level.Map[y][x].Walkable {
			break
		}

		positions = append(positions, game.Pos{X: x, Y: y})
	}
	go ui.displayMovingAnimation(level, 500*time.Millisecond, direction, positions, &ui.textureIndexTiles)
}

// Run main UI loop
func (ui *ui) Run() {
	var newLevel *game.Level
	var ok bool
	ui.prevMouseState = getMouseState()

	for {
		ui.currentMouseState = getMouseState()
		input := game.Input{}
		select {
		case newLevel, ok = <-ui.levelChan:
			if ok {
				switch newLevel.LastEvent {
				case game.Move:
					playRandomSound(ui.sounds.footstep, ui.soundsVolume)
				case game.DoorOpen:
					playRandomSound(ui.sounds.openDoor, ui.soundsVolume)
				case game.DoorClose:
					playRandomSound(ui.sounds.closeDoor, ui.soundsVolume)
				case game.Attack:
					playRandomSound(ui.sounds.swing, ui.soundsVolume)
				case game.Pickup:
					playRandomSound(ui.sounds.pickup, ui.soundsVolume)
				case game.ConsumePotion:
					playRandomSound(ui.sounds.potion, ui.soundsVolume)
				default:
				}
				newLevel.LastEvent = game.Empty
				if ui.state == UIMain {
					ui.draw(newLevel)
				} else if ui.state == UIInventory {
					ui.draw(newLevel)
					ui.drawInventory(newLevel)
				}
			}
		default:
		}
		ui.draw(newLevel)

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				ui.inputChan <- &game.Input{Typ: game.QuitGame}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					ui.inputChan <- &game.Input{Typ: game.CloseWindow, LevelChannel: ui.levelChan}
				}
			case *sdl.MouseButtonEvent:
				if e.State == sdl.RELEASED && e.Button == sdl.BUTTON_LEFT {
					//if clicked on ground item zone
					item := ui.pickupGroundItem(newLevel, e.X, e.Y)
					if item != nil {
						ui.inputChan <- &game.Input{Typ: game.TakeItem, Item: item}
					}
				}
			case *sdl.KeyboardEvent:
				if e.State != sdl.PRESSED {
					break
				}
				switch e.Keysym.Sym {
				case sdl.K_a:
					ui.fire(newLevel, 3)
				case sdl.K_ESCAPE:
					if ui.state == UIMain {
						ui.state = UIMenu
						ui.menuActions()
					}
					ui.state = UIMain
				case sdl.K_UP:
					input = game.Input{Typ: game.Up}
					ui.UpdatePlayer(game.Up)
				case sdl.K_DOWN:
					input = game.Input{Typ: game.Down}
					ui.UpdatePlayer(game.Down)
				case sdl.K_LEFT:
					input = game.Input{Typ: game.Left}
					ui.UpdatePlayer(game.Left)
				case sdl.K_RIGHT:
					input = game.Input{Typ: game.Right}
					ui.UpdatePlayer(game.Right)
				case sdl.K_e:
					input = game.Input{Typ: game.Action}
				case sdl.K_t:
					input = game.Input{Typ: game.TakeAll}
				case sdl.K_i:
					if ui.state == UIMain {
						ui.state = UIInventory
						ui.menuInventory(newLevel)
					}
					ui.state = UIMain
				default:
					input = game.Input{Typ: game.None}
				}
				if input.Typ != game.None {
					ui.inputChan <- &input
				}
			}
		}

		ui.renderer.Present()
		sdl.Delay(1)
		ui.prevMouseState = ui.currentMouseState
	}
}
