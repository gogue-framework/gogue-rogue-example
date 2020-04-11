package main

import (
	"github.com/gogue-framework/gogue/ecs"
	"github.com/gogue-framework/gogue/screens"
	"github.com/gogue-framework/gogue/ui"
	"reflect"
	"runtime"
)

var (
	windowHeight int
	windowWidth  int

	// Global screen manager
	screenManager *screens.ScreenManager

	// Global ECS controller
	ecsController *ecs.Controller

	// Player entity
	player int
)

func init() {
	runtime.LockOSThread()

	windowWidth = 50
	windowHeight = 25
	ui.InitConsole(windowWidth, windowHeight, "Gogue Powered Roguelike", false)

	screenManager = screens.NewScreenManager()
	ecsController = ecs.NewController()

	registerComponents()

	playerPosition := PositionComponent{5, 5}
	playerAppearance := AppearanceComponent{
		Glyph:       ui.NewGlyph("@", "white", "white"),
		Layer:       0,
		Name:        "Player",
		Description: "The player character",
	}
	player = ecsController.CreateEntity([]ecs.Component{playerPosition, playerAppearance})
}

// registerScreens initializes and adds any game screens
func registerScreens() {
	// TitleScreen is the main screen for the game. It is shown when the game first launches
	titleScreen := TitleScreen{}

	// GameScreen is the main gameplay screen
	GameScreen := GameScreen{}

	_ = screenManager.AddScreen("title", &titleScreen)
	_ = screenManager.AddScreen("game", &GameScreen)
}

// registerComponent attaches, via a defined name, any component classes we want to use with the ECS controller
func registerComponents() {
	ecsController.MapComponentClass("position", &PositionComponent{})
	ecsController.MapComponentClass("appearance", &AppearanceComponent{})
}

func registerSystems() {
	render := SystemRender{ecsController: ecsController}
	ecsController.AddSystem(render, 1)
}

func main() {
	// Register all screens
	registerScreens()
	registerSystems()
	excludedSystems := []reflect.Type{}

	// Set the current screen to the title
	_ = screenManager.SetScreenByName("title")
	screenManager.CurrentScreen.Render()

	ui.Refresh()

	for {
		// First up, clear everything off the screen
		for i := 0; i <= 2; i++ {
			ui.ClearWindow(windowWidth, windowHeight, i)
		}

		screenManager.CurrentScreen.HandleInput()

		// Check if the current screen requires the ECS. If so, process all systems. If not, handle any input the screen
		// requires, and do not process and ECS systems
		if screenManager.CurrentScreen.UseEcs() {
			ecsController.Process(excludedSystems)
		}

		screenManager.CurrentScreen.Render()
		ui.Refresh()
	}

	ui.CloseConsole()
}
