package main

import (
	"github.com/gogue-framework/gogue/camera"
	"github.com/gogue-framework/gogue/data"
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

	// Data Loading
	dataLoader *data.FileLoader
	entityLoader *data.EntityLoader
	enemies map[string]interface{}
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
	playerBlocks := BlockingComponent{}
	player = ecsController.CreateEntity([]ecs.Component{playerPosition, playerAppearance, playerMovement, playerBlocks})

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

	// Data loading - load data from relevant data definition files
	dataLoader, _ = data.NewFileLoader("gamedata")
	entityLoader = data.NewEntityLoader(ecsController)
	loadData()
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
	ecsController.MapComponentClass("position", PositionComponent{})
	ecsController.MapComponentClass("appearance", AppearanceComponent{})
	ecsController.MapComponentClass("movement", MovementComponent{})
	ecsController.MapComponentClass("blocking", BlockingComponent{})
	ecsController.MapComponentClass("simpleai", SimpleAiComponent{})
}

func registerSystems() {
	render := SystemRender{ecsController: ecsController}
	input := SystemInput{ecsController: ecsController}
	simpleAi := SystemSimpleAi{ecsController: ecsController, mapSurface: gameMap}

	ecsController.AddSystem(input, 0)
	ecsController.AddSystem(simpleAi, 1)
	ecsController.AddSystem(render, 2)
}

// loadData loads game data (enemies definitions, item definitions, map progression data, etc) via Gogues data file and
// entity loader. Any entities loaded are stored in string indexed maps, making it easy to pull out and create an entity
// via its defined name
func loadData() {
	enemies, _ = dataLoader.LoadDataFromFile("enemies.json")
}

// placeEnemies places a handful of enemies at random around the level
func placeEnemies(numEnemies int) {
	for i := 0; i < numEnemies; i++ {
		var entityKeys []string
		// Build an index for each entity, so we can randomly choose one
		for key := range enemies["level_1"].(map[string]interface{}) {
			entityKeys = append(entityKeys, key)
		}

		// Now, randomly pick an entity
		entityKey := entityKeys[rng.Range(0, len(entityKeys))]
		entity := enemies["level_1"].(map[string]interface{})[entityKey].(map[string]interface{})

		// Create the entity based off the chosen item
		loadedEntity := entityLoader.CreateSingleEntity(entity)

		randomTile := gameMap.FloorTiles[rng.Range(0, len(gameMap.FloorTiles))]
		ecsController.UpdateComponent(loadedEntity, PositionComponent{}.TypeOf(), PositionComponent{X: randomTile.X, Y: randomTile.Y})
	}
}


func main() {
	// Register all screens
	registerScreens()

	// Initialize the game map
	gameMap = SetupGameMap()
	placeEnemies(100)

	registerSystems()
	excludedSystems := []reflect.Type{reflect.TypeOf(SystemRender{})}

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
