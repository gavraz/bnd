package game

import (
	"bnd/engine"
	"fmt"
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

func (p *player) getHit(damage int) {
	if time.Now().After(p.hitCooldown) {
		p.hp -= damage
		fmt.Println("Hit! \nCurrent hp: ", p.hp)
		p.hitCooldown = time.Now().Add(1000 * time.Millisecond)
	}
}
