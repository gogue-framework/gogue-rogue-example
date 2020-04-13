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

func (pc PositionComponent) TypeOf() reflect.Type { return reflect.TypeOf(pc) }

// AppearanceComponent defines what an entity looks like when rendered onto the terminal. It contains a glyph character
// representation, the layer it should be rendered on, and a name and description. An entity must have an appearance
// component to be rendered.
type AppearanceComponent struct {
	Glyph       ui.Glyph
	Layer       int
	Name        string
	Description string
}

func (ac AppearanceComponent) TypeOf() reflect.Type { return reflect.TypeOf(ac) }

// MovementComponent is a flag component (no metadata). Its presence on an entity indicates that the entity can freely move
type MovementComponent struct{}

func (mc MovementComponent) TypeOf() reflect.Type { return reflect.TypeOf(mc) }
