package main

import (
	"github.com/gogue-framework/gogue/ecs"
	"github.com/gogue-framework/gogue/gamemap"
)

// GetBlockingEntities returns true if there is an entity at the location, and that entity has the Blocking component
func GetBlockingEntities(x, y int, entityController *ecs.Controller) bool {
	for _, entity := range entityController.GetEntitiesWithComponent(BlockingComponent{}.TypeOf()) {
		entityPosition := entityController.GetComponent(entity, PositionComponent{}.TypeOf()).(PositionComponent)

		if entityPosition.X == x && entityPosition.Y == y {
			return true
		}
	}

	return false
}

// Debug Utils
// MakeAllVisible makes all map Tiles visible to the player
func MakeAllVisible(mapSurface *gamemap.GameMap) {
	for x := 0; x < mapSurface.Width; x++ {
		for y := 0; y < mapSurface.Height; y++ {
			mapSurface.Tiles[x][y].Visible = true
		}
	}
}