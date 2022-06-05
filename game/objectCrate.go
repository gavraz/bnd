package game

import (
	"bnd/engine"
	"fmt"
	"math/rand"
)

type crate struct {
	engine.StaticObject
	removeCrate func(object engine.StaticObject)
}

func spawnCreate(name string, removeCrate func(object engine.StaticObject)) *crate {
	pos := engine.Vector2{rand.Float64()*1.8 - 0.9, rand.Float64()*1.5 - 0.6}
	fmt.Println(pos)
	crateSize := 0.05
	return &crate{engine.NewStaticObject(engine.GameObjectConf{
		Name:          name,
		Center:        pos,
		Width:         crateSize,
		Height:        crateSize,
		IsPassthrough: true,
	}),
		removeCrate}
}

func (c *crate) remove() {
	c.removeCrate(c)
}
