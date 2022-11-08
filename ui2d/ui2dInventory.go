package ui2d

import (
	"AirPygee/game"
	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) drawInventory(level *game.Level) {
	var locationX, locationY int32

	playerSrcRect := sdl.Rect{X: 0, Y: 0, W: 26, H: 36}
	playerX := ((ui.invWidth - (ui.invWidth / 3)) / 2) + ui.invOffsetX
	playerY := ((ui.invHeight - (ui.invHeight / 3)) / 2) + ui.invOffsetY

	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("panel_beige.png"), &sdl.Rect{X: ui.invOffsetX, Y: ui.invOffsetY, W: ui.invWidth, H: ui.invHeight}); err != nil {
		panic(err)
	}

	if err := ui.renderer.Copy(ui.pTexture, &playerSrcRect, &sdl.Rect{X: playerX, Y: playerY, W: ui.invWidth / 3, H: ui.invHeight / 3}); err != nil {
		panic(err)
	}

	ui.drawEmptyInventory(level)

	// draw panel items
	//Head
	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("buttonSquare_blue_pressed.png"), &sdl.Rect{X: ui.invHeadX, Y: ui.invHeadY, W: ui.itemW, H: ui.itemH}); err != nil {
		panic(err)
	}
	//RightHand
	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("buttonSquare_blue_pressed.png"), &sdl.Rect{X: ui.invRHandX, Y: ui.invRHandY, W: ui.itemW, H: ui.itemH}); err != nil {
		panic(err)
	}
	//LeftHand
	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("buttonSquare_blue_pressed.png"), &sdl.Rect{X: ui.invLHandX, Y: ui.invLHandY, W: ui.itemW, H: ui.itemH}); err != nil {
		panic(err)
	}
	//Foots
	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("buttonSquare_blue_pressed.png"), &sdl.Rect{X: ui.invFootsX, Y: ui.invFootsY, W: ui.itemW, H: ui.itemH}); err != nil {
		panic(err)
	}
	//Chest
	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("buttonSquare_blue_pressed.png"), &sdl.Rect{X: ui.invChestX, Y: ui.invChestY, W: ui.itemW, H: ui.itemH}); err != nil {
		panic(err)
	}
	//Legs
	if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("buttonSquare_blue_pressed.png"), &sdl.Rect{X: ui.invLegsX, Y: ui.invLegsY, W: ui.itemW, H: ui.itemH}); err != nil {
		panic(err)
	}

	for i, item := range level.Player.EquippedItems {
		itemSrcRect := ui.textureIndex[item.Rune][0]

		switch item.Location {
		case game.Head:
			locationX = ui.invHeadX
			locationY = ui.invHeadY
		case game.RightHand:
			locationX = ui.invRHandX
			locationY = ui.invRHandY
		case game.LeftHand:
			locationX = ui.invLHandX
			locationY = ui.invLHandY
		case game.Foots:
			locationX = ui.invFootsX
			locationY = ui.invFootsY
		case game.Chest:
			locationX = ui.invChestX
			locationY = ui.invChestY
		case game.Legs:
			locationX = ui.invLegsX
			locationY = ui.invLegsY
		default:
			locationX = ui.invOffsetX + int32(i)*ui.itemW
			locationY = ui.invOffsetY + ui.invHeight - ui.itemH
		}

		if item == ui.draggedItem {
			if err := ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{X: int32(ui.currentMouseState.pos.X) - ui.itemW/2, Y: int32(ui.currentMouseState.pos.Y) - ui.itemH/2, W: ui.itemW, H: ui.itemH}); err != nil {
				panic(err)
			}
		} else {
			if err := ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{X: locationX, Y: locationY, W: ui.itemW, H: ui.itemH}); err != nil {
				panic(err)
			}
		}
	}

	var countX int32 = 0
	var countY int32 = 0
	for _, item := range level.Player.Items {
		itemSrcRect := ui.textureIndex[item.Rune][0]
		if countX%5 == 0 {
			countX = 0
			countY++
		}
		locationX = ui.invOffsetX + ui.invWidth + ui.itemW*countX
		locationY = ui.invOffsetY + ui.itemH*countY

		if item == ui.draggedItem {
			if err := ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{X: int32(ui.currentMouseState.pos.X) - ui.itemW/2, Y: int32(ui.currentMouseState.pos.Y) - ui.itemH/2, W: ui.itemW, H: ui.itemH}); err != nil {
				panic(err)
			}
		} else {
			var size int32
			size = ui.itemW
			if item.Name == "Potion" {
				switch item.Size {
				case "Small":
					size = int32(float64(size) * .50)
					locationX += ui.itemW/2 - size/2
					locationY += ui.itemH/2 - size/2
				case "Medium":
					size = int32(float64(size) * .75)
					locationX += ui.itemW/2 - size/2
					locationY += ui.itemH/2 - size/2
				case "Large":
					size = int32(float64(size) * .95)
					locationX += ui.itemW/2 - size/2
					locationY += ui.itemH/2 - size/2
				}
			}
			if err := ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{X: locationX, Y: locationY, W: size, H: size}); err != nil {
				panic(err)
			}
		}
		countX++
	}
}

func (ui *ui) drawEmptyInventory(level *game.Level) {
	var countX, countY, locationX, locationY int32
	for i := 0; i < level.Player.InventorySize; i++ {
		if i%5 == 0 {
			countX = 0
			countY++
		}
		locationX = ui.invOffsetX + ui.invWidth + ui.itemW*countX
		locationY = ui.invOffsetY + ui.itemH*countY
		if err := ui.renderer.Copy(ui.uipack, ui.getRectFromTextureName("buttonSquare_brown_pressed.png"), &sdl.Rect{X: locationX, Y: locationY, W: ui.itemW, H: ui.itemH}); err != nil {
			panic(err)
		}
		countX++
	}
}

func (ui *ui) getInventoryItemRect(id int, level *game.Level) *sdl.Rect {
	var locationX, locationY, countX, countY int32

	for i := 0; i < level.Player.InventorySize; i++ {
		if i%5 == 0 {
			countX = 0
			countY++
		}
		locationX = ui.invOffsetX + ui.invWidth + ui.itemW*countX
		locationY = ui.invOffsetY + ui.itemH*countY

		countX++
		if i == id {
			return &sdl.Rect{X: locationX, Y: locationY, W: ui.itemW, H: ui.itemH}
		}
	}
	return nil
}

// getEquippedItemRect based on arbitraries items positions, will return the corresponding
// rectangle in order to compare with click position and then unequip
func (ui *ui) getEquippedItemRect(item *game.Item) *sdl.Rect {
	var locationX, locationY int32

	switch item.Location {
	case game.Head:
		locationX = ui.invHeadX
		locationY = ui.invHeadY
	case game.Foots:
		locationX = ui.invFootsX
		locationY = ui.invFootsY
	case game.LeftHand:
		locationX = ui.invLHandX
		locationY = ui.invLHandY
	case game.RightHand:
		locationX = ui.invRHandX
		locationY = ui.invRHandY
	case game.Chest:
		locationX = ui.invChestX
		locationY = ui.invChestY
	case game.Legs:
		locationX = ui.invLegsX
		locationY = ui.invLegsY
	}

	return &sdl.Rect{X: locationX, Y: locationY, W: ui.itemW, H: ui.itemH}
}

func (ui *ui) clickValidItem(level *game.Level, mouseX, mouseY int32) *game.Item {
	for i, item := range level.Player.Items {
		itemRect := ui.getInventoryItemRect(i, level)
		if itemRect.HasIntersection(&sdl.Rect{X: mouseX, Y: mouseY, W: 1, H: 1}) {
			ui.dragMode = fromInventory
			return item
		}
	}

	for _, item := range level.Player.EquippedItems {
		itemRect := ui.getEquippedItemRect(item)
		if itemRect.HasIntersection(&sdl.Rect{X: mouseX, Y: mouseY, W: 1, H: 1}) {
			ui.dragMode = fromEquippedItems
			return item
		}
	}

	return nil
}

func (ui *ui) isSlotFree(level *game.Level, itemToEquip *game.Item) bool {
	for _, item := range level.Player.EquippedItems {
		if itemToEquip.Location == item.Location {
			return false
		}
	}
	return true
}

func (ui *ui) hasClickedOnValidEquipSlot(mouseX, mouseY int32, item *game.Item) bool {
	itemRect := ui.getEquippedItemRect(item)
	if itemRect.HasIntersection(&sdl.Rect{X: mouseX, Y: mouseY, W: 1, H: 1}) {
		return true
	} else {
		return false
	}
}

func (ui *ui) hasClickedInBackpackZone(mouseX, mouseY int32) bool {
	itemRect := &sdl.Rect{X: ui.invOffsetX + ui.invWidth, Y: ui.invOffsetY + ui.itemH, W: ui.itemW * 5, H: ui.itemH * 4}
	if itemRect.HasIntersection(&sdl.Rect{X: mouseX, Y: mouseY, W: 1, H: 1}) {
		return true
	} else {
		return false
	}
}

func (ui *ui) hasClickedInEquipZone(mouseX, mouseY int32) bool {
	itemRect := &sdl.Rect{X: ui.invOffsetX, Y: ui.invOffsetY, W: ui.invWidth, H: ui.invHeight}
	if itemRect.HasIntersection(&sdl.Rect{X: mouseX, Y: mouseY, W: 1, H: 1}) {
		return true
	} else {
		return false
	}
}

func (ui *ui) hasClickedOutsideInventoryZone(mouseX, mouseY int32) bool {
	return !ui.hasClickedInBackpackZone(mouseX, mouseY) && !ui.hasClickedInEquipZone(mouseX, mouseY)
}
