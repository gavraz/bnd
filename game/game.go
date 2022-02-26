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
	GetDirection() Direction
	GetVelocity() float32
}

type Point struct {
	X, Y float32
}

type Direction struct {
	DX, DY float32
}

type GObject struct {
	Point
	Velocity float32
	Direction
}

func (o GObject) GetPoint() Point {
	return o.Point
}

func (o GObject) SetPoint(p Point) {
	o.Point = p
}

func (o GObject) GetDirection() Direction {
	return o.Direction
}

func (o GObject) GetVelocity() float32 {
	return o.Velocity
}

type Crate struct {
	GObject

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
