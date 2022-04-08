package game

import "bnd/engine"

type Direction struct {
	v engine.Vector2
}

func (d *Direction) Up() {
	d.v.Y += 1
}

func (d *Direction) Down() {
	d.v.Y += -1
}

func (d *Direction) Left() {
	d.v.X += -1
}

func (d *Direction) Right() {
	d.v.X += 1
}
