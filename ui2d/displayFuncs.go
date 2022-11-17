package ui2d

import (
	"AirPygee/game"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"strconv"
	"time"
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

	// Drawing Damage count
	tex = ui.stringToTexture("Damage "+strconv.Itoa(level.Player.MinDamage)+" - "+strconv.Itoa(level.Player.MaxDamage), sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ = tex.Query()

	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: int32(float64(statsPanel.W) * .10), Y: statsPanelOffsetY + int32(float64(statsPanel.H)*.25), W: w, H: h})
	if err != nil {
		panic(err)
	}

	// Drawing Defense count
	tex = ui.stringToTexture("Armor "+strconv.Itoa(level.Player.Armor), sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ = tex.Query()

	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: int32(float64(statsPanel.W) * .10), Y: statsPanelOffsetY + int32(float64(statsPanel.H)*.45), W: w, H: h})
	if err != nil {
		panic(err)
	}

	// Drawing Critical count
	tex = ui.stringToTexture("Critical "+fmt.Sprintf("%.2f %%", level.Player.Critical), sdl.Color{R: 139, G: 69, B: 19}, FontMedium)
	_, _, w, h, _ = tex.Query()

	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: int32(float64(statsPanel.W) * .10), Y: statsPanelOffsetY + int32(float64(statsPanel.H)*.65), W: w, H: h})
	if err != nil {
		panic(err)
	}

}

func (ui *ui) getColorFromHealth(health float64) (r, g, b, a uint8) {
	switch {
	case health <= 0.25:
		return 255, 0, 0, 0
	case health <= 0.50:
		return 255, 255, 0, 0
	case health <= 0.75:
		return 0, 255, 0, 0
	case health <= 1:
		return 0, 150, 0, 0
	}
	return 0, 0, 0, 0
}

// displayHUD draws general UI with remaining hit points and game instructions
func (ui *ui) displayHUD(level *game.Level) {
	//TODO - fix hardcoded value here
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
	game.CheckError(err)

	// Life gauge using red rect on black rect
	err = ui.renderer.FillRect(&sdl.Rect{X: int32(level.Player.Pos.X)*tileSize + ui.offsetX, Y: int32(level.Player.Pos.Y-1)*tileSize + ui.offsetY + 20, W: tileSize, H: 5})
	game.CheckError(err)

	var gauge float64
	gauge = float64(level.Player.Health) / float64(level.Player.MaxHealth)

	err = ui.renderer.SetDrawColor(ui.getColorFromHealth(gauge))
	game.CheckError(err)

	err = ui.renderer.FillRect(&sdl.Rect{X: int32(level.Player.Pos.X)*tileSize + ui.offsetX, Y: int32(level.Player.Pos.Y-1)*tileSize + ui.offsetY + 20, W: int32(float64(tileSize) * gauge), H: 5})
	game.CheckError(err)
	err = ui.renderer.SetDrawColor(0, 0, 0, 0)
	game.CheckError(err)
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

	err := ui.renderer.Copy(popup, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY, W: popupWidth, H: popupHeight})
	game.CheckError(err)

	// display item specific
	switch item.GetEntity().Type {
	case game.Potions:
		// display item Name
		tex := ui.stringToTexture(item.GetName(), color, FontMedium)
		_, _, w, h, _ := tex.Query()
		err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: mouseX - (popupWidth / 2) - (w / 2), Y: mouseY + int32(float64(popupHeight)*.05), W: w, H: h})
		game.CheckError(err)

		texPotion := ui.stringToTexture("Size: "+item.(game.ConsumableItem).GetSize(), color, FontSmall)
		_, _, w, h, _ = texPotion.Query()
		err = ui.renderer.Copy(texPotion, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY + int32(float64(popupHeight)*.65), W: w, H: h})
		game.CheckError(err)

	case game.Weapons, game.Armors:
		// display item Name
		tex := ui.stringToTexture(item.(game.EquipableItem).ToString(item.(game.EquipableItem).GetRarity())+" "+item.GetName(), color, FontMedium)
		_, _, w, h, _ := tex.Query()
		err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: mouseX - (popupWidth / 2) - (w / 2), Y: mouseY + int32(float64(popupHeight)*.05), W: w, H: h})
		game.CheckError(err)

		tex = ui.stringToTexture(fmt.Sprintf("Damage: %d - %d", item.(game.EquipableItem).GetStats().MinDamage, item.(game.EquipableItem).GetStats().MaxDamage), color, FontSmall)
		_, _, w, h, _ = tex.Query()
		err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY + int32(float64(popupHeight)*.45), W: w, H: h})
		game.CheckError(err)

		tex = ui.stringToTexture(fmt.Sprintf("Armor: %d", item.(game.EquipableItem).GetStats().Armor), color, FontSmall)
		_, _, w, h, _ = tex.Query()
		err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY + int32(float64(popupHeight)*.55), W: w, H: h})
		game.CheckError(err)

		tex = ui.stringToTexture(fmt.Sprintf("Crit Chance: %.2f %% ", item.(game.EquipableItem).GetStats().Critical), color, FontSmall)
		_, _, w, h, _ = tex.Query()
		err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY + int32(float64(popupHeight)*.65), W: w, H: h})
		game.CheckError(err)

		// display item rarity
		tex = ui.stringToTexture("Rarity: "+rarity, color, FontSmall)
		_, _, w, h, _ = tex.Query()
		err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY + int32(float64(popupHeight)*.75), W: w, H: h})
		game.CheckError(err)

	}

	// display item description
	tex := ui.stringToTexture("Description: "+item.GetDescription(), color, FontSmall)
	_, _, w, h, _ := tex.Query()
	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY + int32(float64(popupHeight)*.85), W: w, H: h})
	game.CheckError(err)
}

// displayMonsters displays monsters on map
func (ui *ui) displayMonsters(level *game.Level) {
	if err := ui.textureAtlas.SetColorMod(255, 255, 255); err != nil {
		panic(err)
	}
	for pos, monster := range level.Monsters {
		if level.Map[pos.Y][pos.X].Visible {
			err := ui.renderer.FillRect(&sdl.Rect{X: int32(level.Monsters[pos].X)*tileSize + ui.offsetX, Y: int32(level.Monsters[pos].Y-1)*tileSize + ui.offsetY + 20, W: tileSize, H: 5})
			game.CheckError(err)

			var gauge float64
			gauge = float64(level.Monsters[pos].Health) / float64(level.Monsters[pos].MaxHealth)

			err = ui.renderer.SetDrawColor(ui.getColorFromHealth(gauge))
			game.CheckError(err)

			// health bar
			err = ui.renderer.FillRect(&sdl.Rect{X: int32(level.Monsters[pos].X)*tileSize + ui.offsetX, Y: int32(level.Monsters[pos].Y-1)*tileSize + ui.offsetY + 20, W: int32(float64(tileSize) * gauge), H: 5})
			game.CheckError(err)

			ui.textureIndexMonsters.mu.RLock()
			monsterSrcRect := ui.textureIndexMonsters.rects[monster.Rune][0]
			ui.textureIndexMonsters.mu.RUnlock()

			err = ui.renderer.SetDrawColor(0, 0, 0, 0)
			game.CheckError(err)

			err = ui.renderer.Copy(ui.textureAtlas, &monsterSrcRect, &sdl.Rect{X: int32(pos.X)*tileSize + ui.offsetX, Y: int32(pos.Y)*tileSize + ui.offsetY, W: tileSize, H: tileSize})
			game.CheckError(err)
		}
	}
}

// displayItems displays items on Map
func (ui *ui) displayItems(level *game.Level) {
	for pos, items := range level.Items {
		if level.Map[pos.Y][pos.X].Visible {
			for _, item := range items {
				ui.textureIndexItems.mu.RLock()
				itemSrcRect := ui.textureIndexItems.rects[item.GetRune()][0]
				ui.textureIndexItems.mu.RUnlock()
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
	game.CheckError(err)

	i := level.EventPos
	count := 0

	for {
		event := level.Events[i]
		if event != "" {
			tex := ui.stringToTexture(event, sdl.Color{R: 100, G: 50}, FontSmall)
			_, _, w, h, _ := tex.Query()
			err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: textStartX, Y: int32(count*fontSizeY) + (int32(ui.winHeight) - (int32(ui.winHeight) - textStartY)), W: w, H: h})
			game.CheckError(err)

		}
		i = (i + 1) % (len(level.Events))
		count++
		if i == level.EventPos {
			break
		}
	}
}

func (ui *ui) displayDamages() {
	for _, damage := range ui.damagesToDisplay {
		_, _, w, h, _ := damage.tex.Query()
		err := ui.renderer.Copy(damage.tex, nil, &sdl.Rect{X: int32(damage.pos.X)*tileSize + ui.offsetX, Y: int32(damage.pos.Y)*tileSize + ui.offsetY, W: w, H: h})
		game.CheckError(err)
	}
}

func (ui *ui) buildAnimation(animation rune, texs ...*sdl.Rect) {
	for _, tex := range texs {
		ui.animations[animation] = append(ui.animations[animation], tex)
	}
}

func (ui *ui) displayTileAnimation(level *game.Level, duration time.Duration, p game.Pos, animation rune, textureIndex *TextureIndex) {
	tempTile := level.Map[p.Y][p.X].AnimRune
	level.Map[p.Y][p.X].AnimRune = animation
	numFrames := len(ui.animations[animation])

	for start := time.Now(); time.Since(start) < duration; {

		for i := 0; i < numFrames; i++ {
			if int(time.Since(start).Nanoseconds())%numFrames == i {
				textureIndex.mu.Lock()
				textureIndex.rects[animation] = nil
				textureIndex.rects[animation] = append(textureIndex.rects[animation], *ui.animations[animation][i])
				textureIndex.mu.Unlock()
			}
		}
	}
	level.Map[p.Y][p.X].AnimRune = tempTile
}

func (ui *ui) displayMovingAnimation(level *game.Level, duration time.Duration, animation rune, poss []game.Pos, textureIndex *TextureIndex) {
	for _, pos := range poss {
		ui.displayTileAnimation(level, duration, pos, animation, textureIndex)
	}
}

func (ui *ui) addAttackResult(damage int, duration time.Duration, isCritical bool, p game.Pos, who *game.Character) {
	now := time.Now().String()

	tex := ui.stringToTexture(strconv.Itoa(damage), sdl.Color{R: 255}, FontMedium)
	ui.damagesToDisplay[now] = &Damage{pos: p, tex: tex, isCritical: isCritical}
	for start := time.Now(); time.Since(start) < duration; {

	}
	delete(ui.damagesToDisplay, now)
}
