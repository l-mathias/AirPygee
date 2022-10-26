package ui2d

import (
	"AirPygee/game"
	"github.com/veandco/go-sdl2/sdl"
	"strconv"
)

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

func (ui *ui) displayEvents(level *game.Level) {
	textStartX := int32(float64(ui.winWidth) * .015)
	textStartY := int32(float64(ui.winHeight) * .68)
	textWidth := int32(float64(ui.winWidth) * .25)
	_, fontSizeY, _ := ui.fontSmall.SizeUTF8("A")

	err := ui.renderer.Copy(ui.borders, nil, &sdl.Rect{X: 0, Y: int32(ui.winHeight) - (int32(ui.winHeight) - textStartY + int32(fontSizeY)), W: textWidth, H: int32(ui.winHeight) - textStartY + int32(fontSizeY)})
	if err != nil {
		panic(err)
	}

	i := level.EventPos
	count := 0

	for {
		event := level.Events[i]
		if event != "" {
			tex := ui.stringToTexture(event, sdl.Color{R: 255}, FontSmall)
			_, _, w, h, _ := tex.Query()
			//err := ui.renderer.Copy(tex, nil, &sdl.Rect{X: textStartX, Y: int32(count*fontSizeY) + textStartY, W: w, H: h})

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

	// draw inventory
	inventoryStart := int32(float64(ui.winWidth) * 0.7)
	inventoryWidth := int32(ui.winWidth) - inventoryStart

	err = ui.renderer.Copy(ui.borders, nil, &sdl.Rect{X: inventoryStart, Y: int32(ui.winHeight - 32), W: inventoryWidth, H: 32})
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
