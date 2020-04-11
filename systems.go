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
