package ui2d

import (
	"AirPygee/game"
	"github.com/veandco/go-sdl2/sdl"
	"strconv"
)

// displayStats draws general UI with remaining hit points and game instructions
func (ui *ui) displayStats(level *game.Level) {
	firstFrameX := 512
	for i := 0; i < 4; i++ {
		if err := ui.renderer.Copy(ui.tileMap, &sdl.Rect{X: int32(firstFrameX + i*16), Y: 68, W: 16, H: 16}, &sdl.Rect{X: int32(i * 32), Y: 0, W: tileSize, H: tileSize}); err != nil {
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
	if err := ui.renderer.Copy(ui.textureAtlas, &sdl.Rect{X: tileSize, Y: 0, W: tileSize, H: tileSize}, &sdl.Rect{Y: int32(ui.winHeight / 2), W: tileSize, H: tileSize}); err != nil {
		panic(err)
	}

	// Drawing Hit points count
	tex = ui.stringToTexture("Life "+strconv.Itoa(level.Player.Health), sdl.Color{R: 255}, FontMedium)
	_, _, w, h, _ = tex.Query()
	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: tileSize, Y: int32(ui.winHeight / 2), W: w, H: h})
	if err != nil {
		panic(err)
	}

	// Life gauge using red rect on black rect
	if err := ui.renderer.Copy(ui.textureAtlas, &sdl.Rect{X: 928, Y: 1600, W: tileSize, H: tileSize}, &sdl.Rect{X: int32(level.Player.Pos.X)*tileSize + ui.offsetX, Y: int32(level.Player.Pos.Y-1)*tileSize + ui.offsetY + 20, W: tileSize, H: 5}); err != nil {
		panic(err)
	}

	var gauge float64
	gauge = float64(level.Player.Health) / float64(level.Player.MaxHealth)

	if err := ui.renderer.Copy(ui.textureAtlas, &sdl.Rect{X: 1024, Y: 1600, W: tileSize, H: tileSize}, &sdl.Rect{X: int32(level.Player.Pos.X)*tileSize + ui.offsetX, Y: int32(level.Player.Pos.Y-1)*tileSize + ui.offsetY + 20, W: int32(float64(tileSize) * gauge), H: 5}); err != nil {
		panic(err)
	}

}

// displayMonsters displays monsters on map
func (ui *ui) displayMonsters(level *game.Level) {
	if err := ui.textureAtlas.SetColorMod(255, 255, 255); err != nil {
		panic(err)
	}
	for pos, monster := range level.Monsters {
		if level.Map[pos.Y][pos.X].Visible {

			if err := ui.renderer.Copy(ui.textureAtlas, &sdl.Rect{X: 928, Y: 1600, W: tileSize, H: tileSize}, &sdl.Rect{X: int32(level.Monsters[pos].X)*tileSize + ui.offsetX, Y: int32((level.Monsters[pos].Y-1))*tileSize + ui.offsetY + 20, W: tileSize, H: 5}); err != nil {
				panic(err)
			}
			var gauge float64
			gauge = float64(level.Monsters[pos].Health) / float64(level.Monsters[pos].MaxHealth)

			if err := ui.renderer.Copy(ui.textureAtlas, &sdl.Rect{X: 1024, Y: 1600, W: tileSize, H: tileSize}, &sdl.Rect{X: int32(level.Monsters[pos].X)*tileSize + ui.offsetX, Y: int32((level.Monsters[pos].Y-1))*tileSize + ui.offsetY + 20, W: tileSize * int32(gauge), H: 5}); err != nil {
				panic(err)
			}

			monsterSrcRect := ui.textureIndex[monster.Rune][0]
			err := ui.renderer.Copy(ui.textureAtlas, &monsterSrcRect, &sdl.Rect{X: int32(pos.X)*tileSize + ui.offsetX, Y: int32(pos.Y)*tileSize + ui.offsetY, W: tileSize, H: tileSize})
			if err != nil {
				panic(err)
			}
		}
	}
}

// displayItems displays items on Map
func (ui *ui) displayItems(level *game.Level) {
	for pos, items := range level.Items {
		if level.Map[pos.Y][pos.X].Visible {
			for _, item := range items {
				itemSrcRect := ui.textureIndex[item.Rune][0]
				var size int32
				size = tileSize
				if item.Name == "Potion" {
					switch item.Size {
					case "Small":
						size = int32(float64(size) * .50)
					case "Medium":
						size = int32(float64(size) * .75)
					case "Large":
						size = int32(float64(size) * .95)
					}
				}
				err := ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{X: int32(pos.X)*tileSize + ui.offsetX, Y: int32(pos.Y)*tileSize + ui.offsetY, W: size, H: size})
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

// displayEvents for drawing event list during game
func (ui *ui) displayEvents(level *game.Level) {
	textStartX := int32(float64(ui.winWidth) * .015)
	textStartY := int32(float64(ui.winHeight) * .68)
	textWidth := int32(float64(ui.winWidth) * .25)
	_, fontSizeY, _ := ui.fontSmall.SizeUTF8("A")

	err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("panel_beige.png"), &sdl.Rect{X: 0, Y: int32(ui.winHeight) - (int32(ui.winHeight) - textStartY + int32(fontSizeY)), W: textWidth, H: int32(ui.winHeight) - textStartY + int32(fontSizeY)})
	if err != nil {
		panic(err)
	}

	i := level.EventPos
	count := 0

	for {
		event := level.Events[i]
		if event != "" {
			tex := ui.stringToTexture(event, sdl.Color{R: 100, G: 50}, FontSmall)
			_, _, w, h, _ := tex.Query()

			err := ui.renderer.Copy(tex, nil, &sdl.Rect{X: textStartX, Y: int32(count*fontSizeY) + (int32(ui.winHeight) - (int32(ui.winHeight) - textStartY)), W: w, H: h})
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
}

func (ui *ui) highlightPrevious() {
	for i, b := range ui.menuButtons {
		if b.highlighted {
			b.highlighted = false
			if i-1 < 0 {
				ui.menuButtons[len(ui.menuButtons)-1].highlighted = true
				return
			} else {
				ui.menuButtons[i-1].highlighted = true
				return
			}
		}
	}
}

func (ui *ui) highlightNext() {
	for i, b := range ui.menuButtons {
		if b.highlighted {
			b.highlighted = false
			if i+1 == len(ui.menuButtons) {
				ui.menuButtons[0].highlighted = true
				return
			} else {
				ui.menuButtons[i+1].highlighted = true
				return
			}
		}
	}
}

func (ui *ui) getHighlightedButton() *menuButton {
	for _, b := range ui.menuButtons {
		if b.highlighted {
			return b
		}
	}
	return nil
}

func (ui *ui) doMenuAction() {
	button := ui.getHighlightedButton()

	switch button.name {
	case "Continue":
		ui.state = UIMain
	case "Quit":
		ui.inputChan <- &game.Input{Typ: game.CloseWindow, LevelChannel: ui.levelChan}
		ui.inputChan <- &game.Input{Typ: game.QuitGame}
	}
}

func (ui *ui) menuActions() {
	ui.displayMenu()
	for ui.state == UIMenu {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				ui.inputChan <- &game.Input{Typ: game.QuitGame}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					ui.inputChan <- &game.Input{Typ: game.CloseWindow, LevelChannel: ui.levelChan}
				}
			case *sdl.KeyboardEvent:
				if e.State != sdl.PRESSED {
					break
				}
				switch e.Keysym.Sym {
				case sdl.K_RETURN:
					ui.doMenuAction()
				case sdl.K_UP:
					ui.highlightPrevious()
					ui.displayMenu()
				case sdl.K_DOWN:
					ui.highlightNext()
					ui.displayMenu()
				case sdl.K_ESCAPE:
					return
				}
			}
		}
	}
}

func (ui *ui) displayMenu() {
	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("panel_beige.png"), &sdl.Rect{X: ui.invOffsetX, Y: ui.invOffsetY, W: ui.invWidth, H: ui.invHeight}); err != nil {
		panic(err)
	}

	buttonStandard := ui.getRectFromTextureName("buttonLong_brown.png")
	buttonHighlighted := ui.getRectFromTextureName("buttonLong_blue.png")

	for _, b := range ui.menuButtons {
		var button *sdl.Rect
		if b.highlighted {
			button = buttonHighlighted
		} else {
			button = buttonStandard
		}
		if err := ui.renderer.Copy(ui.uipack, button, b.buttonRect); err != nil {
			panic(err)
		}
		err := ui.renderer.Copy(b.buttonTexture, nil, b.buttonTextRect)
		if err != nil {
			panic(err)
		}
	}
	ui.renderer.Present()

}
