package main

import (
	"github.com/gogue-framework/gogue/ecs"
	"github.com/gogue-framework/gogue/gamemap"
	"github.com/gogue-framework/gogue/ui"
)

type SystemRender struct {
	ecsController *ecs.Controller
}

func (sr SystemRender) Process() {
	// Render all entities present in the global entity controller
	for e := range sr.ecsController.GetEntities() {
		if sr.ecsController.HasComponent(e, PositionComponent{}.TypeOf()) && sr.ecsController.HasComponent(e, AppearanceComponent{}.TypeOf()) {

			pos := sr.ecsController.GetComponent(e, PositionComponent{}.TypeOf()).(PositionComponent)
			appearance := sr.ecsController.GetComponent(e, AppearanceComponent{}.TypeOf()).(AppearanceComponent)

			// Get the coordinates of the entity, in reference to where the camera currently is.
			cameraX, cameraY := gameCamera.ToCameraCoordinates(pos.X, pos.Y)

			if gameMap.IsVisibleToPlayer(pos.X, pos.Y) {
				// Clear the cell this entity occupies, so it is the only glyph drawn there
				for i := 1; i <= 2; i++ {
					ui.ClearArea(cameraX, cameraY, 1, 1, i)
				}
				ui.PrintGlyph(cameraX, cameraY, appearance.Glyph, "", appearance.Layer)
			}
		}
	}
}

type SystemInput struct {
	ecsController *ecs.Controller
}

func (si SystemInput) Process() {
	// Get all entities that can be controlled by the player. Most of the time, this will just be the player entity, but
	// it may be possible to take control of other entities. This system will block until the player has taken an action
	key := ui.ReadInput(false)

	if key == ui.KeyEscape {
		// Shift the screen back to the title screen
		_ = screenManager.SetScreenByName("title")

		// Render the screen
		screenManager.CurrentScreen.Render()

		return
	}

	if key == ui.KeyZ {
		// Make all Tiles visible
		MakeAllVisible(gameMap)
	}

	// Fow now, just handle movement
	if ecsController.HasComponent(player, PositionComponent{}.TypeOf()) && ecsController.HasComponent(player, AppearanceComponent{}.TypeOf()) {

		// Check if the player is currently capable of movement (has a movement component)
		if ecsController.HasComponent(player, MovementComponent{}.TypeOf()) {
			pos := ecsController.GetComponent(player, PositionComponent{}.TypeOf()).(PositionComponent)

			// Clear the existing appearance from the screen, since it will be moving. This will prevent artifact trails.
			ui.PrintGlyph(pos.X, pos.Y, ui.EmptyGlyph, "", 1)

			var dx, dy, newX, newY int

			switch key {
			case ui.KeyRight, ui.KeyL:
				dx, dy = 1, 0
			case ui.KeyLeft, ui.KeyH:
				dx, dy = -1, 0
			case ui.KeyUp, ui.KeyK:
				dx, dy = 0, -1
			case ui.KeyDown, ui.KeyJ:
				dx, dy = 0, 1
			}

			if !gameMap.IsBlocked(pos.X + dx, pos.Y + dy) && !GetBlockingEntities(pos.X + dx, pos.Y + dy, si.ecsController) {
				newX = dx + pos.X
				newY = dy + pos.Y

				cameraX, cameraY := gameCamera.ToCameraCoordinates(pos.X, pos.Y)
				ui.PrintGlyph(cameraX, cameraY, ui.EmptyGlyph, "", 1)

				newPos := PositionComponent{X: newX, Y: newY}
				si.ecsController.UpdateComponent(player, PositionComponent{}.TypeOf(), newPos)
			}
		}
	}
}

type SystemSimpleAi struct {
	ecsController *ecs.Controller
	mapSurface    *gamemap.GameMap
}

func (sas SystemSimpleAi) Process() {
	// Process all entities that have the simple AI component attached to them
	// For now, just have them print something
	for _, entity := range sas.ecsController.GetEntitiesWithComponent(SimpleAiComponent{}.TypeOf()) {
		//Act
		if sas.ecsController.HasComponent(entity, AppearanceComponent{}.TypeOf()) {
			// For the time being, just have the AI move around randomly. This will be fleshed out in time.
			pos := sas.ecsController.GetComponent(entity, PositionComponent{}.TypeOf()).(PositionComponent)
			dx := rng.RangeNegative(-1, 1)
			dy := rng.RangeNegative(-1, 1)

			var newX, newY int

			if !sas.mapSurface.IsBlocked(pos.X + dx, pos.Y + dy) && !GetBlockingEntities(pos.X + dx, pos.Y + dy, sas.ecsController){
				newX = dx + pos.X
				newY = dy + pos.Y
			} else {
				newX = pos.X
				newY = pos.Y
			}

			cameraX, cameraY := gameCamera.ToCameraCoordinates(pos.X, pos.Y)
			ui.PrintGlyph(cameraX, cameraY, ui.EmptyGlyph, "", 1)


			newPos := PositionComponent{X: newX, Y: newY}
			sas.ecsController.UpdateComponent(entity, PositionComponent{}.TypeOf(), newPos)
		}
	}
}
