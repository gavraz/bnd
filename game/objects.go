package game

import (
	"bnd/engine"
	"time"
)

type player struct {
	engine.DynamicObject
	hitCooldown time.Time
	hp          int
}

type crate struct {
	engine.DynamicObject
}

type wall struct {
	engine.StaticObject
}
