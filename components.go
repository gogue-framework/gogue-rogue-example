package main

import (
	"github.com/gogue-framework/gogue/ui"
	"reflect"
)

// PositionComponent represents a location on the GameMap. An entity must have a PositionComponent to be rendered.
type PositionComponent struct {
	X int
	Y int
}

func (ps PositionComponent) TypeOf() reflect.Type { return reflect.TypeOf(ps)}

// AppearanceComponent defines what an entity looks like when rendered onto the terminal. It contains a glyph character
// representation, the layer it should be rendered on, and a name and description. An entity must have an appearance
// component to be rendered.
type AppearanceComponent struct {
	Glyph ui.Glyph
	Layer int
	Name string
	Description string
}

func (as AppearanceComponent) TypeOf() reflect.Type { return reflect.TypeOf(as) }
