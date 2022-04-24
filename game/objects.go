package game

import "bnd/engine"

type player struct {
	engine.DynamicObject
	hp int
}

type crate struct {
	engine.DynamicObject
}

type wall struct {
	engine.StaticObject
}
