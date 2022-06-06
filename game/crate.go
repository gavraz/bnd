package game

import (
	"bnd/engine"
)

type crate struct {
	engine.StaticObject
	removeCrate func(object engine.StaticObject)
}

func (c *crate) remove() {
	c.removeCrate(c)
}
