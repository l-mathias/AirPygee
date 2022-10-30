package ui2d

//TODO - add damage on top of character when combat
//TODO - add Player character selection
//TODO - improve life gauge using shading
//TODO - Improve fog of war effect using transparent texture or special tiles
//TODO - Fix UI bug when dead, should recenter camera

import (
	"AirPygee/game"
	"bufio"
	"encoding/xml"
	"fmt"
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

type ui struct {
	state               uiState
	sounds              sounds
	winWidth, winHeight int
	renderer            *sdl.Renderer
	window              *sdl.Window
	textureAtlas        *sdl.Texture
	tileMap             *sdl.Texture

	//player
	pTexture                                          *sdl.Texture
	pWidthTex, pHeightTex                             int32
	pFromX, pFromY, pFramesX, pFramesY, pCurrentFrame int32
	pSrc                                              sdl.Rect
	pDest                                             sdl.Rect

	// UI Theme
	uipack       *sdl.Texture
	texturesList SubTextures

	textureIndex     map[rune][]sdl.Rect
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
	draggedItem       *game.Item
	dragMode          dragMode
	currentMouseState *mouseState
	prevMouseState    *mouseState

	// Fonts
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
	ui.uipack = ui.imgFileToTexture("ui2d/assets/uipack_rpg_sheet.png")

	ui.loadTextureIndex()
	ui.LoadPlayer()
	ui.loadSpritesheetFromXml()

	ui.centerX = -1
	ui.centerY = -1

	ui.fontSmall, err = ttf.OpenFont("ui2d/assets/Kingthings_Foundation.ttf", int(float64(ui.winWidth)*0.010))
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
	ui.sounds.closeDoor = buildSoundsVariations("ui2d/assets/audio/sounds/Kenney/doorClose*.ogg")
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
	openDoor  []*mix.Chunk
	closeDoor []*mix.Chunk
	footstep  []*mix.Chunk
	swing     []*mix.Chunk
	pickup    []*mix.Chunk
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
			rects = append(rects, sdl.Rect{X: int32(x) * tileSize, Y: int32(y) * tileSize, W: tileSize, H: tileSize})
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

func (ui *ui) LoadPlayer() {
	ui.pCurrentFrame = 0
	ui.pFramesX = 4
	ui.pFramesY = 4

	image, err := img.Load("ui2d/assets/george.png")
	if err != nil {
		panic(err)
	}
	defer image.Free()
	ui.pTexture, err = ui.renderer.CreateTextureFromSurface(image)
	if err != nil {
		panic(err)
	}

	_, _, imageWidth, imageHeight, _ := ui.pTexture.Query()
	ui.pWidthTex = imageWidth / ui.pFramesX
	ui.pHeightTex = imageHeight / ui.pFramesY
}

func (ui *ui) UpdatePlayer(input game.InputType) {
	ui.pCurrentFrame++
	if ui.pCurrentFrame >= ui.pFramesY {
		ui.pCurrentFrame = 0
	}
	switch input {
	case game.Up:
		ui.pFromX = 2 * ui.pWidthTex
	case game.Down:
		ui.pFromX = 0
	case game.Left:
		ui.pFromX = ui.pWidthTex
	case game.Right:
		ui.pFromX = 3 * ui.pWidthTex
	}
}

func (ui *ui) drawPlayer(level *game.Level) {
	p := level.Player
	ui.pFromY = ui.pCurrentFrame * ui.pWidthTex

	ui.pSrc = sdl.Rect{X: ui.pFromX, Y: ui.pFromY, W: ui.pWidthTex, H: ui.pHeightTex}
	ui.pDest = sdl.Rect{X: int32(p.X)*tileSize + ui.offsetX - ((ui.pWidthTex - tileSize) / 2), Y: int32(p.Y)*tileSize + ui.offsetY - ((ui.pHeightTex - tileSize) / 2), W: ui.pWidthTex, H: ui.pHeightTex}

	err := ui.renderer.Copy(ui.pTexture, &ui.pSrc, &ui.pDest)
	if err != nil {
		panic(err)
	}
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
	if err != nil {
		fmt.Println(err)
	}

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

func (ui *ui) drawInventory(level *game.Level) {
	var locationX, locationY int32

	playerSrcRect := sdl.Rect{X: 0, Y: 0, W: 48, H: 48}
	playerX := ((ui.invWidth - (ui.invWidth / 2)) / 2) + ui.invOffsetX
	playerY := ((ui.invHeight - (ui.invHeight / 2)) / 2) + ui.invOffsetY

	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("panel_beige.png"), &sdl.Rect{X: ui.invOffsetX, Y: ui.invOffsetY, W: ui.invWidth, H: ui.invHeight}); err != nil {
		panic(err)
	}

	if err := ui.renderer.Copy(ui.pTexture, &playerSrcRect, &sdl.Rect{X: playerX, Y: playerY, W: ui.invWidth / 2, H: ui.invHeight / 2}); err != nil {
		panic(err)
	}

	ui.drawEmptyInventory(level)

	// draw panel items
	//Head
	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("buttonSquare_blue_pressed.png"), &sdl.Rect{X: ui.invHeadX, Y: ui.invHeadY, W: ui.itemW, H: ui.itemH}); err != nil {
		panic(err)
	}
	//RightHand
	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("buttonSquare_blue_pressed.png"), &sdl.Rect{X: ui.invRHandX, Y: ui.invRHandY, W: ui.itemW, H: ui.itemH}); err != nil {
		panic(err)
	}
	//LeftHand
	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("buttonSquare_blue_pressed.png"), &sdl.Rect{X: ui.invLHandX, Y: ui.invLHandY, W: ui.itemW, H: ui.itemH}); err != nil {
		panic(err)
	}
	//Foots
	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("buttonSquare_blue_pressed.png"), &sdl.Rect{X: ui.invFootsX, Y: ui.invFootsY, W: ui.itemW, H: ui.itemH}); err != nil {
		panic(err)
	}
	//Chest
	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("buttonSquare_blue_pressed.png"), &sdl.Rect{X: ui.invChestX, Y: ui.invChestY, W: ui.itemW, H: ui.itemH}); err != nil {
		panic(err)
	}
	//Legs
	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("buttonSquare_blue_pressed.png"), &sdl.Rect{X: ui.invLegsX, Y: ui.invLegsY, W: ui.itemW, H: ui.itemH}); err != nil {
		panic(err)
	}

	for i, item := range level.Player.EquippedItems {
		itemSrcRect := ui.textureIndex[item.Rune][0]

		switch item.Location {
		case game.Head:
			locationX = ui.invHeadX
			locationY = ui.invHeadY
		case game.RightHand:
			locationX = ui.invRHandX
			locationY = ui.invRHandY
		case game.LeftHand:
			locationX = ui.invLHandX
			locationY = ui.invLHandY
		case game.Foots:
			locationX = ui.invFootsX
			locationY = ui.invFootsY
		case game.Chest:
			locationX = ui.invChestX
			locationY = ui.invChestY
		case game.Legs:
			locationX = ui.invLegsX
			locationY = ui.invLegsY
		default:
			locationX = ui.invOffsetX + int32(i)*ui.itemW
			locationY = ui.invOffsetY + ui.invHeight - ui.itemH
		}

		if item == ui.draggedItem {
			if err := ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{X: int32(ui.currentMouseState.pos.X) - ui.itemW/2, Y: int32(ui.currentMouseState.pos.Y) - ui.itemH/2, W: ui.itemW, H: ui.itemH}); err != nil {
				panic(err)
			}
		} else {
			if err := ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{X: locationX, Y: locationY, W: ui.itemW, H: ui.itemH}); err != nil {
				panic(err)
			}
		}
	}

	var countX int32 = 0
	var countY int32 = 0
	for _, item := range level.Player.Items {
		itemSrcRect := ui.textureIndex[item.Rune][0]
		if countX%5 == 0 {
			countX = 0
			countY++
		}
		locationX = ui.invOffsetX + ui.invWidth + ui.itemW*countX
		locationY = ui.invOffsetY + ui.itemH*countY

		if item == ui.draggedItem {
			if err := ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{X: int32(ui.currentMouseState.pos.X) - ui.itemW/2, Y: int32(ui.currentMouseState.pos.Y) - ui.itemH/2, W: ui.itemW, H: ui.itemH}); err != nil {
				panic(err)
			}
		} else {
			if err := ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{X: locationX, Y: locationY, W: ui.itemW, H: ui.itemH}); err != nil {
				panic(err)
			}
		}
		countX++
	}
}

func (ui *ui) drawEmptyInventory(level *game.Level) {
	var countX, countY, locationX, locationY int32
	for i := 0; i < level.Player.InventorySize; i++ {
		if i%5 == 0 {
			countX = 0
			countY++
		}
		locationX = ui.invOffsetX + ui.invWidth + ui.itemW*countX
		locationY = ui.invOffsetY + ui.itemH*countY
		if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("buttonSquare_brown_pressed.png"), &sdl.Rect{X: locationX, Y: locationY, W: ui.itemW, H: ui.itemH}); err != nil {
			panic(err)
		}
		countX++
	}
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
					dstRect := sdl.Rect{X: int32(x)*tileSize + ui.offsetX, Y: int32(y)*tileSize + ui.offsetY, W: tileSize, H: tileSize}
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
			err := ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{X: int32(ui.winWidth) - tileSize - int32(i)*tileSize, Y: 0, W: tileSize, H: tileSize})
			if err != nil {
				panic(err)
			}
		}
		// drawing help letter T
		if err := ui.renderer.Copy(ui.tileMap, &sdl.Rect{X: 358, Y: 34, W: 16, H: 16}, &sdl.Rect{X: int32(ui.winWidth) - tileSize - int32(len(groundItems))*tileSize, Y: 0, W: tileSize, H: tileSize}); err != nil {
			panic(err)
		}
	}

	if level.Map[level.FrontOf().Y][level.FrontOf().X].Actionable {
		// drawing help letter E
		if err := ui.renderer.Copy(ui.tileMap, &sdl.Rect{X: 324, Y: 34, W: 16, H: 16}, &sdl.Rect{X: int32(ui.winWidth) - tileSize - tileSize, Y: 0, W: tileSize, H: tileSize}); err != nil {
			panic(err)
		}

	}
	//ui.renderer.Present()
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
	return &sdl.Rect{X: int32(ui.winWidth) - tileSize - int32(i)*tileSize, Y: 0, W: tileSize, H: tileSize}
}

func (ui *ui) getInventoryItemRect(id int, level *game.Level) *sdl.Rect {
	var locationX, locationY, countX, countY int32

	for i := 0; i < level.Player.InventorySize; i++ {
		if i%5 == 0 {
			countX = 0
			countY++
		}
		locationX = ui.invOffsetX + ui.invWidth + ui.itemW*countX
		locationY = ui.invOffsetY + ui.itemH*countY

		countX++
		if i == id {
			return &sdl.Rect{X: locationX, Y: locationY, W: ui.itemW, H: ui.itemH}
		}
	}
	return nil
}

// getEquippedItemRect based on arbitraries items positions, will return the corresponding
// rectangle in order to compare with click position and then unequip
func (ui *ui) getEquippedItemRect(item *game.Item) *sdl.Rect {
	var locationX, locationY int32

	switch item.Location {
	case game.Head:
		locationX = ui.invHeadX
		locationY = ui.invHeadY
	case game.Foots:
		locationX = ui.invFootsX
		locationY = ui.invFootsY
	case game.LeftHand:
		locationX = ui.invLHandX
		locationY = ui.invLHandY
	case game.RightHand:
		locationX = ui.invRHandX
		locationY = ui.invRHandY
	case game.Chest:
		locationX = ui.invChestX
		locationY = ui.invChestY
	case game.Legs:
		locationX = ui.invLegsX
		locationY = ui.invLegsY
	}

	return &sdl.Rect{X: locationX, Y: locationY, W: ui.itemW, H: ui.itemH}
}

func (ui *ui) clickValidItem(level *game.Level, mouseX, mouseY int32) *game.Item {
	for i, item := range level.Player.Items {
		itemRect := ui.getInventoryItemRect(i, level)
		if itemRect.HasIntersection(&sdl.Rect{X: mouseX, Y: mouseY, W: 1, H: 1}) {
			ui.dragMode = fromInventory
			return item
		}
	}

	for _, item := range level.Player.EquippedItems {
		itemRect := ui.getEquippedItemRect(item)
		if itemRect.HasIntersection(&sdl.Rect{X: mouseX, Y: mouseY, W: 1, H: 1}) {
			ui.dragMode = fromEquippedItems
			return item
		}
	}

	return nil
}

func (ui *ui) isSlotFree(level *game.Level, itemToEquip *game.Item) bool {
	for _, item := range level.Player.EquippedItems {
		if itemToEquip.Location == item.Location {
			return false
		}
	}
	return true
}

func (ui *ui) hasClickedOnValidEquipSlot(mouseX, mouseY int32, item *game.Item) bool {
	itemRect := ui.getEquippedItemRect(item)
	if itemRect.HasIntersection(&sdl.Rect{X: mouseX, Y: mouseY, W: 1, H: 1}) {
		return true
	} else {
		return false
	}
}

func (ui *ui) hasClickedInBackpackZone(mouseX, mouseY int32) bool {
	itemRect := &sdl.Rect{X: ui.invOffsetX + ui.invWidth, Y: ui.invOffsetY + ui.itemH, W: ui.itemW * 5, H: ui.itemH * 4}
	if itemRect.HasIntersection(&sdl.Rect{X: mouseX, Y: mouseY, W: 1, H: 1}) {
		return true
	} else {
		return false
	}
}

func (ui *ui) hasClickedInEquipZone(mouseX, mouseY int32) bool {
	itemRect := &sdl.Rect{X: ui.invOffsetX, Y: ui.invOffsetY, W: ui.invWidth, H: ui.invHeight}
	if itemRect.HasIntersection(&sdl.Rect{X: mouseX, Y: mouseY, W: 1, H: 1}) {
		return true
	} else {
		return false
	}
}

func (ui *ui) hasClickedOutsideInventoryZone(mouseX, mouseY int32) bool {
	return !ui.hasClickedInBackpackZone(mouseX, mouseY) && !ui.hasClickedInEquipZone(mouseX, mouseY)
}

// pickupGroundItem will check if clicked on a ground item, then take it
func (ui *ui) pickupGroundItem(level *game.Level, mouseX, mouseY int32) *game.Item {
	items := level.Items[level.Player.Pos]
	for i, item := range items {
		itemRect := ui.getGroundItemRect(i)
		if itemRect.HasIntersection(&sdl.Rect{X: mouseX, Y: mouseY, W: 1, H: 1}) {
			return item
		}
	}
	return nil
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
					playRandomSound(ui.sounds.footstep, soundsVolume)
					//TODO - improve animations
				case game.DoorOpen:
					playRandomSound(ui.sounds.openDoor, soundsVolume)
				case game.DoorClose:
					playRandomSound(ui.sounds.closeDoor, soundsVolume)
				case game.Attack:
					playRandomSound(ui.sounds.swing, soundsVolume)
				case game.Pickup:
					playRandomSound(ui.sounds.pickup, soundsVolume)
				default:
				}
				newLevel.LastEvent = game.Empty
				if ui.state == UIMain {
					ui.draw(newLevel)
				} else {
					ui.draw(newLevel)
					ui.drawInventory(newLevel)
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
			case *sdl.MouseMotionEvent:
				if ui.draggedItem != nil && ui.state == UIInventory {
					ui.draw(newLevel)
					ui.drawInventory(newLevel)
				}
			case *sdl.MouseButtonEvent:
				if e.State == sdl.PRESSED && e.Button == sdl.BUTTON_LEFT {
					if ui.state == UIInventory && ui.draggedItem == nil {
						ui.draggedItem = ui.clickValidItem(newLevel, e.X, e.Y)
						if ui.draggedItem == nil {
							ui.draggedItem = ui.clickValidItem(newLevel, e.X, e.Y)
						}
						ui.draw(newLevel)
						ui.drawInventory(newLevel)
					}
				}
				if e.State == sdl.RELEASED && e.Button == sdl.BUTTON_LEFT {
					if ui.state == UIMain {
						item := ui.pickupGroundItem(newLevel, e.X, e.Y)
						if item != nil {
							ui.inputChan <- &game.Input{Typ: game.TakeItem, Item: item}
						}
					} else {
						if ui.draggedItem != nil && ui.dragMode != none {
							var item *game.Item
							if ui.dragMode == fromInventory {
								if ui.hasClickedOnValidEquipSlot(e.X, e.Y, ui.draggedItem) && ui.isSlotFree(newLevel, ui.draggedItem) {
									item = ui.draggedItem
								} else if ui.hasClickedOutsideInventoryZone(e.X, e.Y) {
									ui.inputChan <- &game.Input{Typ: game.Drop, Item: ui.draggedItem}
								}
							} else if ui.hasClickedInBackpackZone(e.X, e.Y) {
								item = ui.draggedItem
							}

							if item != nil {
								ui.inputChan <- &game.Input{Typ: game.Equip, Item: ui.draggedItem}
							}
							ui.draggedItem = nil
							ui.dragMode = none
							ui.draw(newLevel)
							ui.drawInventory(newLevel)
						}
					}
				}
			case *sdl.KeyboardEvent:
				if e.State != sdl.PRESSED {
					break
				}
				switch e.Keysym.Sym {
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
						ui.draw(newLevel)
						ui.drawInventory(newLevel)
					} else {
						ui.state = UIMain
						ui.draw(newLevel)
					}
				default:
					input = game.Input{Typ: game.None}
				}
				if input.Typ != game.None {
					ui.inputChan <- &input
				}
			}
		}
		ui.renderer.Present()
		ui.prevMouseState = ui.currentMouseState
		sdl.Delay(1)
	}
}
