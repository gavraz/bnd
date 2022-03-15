package game

type MoveDirection struct {
	v Vector2
}

func (d *MoveDirection) Up() {
	d.v.Y += 1
}

func (d *MoveDirection) Down() {
	d.v.Y += -1
}

func (d *MoveDirection) Left() {
	d.v.X += -1
}

func (d *MoveDirection) Right() {
	d.v.X += 1
}
