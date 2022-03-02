package game

type Direction struct {
	up    bool
	down  bool
	left  bool
	right bool
}

func (d *Direction) Up() {
	d.up = true
}

func (d *Direction) Down() {
	d.down = true
}

func (d *Direction) Left() {
	d.left = true
}

func (d *Direction) Right() {
	d.right = true
}

func dirToVec2(d Direction) Vector2 {
	var x, y float64
	if d.up {
		y++
	}
	if d.down {
		y--
	}
	if d.left {
		x--
	}
	if d.right {
		x++
	}
	return Vector2{X: x, Y: y}
}
