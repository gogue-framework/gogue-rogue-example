package main

import (
	"github.com/gogue-framework/gogue/screens"
	"github.com/gogue-framework/gogue/ui"
	"runtime"
)

var (
	windowHeight int
	windowWidth  int

	// Global screen manager
	screenManager *screens.ScreenManager
)

func init() {
	runtime.LockOSThread()

	windowWidth = 50
	windowHeight = 25
	ui.InitConsole(windowWidth, windowHeight, "Gogue Powered Roguelike", false)

	screenManager = screens.NewScreenManager()
}

// registerScreens initializes and adds any game screens
func registerScreens() {
	// TitleScreen is the main screen for the game. It is shown when the game first launches
	titleScreen := TitleScreen{}

	// GameScreen is the main gameplay screen
	GameScreen := GameScreen{}

	screenManager.AddScreen("title", &titleScreen)
	screenManager.AddScreen("game", &GameScreen)
}

func main() {
	// Register all screens
	registerScreens()

	// Set the current screen to the title
	screenManager.SetScreenByName("title")
	screenManager.CurrentScreen.Render()

	ui.Refresh()

	for {
		screenManager.CurrentScreen.HandleInput()
		screenManager.CurrentScreen.Render()
		ui.Refresh()
	}

	ui.CloseConsole()
}
