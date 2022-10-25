package game

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Player struct {
	Character
	Filename                                     string
	Width, Height                                int32
	FromX, FromY, FramesX, FramesY, CurrentFrame int32
	Texture                                      *sdl.Texture
	Src                                          sdl.Rect
	Dest                                         sdl.Rect
}

func (p *Player) Update(input InputType) {
	p.CurrentFrame++
	if p.CurrentFrame >= p.FramesY {
		p.CurrentFrame = 0
	}
	switch input {
	case Up:
		p.FromX = 2 * p.Width
	case Down:
		p.FromX = 0
	case Left:
		p.FromX = p.Width
	case Right:
		p.FromX = 3 * p.Width
	}

}

func NewPlayer(file string, framesX, framesY int32) *Player {
	player := &Player{Character: Character{
		Entity:       Entity{Name: "Wizard", Rune: '@'},
		Hitpoints:    20,
		MaxHitpoints: 20,
		Strength:     20,
		Speed:        1.0,
		ActionPoints: 0,
		SightRange:   10,
	}, Filename: file, FramesX: framesX, FramesY: framesY, CurrentFrame: 0}
	return player
}
