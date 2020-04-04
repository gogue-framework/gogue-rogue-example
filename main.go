package main

import (
	"github.com/gogue-framework/gogue"
	"runtime"
)

var (
	windowHeight int
	windowWidth int
)

func init() {
	runtime.LockOSThread()
	
	windowWidth = 50
	windowHeight = 25
	gogue.InitConsole(windowWidth, windowHeight, "Gogue Powered Roguelike", false)
	gogue.SetPrimaryFont(16, "press-start.ttf")
}

func main() {
	gogue.SetCompositionMode(0)

	text := "Welcome to Gogue!"
	gogue.PrintText((windowWidth / 2) -  (len(text) / 2), windowHeight / 2, 0, 0, text, "white", "", 0)
	
	gogue.Refresh()

	for {
		key := gogue.ReadInput(false)

		if key == gogue.KEY_CLOSE || key == gogue.KEY_ESCAPE {
			break
		}

		gogue.Refresh()
	}
	gogue.CloseConsole()
}
