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
	var tex *sdl.Texture
	color := sdl.Color{R: 255, G: 215, B: 0, A: 255}
	statsColor := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	panelWidth := int32(float64(ui.winWidth) * 0.20)
	panelHeight := int32(float64(ui.winHeight) * 0.30)

	statsPanelOffsetY := int32(float64(ui.winHeight) * 0.10)

	err := ui.renderer.Copy(ui.metalPlate, nil, &sdl.Rect{X: 0, Y: statsPanelOffsetY, W: panelWidth, H: panelHeight})
	game.CheckError(err)

	// Drawing Health count
	tex = ui.stringToTexture("Life:", color, FontSmall)
	_, _, w, h, _ := tex.Query()

	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: int32(float64(panelWidth) * .15), Y: statsPanelOffsetY + int32(float64(panelHeight)*.15), W: w, H: h})
	game.CheckError(err)

	tex = ui.stringToTexture(fmt.Sprintf("%v", level.Player.Health), statsColor, FontSmall)
	_, _, w, h, _ = tex.Query()

	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: int32(float64(panelWidth) * .5), Y: statsPanelOffsetY + int32(float64(panelHeight)*.15), W: w, H: h})
	game.CheckError(err)

	// Drawing Damage count
	tex = ui.stringToTexture("Damage:", color, FontSmall)
	_, _, w, h, _ = tex.Query()

	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: int32(float64(panelWidth) * .15), Y: statsPanelOffsetY + int32(float64(panelHeight)*.25), W: w, H: h})
	game.CheckError(err)

	tex = ui.stringToTexture(fmt.Sprintf("%v - %v", level.Player.MinDamage, level.Player.MaxDamage), statsColor, FontSmall)
	_, _, w, h, _ = tex.Query()

	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: int32(float64(panelWidth) * .5), Y: statsPanelOffsetY + int32(float64(panelHeight)*.25), W: w, H: h})
	game.CheckError(err)

	// Drawing Defense count
	tex = ui.stringToTexture("Armor:", color, FontSmall)
	_, _, w, h, _ = tex.Query()

	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: int32(float64(panelWidth) * .15), Y: statsPanelOffsetY + int32(float64(panelHeight)*.35), W: w, H: h})
	game.CheckError(err)

	tex = ui.stringToTexture(fmt.Sprintf("%v", level.Player.Armor), statsColor, FontSmall)
	_, _, w, h, _ = tex.Query()

	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: int32(float64(panelWidth) * .5), Y: statsPanelOffsetY + int32(float64(panelHeight)*.35), W: w, H: h})
	game.CheckError(err)

	// Drawing Critical count
	tex = ui.stringToTexture("Critical:", color, FontSmall)
	_, _, w, h, _ = tex.Query()

	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: int32(float64(panelWidth) * .15), Y: statsPanelOffsetY + int32(float64(panelHeight)*.45), W: w, H: h})
	game.CheckError(err)

	tex = ui.stringToTexture(fmt.Sprintf("%.2f %%", level.Player.Critical), statsColor, FontSmall)
	_, _, w, h, _ = tex.Query()

	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: int32(float64(panelWidth) * .5), Y: statsPanelOffsetY + int32(float64(panelHeight)*.45), W: w, H: h})
	game.CheckError(err)
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
	//arrows keys
	// Move
	firstFrameX := 512
	for i := 0; i < 4; i++ {

		err := ui.renderer.Copy(ui.tileMap, &sdl.Rect{X: int32(firstFrameX + i*16), Y: 68, W: 16, H: 16}, &sdl.Rect{X: int32(i * 32), Y: 0, W: tileSize, H: tileSize})
		game.CheckError(err)

	}

	tex := ui.stringToTexture("Move", sdl.Color{R: 255}, FontSmall)
	_, _, w, h, _ := tex.Query()

	err := ui.renderer.Copy(tex, nil, &sdl.Rect{X: 144, Y: 8, W: w, H: h})
	game.CheckError(err)

	// Inventory

	err = ui.renderer.Copy(ui.tileMap, &sdl.Rect{X: 408, Y: 34, W: 16, H: 16}, &sdl.Rect{X: 98, Y: 32, W: tileSize, H: tileSize})
	game.CheckError(err)

	tex = ui.stringToTexture("Inventory", sdl.Color{R: 255}, FontSmall)
	_, _, w, h, _ = tex.Query()

	err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: 144, Y: 40, W: w, H: h})
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

	switch item.(type) {
	case game.EquipableItem:
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
	default:
	}

	err := ui.renderer.Copy(popup, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY, W: popupWidth, H: popupHeight})
	game.CheckError(err)

	// display item specific
	switch item.(type) {
	case game.ConsumableItem:
		// display item Name
		tex := ui.stringToTexture(item.GetName(), color, FontMedium)
		_, _, w, h, _ := tex.Query()
		err = ui.renderer.Copy(tex, nil, &sdl.Rect{X: mouseX - (popupWidth / 2) - (w / 2), Y: mouseY + int32(float64(popupHeight)*.05), W: w, H: h})
		game.CheckError(err)

		texPotion := ui.stringToTexture("Size: "+item.(game.ConsumableItem).GetSize(), color, FontSmall)
		_, _, w, h, _ = texPotion.Query()
		err = ui.renderer.Copy(texPotion, nil, &sdl.Rect{X: mouseX - popupWidth, Y: mouseY + int32(float64(popupHeight)*.65), W: w, H: h})
		game.CheckError(err)

	case game.EquipableItem:
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
	err := ui.textureAtlas.SetColorMod(255, 255, 255)
	game.CheckError(err)
	for pos, monster := range level.Monsters {
		if level.Map[pos.Y][pos.X].Visible {
			err = ui.renderer.FillRect(&sdl.Rect{X: int32(level.Monsters[pos].X)*tileSize + ui.offsetX, Y: int32(level.Monsters[pos].Y-1)*tileSize + ui.offsetY + 20, W: tileSize, H: 5})
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

			err = ui.renderer.Copy(ui.textureAtlas, monsterSrcRect, &sdl.Rect{X: int32(pos.X)*tileSize + ui.offsetX, Y: int32(pos.Y)*tileSize + ui.offsetY, W: tileSize, H: tileSize})
			game.CheckError(err)
		}
	}
}

// displayItems displays items on Map
func (ui *ui) displayItems(level *game.Level) {
	for pos, items := range level.Items {
		if level.Map[pos.Y][pos.X].Visible {
			for _, item := range items {
				var size int32
				var itemSrcRect *sdl.Rect
				var itemSrcTex *sdl.Texture
				size = tileSize
				switch item.(type) {
				case game.ConsumableItem:
					ui.textureIndexItems.mu.RLock()
					itemSrcRect = ui.textureIndexItems.rects[item.GetRune()][0]
					ui.textureIndexItems.mu.RUnlock()
					itemSrcTex = ui.textureAtlas
					switch item.(game.ConsumableItem).GetSize() {
					case "Small":
						size = int32(float64(size) * .50)
					case "Medium":
						size = int32(float64(size) * .75)
					case "Large":
						size = int32(float64(size) * .95)
					}
				case game.EquipableItem:
					ui.textureIndexItems.mu.RLock()
					itemSrcRect = ui.textureIndexItems.rects[item.GetRune()][0]
					ui.textureIndexItems.mu.RUnlock()
					itemSrcTex = ui.textureAtlas
				case game.OpenableItem:
					itemSrcRect = ui.textureIndexChests.rects[rune(item.(game.OpenableItem).GetSize())][0]
					itemSrcTex = ui.chestsTex
					size = tileSize
				}

				err := ui.renderer.Copy(itemSrcTex, itemSrcRect, &sdl.Rect{X: int32(pos.X)*tileSize + ui.offsetX, Y: int32(pos.Y)*tileSize + ui.offsetY, W: size, H: size})
				game.CheckError(err)
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

// displayTileAnimation if passed a zero duration, will only play all frames once
func (ui *ui) displayTileAnimation(level *game.Level, duration time.Duration, tick time.Duration, p game.Pos, animation rune, textureIndex *TextureIndex, tex *sdl.Texture) {
	tempTile := level.Map[p.Y][p.X].AnimRune
	numFrames := len(textureIndex.rects[animation])
	started := true
	currentFrame := 0
	timeStart := time.Now()

	for range time.Tick(tick) {
		textureIndex.mu.RLock()

		ui.animations[animation] = &Animation{
			rect: textureIndex.rects[animation][currentFrame],
			tex:  tex,
		}
		textureIndex.mu.RUnlock()

		if started {
			level.Map[p.Y][p.X].AnimRune = animation
			started = false
		}

		currentFrame++
		if currentFrame == numFrames {
			if duration == 0 {
				return
			} else {
				currentFrame = 0
			}
		}
		if time.Since(timeStart) >= duration && duration > 0 {
			break
		}
	}
	level.Map[p.Y][p.X].AnimRune = tempTile
}

func (ui *ui) displayPlayerAnimation(duration time.Duration, tick time.Duration, animation rune, textureIndex *TextureIndex, tex *sdl.Texture) {
	numFrames := len(textureIndex.rects[animation])
	started := true
	currentFrame := 0
	timeStart := time.Now()

	for range time.Tick(tick) {
		textureIndex.mu.RLock()
		ui.currentAnim = &Animation{
			rect: textureIndex.rects[animation][currentFrame],
			tex:  tex,
		}
		textureIndex.mu.RUnlock()

		if started {
			ui.pAnimated = true
			started = false
		}

		currentFrame++
		if currentFrame == numFrames {
			if duration > 0 {
				currentFrame = 0
			} else {
				return
			}
		}
		if time.Since(timeStart) >= duration && duration > 0 {
			break
		}
	}
	ui.pAnimated = false
}

func (ui *ui) displayMovingAnimation(level *game.Level, duration time.Duration, tick time.Duration, poss []game.Pos, animation rune, textureIndex *TextureIndex, tex *sdl.Texture) {
	for _, pos := range poss {
		ui.displayTileAnimation(level, duration, tick, pos, animation, textureIndex, tex)
	}
}

func (ui *ui) addAttackResult(damage int, duration time.Duration, isCritical bool, p game.Pos) {
	now := time.Now().String()

	tex := ui.stringToTexture(strconv.Itoa(damage), sdl.Color{R: 255}, FontMedium)
	ui.damagesToDisplay[now] = &Damage{pos: p, tex: tex, isCritical: isCritical}
	for start := time.Now(); time.Since(start) < duration; {

	}
	delete(ui.damagesToDisplay, now)
}
