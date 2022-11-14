package main

import (
	"AirPygee/game"
	"AirPygee/ui2d"
)

//func init() {
//	runtime.LockOSThread()
//}

func main() {

	// For multiple UI but doesn't work on MAC because of sdl.PollEvents
	//game := game.NewGame(1, "game/maps/level1.map")
	//
	//for i := 0; i < 1; i++ {
	//	go func(i int) {
	//		runtime.LockOSThread()
	//		ui := ui2d.NewUI(game.InputChan, game.LevelChans[i])
	//		ui.Run()
	//	}(i)
	//}
	//game.Run()

	game := game.NewGame(1)

	go game.Run()

	ui := ui2d.NewUI(game.InputChan, game.LevelChans[0])
	ui.Run()

}
