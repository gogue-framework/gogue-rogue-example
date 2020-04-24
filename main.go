package main

import (
	"github.com/gogue-framework/gogue/camera"
	"github.com/gogue-framework/gogue/ecs"
	"github.com/gogue-framework/gogue/fov"
	"github.com/gogue-framework/gogue/gamemap"
	"github.com/gogue-framework/gogue/gamemap/maptypes"
	"github.com/gogue-framework/gogue/randomnumbergenerator"
	"github.com/gogue-framework/gogue/screens"
	"github.com/gogue-framework/gogue/ui"
	"reflect"
	"runtime"
)

var (
	windowHeight int
	windowWidth  int

	mapWidth int
	mapHeight int

	rng *randomnumbergenerator.RNG

	// Global screen manager
	screenManager *screens.ScreenManager

	// Global ECS controller
	ecsController *ecs.Controller

	// GameMap
	gameMap *gamemap.GameMap
	wallGlyph ui.Glyph
	floorGlyph ui.Glyph

	// Game Camera
	gameCamera *camera.GameCamera

	// Player entity
	player int
	playerFOV fov.FieldOfVision
	torchRadius int
)

func init() {
	runtime.LockOSThread()

	windowWidth = 50
	windowHeight = 25
	ui.InitConsole(windowWidth, windowHeight, "Gogue Powered Roguelike", false)

	rng = randomnumbergenerator.NewRNG()

	screenManager = screens.NewScreenManager()
	ecsController = ecs.NewController()

	registerComponents()

	playerPosition := PositionComponent{5, 5}
	playerAppearance := AppearanceComponent{
		Glyph:       ui.NewGlyph("@", "white", "white"),
		Layer:       1,
		Name:        "Player",
		Description: "The player character",
	}
	playerMovement := MovementComponent{}
	player = ecsController.CreateEntity([]ecs.Component{playerPosition, playerAppearance, playerMovement})

	wallGlyph = ui.NewGlyph("#", "white", "gray")
	floorGlyph = ui.NewGlyph(".", "white", "gray")

	// Set the map width and height
	mapWidth = 100
	mapHeight = 100

	// Initialize the game camera
	gameCamera, _ = camera.NewGameCamera(1, 1, windowWidth, windowHeight)

	// Initialize the players FOV
	playerFOV.InitializeFOV()
	torchRadius = 5
	playerFOV.SetTorchRadius(torchRadius)
}

// SetupGameMap initializes a new GameMap, with a fixed width and height (this can be larger than the game window)
// Initially, all map tiles are set to floor
func SetupGameMap() *gamemap.GameMap {
	gameMap = &gamemap.GameMap{Width:mapWidth, Height:mapHeight}
	gameMap.InitializeMap()
	maptypes.GenerateCavern(gameMap, wallGlyph, floorGlyph, 50)

	// Set a random starting position for the player
	randomTile := gameMap.FloorTiles[rng.Range(0, len(gameMap.FloorTiles))]
	newPos := PositionComponent{X: randomTile.X, Y: randomTile.Y}
	ecsController.UpdateComponent(player, PositionComponent{}.TypeOf(), newPos)

	return gameMap
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
	ecsController.MapComponentClass("movement", &MovementComponent{})
}

func registerSystems() {
	render := SystemRender{ecsController: ecsController}
	input := SystemInput{ecsController: ecsController}

	ecsController.AddSystem(input, 0)
	ecsController.AddSystem(render, 1)
}

func main() {
	// Register all screens
	registerScreens()
	registerSystems()
	excludedSystems := []reflect.Type{reflect.TypeOf(SystemRender{})}

	// Initialize the game map
	gameMap = SetupGameMap()

	// Set the current screen to the title
	_ = screenManager.SetScreenByName("title")
	screenManager.CurrentScreen.Render()

	ui.Refresh()

	for {
		// First up, clear everything off the screen
		for i := 0; i <= 2; i++ {
			ui.ClearWindow(windowWidth, windowHeight, i)
		}

		// Process the players Field of Vision. This will dictate what they can and cant see each turn.
		playerPosition := ecsController.GetComponent(player, PositionComponent{}.TypeOf()).(PositionComponent)
		playerFOV.SetAllInvisible(gameMap)
		playerFOV.RayCast(playerPosition.X, playerPosition.Y, gameMap)

		// Check if the current screen requires the ECS. If so, process all systems. If not, handle any input the screen
		// requires, and do not process and ECS systems
		if screenManager.CurrentScreen.UseEcs() {
			ecsController.Process(excludedSystems)
		} else {
			screenManager.CurrentScreen.HandleInput()
		}

		screenManager.CurrentScreen.Render()
		ecsController.ProcessSystem(reflect.TypeOf(SystemRender{}))

		ui.Refresh()
	}

	ui.CloseConsole()
}
