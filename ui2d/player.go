package ui2d

import (
	"AirPygee/game"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

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
