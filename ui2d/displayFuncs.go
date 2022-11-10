package ui2d

import (
	"AirPygee/game"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"strconv"
)

// displayStats is displaying player stats
func (ui *ui) displayStats(level *game.Level) {

	statsPanel := ui.getRectFromTextureName("panel_brown.png")
	statsPanelOffsetY := int32(float64(ui.winHeight) * 0.10)
	if err := ui.renderer.Copy(ui.uipack, statsPanel, &sdl.Rect{X: 0, Y: statsPanelOffsetY, W: int32(float64(ui.winWidth) * 0.20), H: int32(float64(ui.winHeight) * 0.30)}); err != nil {
		panic(err)
	}

	// Drawing Health count
	tex := ui.stringToTexture("Life "+strconv.Itoa(level.Player.Health), sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ := tex.Query()
	err := ui.renderer.Copy(tex, nil, &sdl.Rect{X: int32(float64(statsPanel.W) * .10), Y: statsPanelOffsetY + int32(float64(statsPanel.H)*.05), W: w, H: h})
	if err != nil {
		panic(err)
	}

	// Drawing Strength count
	tex = ui.stringToTexture("Strength "+strconv.Itoa(level.Player.Strength), sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ = tex.Query()
	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: int32(float64(statsPanel.W) * .10), Y: statsPanelOffsetY + int32(float64(statsPanel.H)*.25), W: w, H: h})
	if err != nil {
		panic(err)
	}

}

// displayHUD draws general UI with remaining hit points and game instructions
func (ui *ui) displayHUD(level *game.Level) {
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

func (ui *ui) displayPopupItem(item game.Item, mouseX, mouseY int32) {
	popup := ui.getSinglePixel(sdl.Color{A: 128})
	popup.SetBlendMode(sdl.BLENDMODE_BLEND)

	popupWidth := int32(float64(ui.winWidth) * .25)
	popupHeight := int32(float64(ui.winHeight) * .25)
	color := sdl.Color{R: 225, G: 225, B: 225}
	var rarity string

	if item.GetEntity().Type != game.Potions {
		switch item.(game.EquipableItem).GetRarity() {
		case game.Common:
			color = sdl.Color{R: 225, G: 225, B: 225}
			rarity = "Common"
		case game.Uncommon:
			color = sdl.Color{R: 0, G: 225, B: 0}
			rarity = "Uncommon"
		case game.Rare:
			color = sdl.Color{R: 0, G: 0, B: 225}
			rarity = "Rare"
		case game.Epic:
			color = sdl.Color{R: 225, G: 0, B: 225}
			rarity = "Epic"
		case game.Legendary:
			color = sdl.Color{R: 225, G: 225, B: 0}
			rarity = "Legendary"
		}
	}

	ui.renderer.Copy(popup, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY, W: popupWidth, H: popupHeight})

	// display item Name
	tex := ui.stringToTexture(item.GetName(), color, FontMedium)
	_, _, w, h, _ := tex.Query()
	ui.renderer.Copy(tex, nil, &sdl.Rect{X: mouseX - (popupWidth / 2) - (w / 2), Y: mouseY + int32(float64(popupHeight)*.05), W: w, H: h})

	// display item specific
	switch item.GetEntity().Type {
	case game.Potions:
		texPotion := ui.stringToTexture("Size: "+item.(game.ConsumableItem).GetSize(), color, FontSmall)
		_, _, w, h, _ := texPotion.Query()
		ui.renderer.Copy(texPotion, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY + int32(float64(popupHeight)*.65), W: w, H: h})

	case game.Weapons, game.Armors:
		tex = ui.stringToTexture(fmt.Sprintf("Strength: %d - %d", item.(game.EquipableItem).GetStats().MinStrength, item.(game.EquipableItem).GetStats().MaxStrength), color, FontSmall)
		_, _, w, h, _ = tex.Query()
		ui.renderer.Copy(tex, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY + int32(float64(popupHeight)*.45), W: w, H: h})

		tex = ui.stringToTexture(fmt.Sprintf("Defense: %d - %d", item.(game.EquipableItem).GetStats().MinDefense, item.(game.EquipableItem).GetStats().MaxDefense), color, FontSmall)
		_, _, w, h, _ = tex.Query()
		ui.renderer.Copy(tex, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY + int32(float64(popupHeight)*.55), W: w, H: h})

		tex = ui.stringToTexture(fmt.Sprintf("Critical: %.2f %% - %.2f %%", item.(game.EquipableItem).GetStats().MinCritical, item.(game.EquipableItem).GetStats().MaxCritical), color, FontSmall)
		_, _, w, h, _ = tex.Query()
		ui.renderer.Copy(tex, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY + int32(float64(popupHeight)*.65), W: w, H: h})

		// display item rarity
		tex = ui.stringToTexture("Rarity: "+rarity, color, FontSmall)
		_, _, w, h, _ = tex.Query()
		ui.renderer.Copy(tex, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY + int32(float64(popupHeight)*.75), W: w, H: h})
	}

	// display item description
	tex = ui.stringToTexture("Description: "+item.GetDescription(), color, FontSmall)
	_, _, w, h, _ = tex.Query()
	ui.renderer.Copy(tex, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY + int32(float64(popupHeight)*.85), W: w, H: h})

}

// displayMonsters displays monsters on map
func (ui *ui) displayMonsters(level *game.Level) {
	if err := ui.textureAtlas.SetColorMod(255, 255, 255); err != nil {
		panic(err)
	}
	for pos, monster := range level.Monsters {
		if level.Map[pos.Y][pos.X].Visible {

			if err := ui.renderer.Copy(ui.textureAtlas, &sdl.Rect{X: 928, Y: 1600, W: tileSize, H: tileSize}, &sdl.Rect{X: int32(level.Monsters[pos].X)*tileSize + ui.offsetX, Y: int32(level.Monsters[pos].Y-1)*tileSize + ui.offsetY + 20, W: tileSize, H: 5}); err != nil {
				panic(err)
			}
			var gauge float64
			gauge = float64(level.Monsters[pos].Health) / float64(level.Monsters[pos].MaxHealth)

			if err := ui.renderer.Copy(ui.textureAtlas, &sdl.Rect{X: 1024, Y: 1600, W: tileSize, H: tileSize}, &sdl.Rect{X: int32(level.Monsters[pos].X)*tileSize + ui.offsetX, Y: int32(level.Monsters[pos].Y-1)*tileSize + ui.offsetY + 20, W: tileSize * int32(gauge), H: 5}); err != nil {
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
				itemSrcRect := ui.textureIndex[item.GetRune()][0]
				var size int32
				size = tileSize
				if item.GetName() == "Potion" {
					switch item.(game.ConsumableItem).GetSize() {
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

			err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: textStartX, Y: int32(count*fontSizeY) + (int32(ui.winHeight) - (int32(ui.winHeight) - textStartY)), W: w, H: h})
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
