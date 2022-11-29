package ui2d

// TODO - implements adjustable volume

import (
	"AirPygee/game"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) buildDifficultyButtons() {
	button := ui.getRectFromTextureName("buttonLong_brown.png")
	button.W /= 2
	button.H /= 2

	//Easy
	tex := ui.stringToTexture("Easy", sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ := tex.Query()

	ui.difficultyButtons = append(ui.difficultyButtons, &menuButton{
		name:           "Easy",
		buttonRect:     &sdl.Rect{X: ui.invOffsetX + button.W/2, Y: ui.invOffsetY + button.H*8, W: button.W, H: button.H},
		buttonTexture:  tex,
		buttonTextRect: &sdl.Rect{X: ui.invOffsetX + button.W/2 + (button.W/2 - w/2), Y: ui.invOffsetY + button.H*8 + (button.H / 2) - (h / 2), W: w, H: h},
		highlighted:    true,
	})

	//Medium
	tex = ui.stringToTexture("Medium", sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ = tex.Query()

	ui.difficultyButtons = append(ui.difficultyButtons, &menuButton{
		name:           "Medium",
		buttonRect:     &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - button.W/2, Y: ui.invOffsetY + button.H*8, W: button.W, H: button.H},
		buttonTexture:  tex,
		buttonTextRect: &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - w/2, Y: ui.invOffsetY + button.H*8 + (button.H / 2) - (h / 2), W: w, H: h},
		highlighted:    false,
	})

	//Hard
	tex = ui.stringToTexture("Hard", sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ = tex.Query()

	ui.difficultyButtons = append(ui.difficultyButtons, &menuButton{
		name:           "Hard",
		buttonRect:     &sdl.Rect{X: ui.invOffsetX + ui.invWidth - button.W - button.W/2, Y: ui.invOffsetY + button.H*8, W: button.W, H: button.H},
		buttonTexture:  tex,
		buttonTextRect: &sdl.Rect{X: ui.invOffsetX + ui.invWidth - button.W - button.W/2 + (button.W/2 - w/2), Y: ui.invOffsetY + button.H*8 + (button.H / 2) - (h / 2), W: w, H: h},
		highlighted:    false,
	})
}

func (ui *ui) buildStartMenuButtons() {
	button := ui.getRectFromTextureName("buttonLong_brown.png")

	// Start button
	tex := ui.stringToTexture("Start", sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ := tex.Query()

	ui.startMenuButtons = append(ui.startMenuButtons, &menuButton{
		name:           "Start",
		buttonRect:     &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - button.W/2, Y: ui.invOffsetY + button.H, W: button.W, H: button.H},
		buttonTexture:  tex,
		buttonTextRect: &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - w/2, Y: ui.invOffsetY + button.H + (button.H / 2) - (h / 2), W: w, H: h},
		highlighted:    true,
	})

	// Difficulty button
	tex = ui.stringToTexture("Difficulty", sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ = tex.Query()

	ui.startMenuButtons = append(ui.startMenuButtons, &menuButton{
		name:           "Difficulty",
		buttonRect:     &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - button.W/2, Y: ui.invOffsetY + button.H*3, W: button.W, H: button.H},
		buttonTexture:  tex,
		buttonTextRect: &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - w/2, Y: ui.invOffsetY + button.H*3 + (button.H / 2) - (h / 2), W: w, H: h},
	})

	// Quit button

	tex = ui.stringToTexture("Quit", sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ = tex.Query()

	ui.startMenuButtons = append(ui.startMenuButtons, &menuButton{
		name:           "Quit",
		buttonRect:     &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - button.W/2, Y: ui.invOffsetY + button.H*5, W: button.W, H: button.H},
		buttonTexture:  tex,
		buttonTextRect: &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - w/2, Y: ui.invOffsetY + button.H*5 + (button.H / 2) - (h / 2), W: w, H: h},
	})

}

func (ui *ui) startMenuActions(level *game.Level) {
	ui.displayStartMenu()
	for ui.state == UIStartMenu || ui.state == UIStartMenuDifficulty {
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
				if ui.state == UIStartMenuDifficulty {
					switch e.Keysym.Sym {
					case sdl.K_RETURN:
						ui.doDifficultyMenuAction()
						ui.state = UIStartMenu
						ui.displayStartMenu()
					case sdl.K_LEFT:
						ui.highlightLeftDifficulty()
						ui.displayDifficulty()
					case sdl.K_RIGHT:
						ui.highlightRightDifficulty()
						ui.displayDifficulty()
					}
				} else if ui.state == UIStartMenu {
					switch e.Keysym.Sym {
					case sdl.K_RETURN:
						ui.doStartMenuAction(level)
					case sdl.K_UP:
						ui.highlightPreviousStartMenu()
						ui.displayStartMenu()
					case sdl.K_DOWN:
						ui.highlightNextStartMenu()
						ui.displayStartMenu()
					}
				}
			}
		}
	}
}

func (ui *ui) highlightLeftDifficulty() {
	for i, b := range ui.difficultyButtons {
		if b.highlighted {
			b.highlighted = false
			if i-1 < 0 {
				ui.difficultyButtons[len(ui.difficultyButtons)-1].highlighted = true
				return
			} else {
				ui.difficultyButtons[i-1].highlighted = true
				return
			}
		}
	}
}

func (ui *ui) highlightRightDifficulty() {
	for i, b := range ui.difficultyButtons {
		if b.highlighted {
			b.highlighted = false
			if i+1 == len(ui.difficultyButtons) {
				ui.difficultyButtons[0].highlighted = true
				return
			} else {
				ui.difficultyButtons[i+1].highlighted = true
				return
			}
		}
	}
}

func (ui *ui) displayStartMenu() {
	image, err := img.Load("ui2d/assets/startMenuBackground.jpeg")
	game.CheckError(err)

	defer image.Free()

	menuTex, err := ui.renderer.CreateTextureFromSurface(image)
	game.CheckError(err)

	err = ui.renderer.Copy(menuTex, nil, nil)
	game.CheckError(err)

	err = ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("panel_beige.png"), &sdl.Rect{X: ui.invOffsetX, Y: ui.invOffsetY, W: ui.invWidth, H: ui.invHeight})
	game.CheckError(err)

	buttonStandard := ui.getRectFromTextureName("buttonLong_brown.png")
	buttonHighlighted := ui.getRectFromTextureName("buttonLong_grey.png")

	for _, b := range ui.startMenuButtons {
		var button *sdl.Rect
		if b.highlighted {
			button = buttonHighlighted
		} else {
			button = buttonStandard
		}
		err = ui.renderer.Copy(ui.uipack, button, b.buttonRect)
		game.CheckError(err)
		err = ui.renderer.Copy(b.buttonTexture, nil, b.buttonTextRect)
		game.CheckError(err)
	}
	ui.renderer.Present()

}

func (ui *ui) highlightPreviousStartMenu() {
	for i, b := range ui.startMenuButtons {
		if b.highlighted {
			b.highlighted = false
			if i-1 < 0 {
				ui.startMenuButtons[len(ui.startMenuButtons)-1].highlighted = true
				return
			} else {
				ui.startMenuButtons[i-1].highlighted = true
				return
			}
		}
	}
}

func (ui *ui) highlightNextStartMenu() {
	for i, b := range ui.startMenuButtons {
		if b.highlighted {
			b.highlighted = false
			if i+1 == len(ui.startMenuButtons) {
				ui.startMenuButtons[0].highlighted = true
				return
			} else {
				ui.startMenuButtons[i+1].highlighted = true
				return
			}
		}
	}
}

func (ui *ui) getDifficultyHighlightedButton() *menuButton {
	for _, b := range ui.difficultyButtons {
		if b.highlighted {
			return b
		}
	}
	return nil
}

func (ui *ui) getStartMenuHighlightedButton() *menuButton {
	for _, b := range ui.startMenuButtons {
		if b.highlighted {
			return b
		}
	}
	return nil
}

func (ui *ui) displayDifficulty() {
	buttonStandard := ui.getRectFromTextureName("buttonLong_brown.png")
	buttonHighlighted := ui.getRectFromTextureName("buttonLong_grey.png")

	for _, b := range ui.difficultyButtons {
		var button *sdl.Rect
		if b.highlighted {
			button = buttonHighlighted
		} else {
			button = buttonStandard
		}
		err := ui.renderer.Copy(ui.uipack, button, b.buttonRect)
		game.CheckError(err)
		err = ui.renderer.Copy(b.buttonTexture, nil, b.buttonTextRect)
		game.CheckError(err)
	}
	ui.renderer.Present()

}

func (ui *ui) doDifficultyMenuAction() {
	button := ui.getDifficultyHighlightedButton()

	switch button.name {
	case "Easy":
		ui.inputChan <- &game.Input{Typ: game.SetDifficulty, Difficulty: 1}
	case "Medium":
		ui.inputChan <- &game.Input{Typ: game.SetDifficulty, Difficulty: 2}
	case "Hard":
		ui.inputChan <- &game.Input{Typ: game.SetDifficulty, Difficulty: 3}
	}
}

func (ui *ui) doStartMenuAction(level *game.Level) {
	button := ui.getStartMenuHighlightedButton()

	switch button.name {
	case "Start":
		ui.state = UIMain
		ui.inputChan <- &game.Input{Typ: game.Restart, LevelChannel: ui.levelChan}
	case "Difficulty":
		ui.state = UIStartMenuDifficulty
		ui.displayDifficulty()
	//TODO
	case "Quit":
		ui.inputChan <- &game.Input{Typ: game.CloseWindow, LevelChannel: ui.levelChan}
	}
}
