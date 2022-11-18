package ui2d

import (
	"AirPygee/game"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) LoadPlayer() {
	ui.pCurrentFrame = 0
	ui.pFramesX = 3
	ui.pFramesY = 4

	image, err := img.Load("ui2d/assets/chara2.png")
	if err != nil {
		panic(err)
	}
	defer image.Free()
	image.W /= 4
	image.H /= 2
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
	if ui.pCurrentFrame >= ui.pFramesX {
		ui.pCurrentFrame = 0
	}
	switch input {
	case game.Up:
		ui.pFromY = 3 * ui.pHeightTex
	case game.Down:
		ui.pFromY = 0
	case game.Left:
		ui.pFromY = ui.pHeightTex
	case game.Right:
		ui.pFromY = 2 * ui.pHeightTex
	}
}

func (ui *ui) drawPlayer(level *game.Level) {
	p := level.Player
	ui.pFromX = ui.pCurrentFrame * ui.pWidthTex

	ui.pSrc = sdl.Rect{X: ui.pFromX, Y: ui.pFromY, W: ui.pWidthTex, H: ui.pHeightTex}
	ui.pDest = sdl.Rect{X: int32(p.X)*tileSize + ui.offsetX - (int32(float64(ui.pWidthTex)*1.25) - tileSize), Y: int32(p.Y)*tileSize + ui.offsetY - (int32(float64(ui.pHeightTex)*1.25) - tileSize), W: int32(float64(ui.pWidthTex) * 1.25), H: int32(float64(ui.pHeightTex) * 1.25)}

	err := ui.renderer.Copy(ui.pTexture, &ui.pSrc, &ui.pDest)
	game.CheckError(err)
}
