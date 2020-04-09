package main

import (
	"github.com/gogue-framework/gogue/ui"
	"runtime"
)

var (
	windowHeight int
	windowWidth  int
)

func init() {
	runtime.LockOSThread()

	windowWidth = 50
	windowHeight = 25
	ui.InitConsole(windowWidth, windowHeight, "Gogue Powered Roguelike", false)
	ui.SetPrimaryFont(16, "press-start.ttf")
}

func main() {
	ui.SetCompositionMode(0)

	text := "Welcome to Gogue!"
	ui.PrintText((windowWidth/2)-(len(text)/2), windowHeight/2, 0, 0, text, "white", "", 0)

	ui.Refresh()

	for {
		key := ui.ReadInput(false)

		if key == ui.KeyClose || key == ui.KeyEscape {
			break
		}

		ui.Refresh()
	}
	ui.CloseConsole()
}
