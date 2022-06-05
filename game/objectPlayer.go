package game

import (
	"bnd/engine"
	"fmt"
)

type player struct {
	engine.DynamicObject
	hp int
}

func (p *player) applyDamage(damage int) {
	p.hp -= damage
	fmt.Println("Hit! \nCurrent hp: ", p.hp)
}

func (p *player) OnCollision(collider engine.Object) {
	if c, ok := collider.(*crate); ok {
		c.remove()
	}
}
