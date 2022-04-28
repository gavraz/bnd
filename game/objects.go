package game

import (
	"bnd/engine"
	"fmt"
)

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

func (p *player) applyDamage(damage int) {
	p.hp -= damage
	fmt.Println("Hit! \nCurrent hp: ", p.hp)
}
