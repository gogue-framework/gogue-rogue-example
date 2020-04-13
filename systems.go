package main

import (
	"github.com/gogue-framework/gogue/ecs"
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

			// Clear the cell this entity occupies, so it is the only glyph drawn there
			for i := 0; i <= 2; i++ {
				ui.ClearArea(pos.X, pos.Y, 1, 1, i)
			}
			ui.PrintGlyph(pos.X, pos.Y, appearance.Glyph, "", appearance.Layer)
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

			newX = dx + pos.X
			newY = dy + pos.Y

			ui.PrintGlyph(pos.X, pos.Y, ui.EmptyGlyph, "", 1)

			newPos := PositionComponent{X: newX, Y: newY}
			ecsController.UpdateComponent(player, PositionComponent{}.TypeOf(), newPos)
		}
	}
}
