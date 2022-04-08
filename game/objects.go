package game

import "bnd/engine"

type player struct {
	engine.DynamicObject

	hp int
}

func NewPlayer(object engine.DynamicObject, hp int) *player {
	return &player{
		DynamicObject: object,
		hp:            hp,
	}
}

type crate struct {
	engine.DynamicObject
}

type wall struct {
	engine.StaticObject
}

type fart struct {
	engine.DynamicObject
}
