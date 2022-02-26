package game

type Ability int

const (
	None Ability = iota
	Fire
	Speed
	Bomb
	Oil
	Shield
)

type Object interface {
	GetPoint() Point
	SetPoint(p Point)
	GetDirection() Vector2
	SetDirection(d Vector2)
	GetVelocity() Vector2
	SetVelocity(v Vector2)
	UpdateVelocity()
	GetAcceleration() Vector2
	SetAcceleration(a Vector2)
}

type Point struct {
	X, Y float32
}

type Vector2 struct {
	X, Y float32
}

type GObject struct {
	Point
	Velocity     Vector2
	Acceleration Vector2
	Direction    Vector2
}

func (o *GObject) GetPoint() Point {
	return o.Point
}

func (o *GObject) SetPoint(p Point) {
	o.Point = p
}

func (o *GObject) GetDirection() Vector2 {
	return o.Direction
}

func (o *GObject) SetDirection(d Vector2) {
	o.Direction = d
}

func (o *GObject) GetVelocity() Vector2 {
	return o.Velocity
}

func (o *GObject) GetAcceleration() Vector2 {
	return o.Acceleration
}

func (g *GObject) UpdateVelocity() {
	g.Velocity.X += g.Acceleration.X
	g.Velocity.Y += g.Acceleration.Y
}

func (o *GObject) SetVelocity(a Vector2) {
	o.Velocity = a
}

func (g *GObject) SetAcceleration(a Vector2) {
	g.Acceleration = a
}

type Crate struct {
	Object

	ability Ability
}

type Player struct {
	Object

	hp        int
	primary   Ability
	secondary Ability
}

func NewPlayer(object Object, hp int) *Player {
	return &Player{
		Object: object,
		hp:     hp,
	}
}

type Bullet struct {
	GObject
}
