package ui2d

// TODO - implements adjustable volume

import (
	"AirPygee/game"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) buildMenuButtons() {
	button := ui.getRectFromTextureName("buttonLong_brown.png")

	// Quit button
	tex := ui.stringToTexture("Quit", sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ := tex.Query()

	ui.menuButtons = append(ui.menuButtons, &menuButton{
		name:           "Quit",
		buttonRect:     &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - button.W/2, Y: ui.invOffsetY + button.H, W: button.W, H: button.H},
		buttonTexture:  tex,
		buttonTextRect: &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - w/2, Y: ui.invOffsetY + button.H + (button.H / 2) - (h / 2), W: w, H: h},
		highlighted:    true,
	})

	// Continue button
	tex = ui.stringToTexture("Continue", sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ = tex.Query()

	ui.menuButtons = append(ui.menuButtons, &menuButton{
		name:           "Continue",
		buttonRect:     &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - button.W/2, Y: ui.invOffsetY + button.H*3, W: button.W, H: button.H},
		buttonTexture:  tex,
		buttonTextRect: &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - w/2, Y: ui.invOffsetY + button.H*3 + (button.H / 2) - (h / 2), W: w, H: h},
	})

	// music button

	tex = ui.stringToTexture("Music", sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ = tex.Query()

	ui.menuButtons = append(ui.menuButtons, &menuButton{
		name:           "Music",
		buttonRect:     &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - button.W/2, Y: ui.invOffsetY + button.H*5, W: button.W, H: button.H},
		buttonTexture:  tex,
		buttonTextRect: &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - w/2, Y: ui.invOffsetY + button.H*5 + (button.H / 2) - (h / 2), W: w, H: h},
	})

	// Sound button

	tex = ui.stringToTexture("Sound effects", sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ = tex.Query()

	ui.menuButtons = append(ui.menuButtons, &menuButton{
		name:           "Sound effects",
		buttonRect:     &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - button.W/2, Y: ui.invOffsetY + button.H*7, W: button.W, H: button.H},
		buttonTexture:  tex,
		buttonTextRect: &sdl.Rect{X: ui.invOffsetX + ui.invWidth/2 - w/2, Y: ui.invOffsetY + button.H*7 + (button.H / 2) - (h / 2), W: w, H: h},
	})

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
					ui.state = UIMain
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
	buttonHighlighted := ui.getRectFromTextureName("buttonLong_grey.png")

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
	case "Music":
		if mix.VolumeMusic(0) == 0 {
			mix.VolumeMusic(ui.musicVolume)
		} else {
			mix.VolumeMusic(0)
		}
	case "Sound effects":
		if ui.soundsVolume > 0 {
			ui.soundsVolume = 0
		} else {
			ui.soundsVolume = 10
		}
	case "Continue":
		ui.state = UIMain
	case "Quit":
		ui.inputChan <- &game.Input{Typ: game.CloseWindow, LevelChannel: ui.levelChan}
	}
}
