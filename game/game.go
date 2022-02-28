package game

import (
	//v "bnd/vector_pointers"
	v "bnd/vector"
)

type Ability int

const (
	None Ability = iota
	Fire
	Speed
	Bomb
	Oil
	Shield
)

type Vector2 = v.Vector2

type Object interface {
	Radius() float64
	GetCenter() Vector2
	SetCenter(p Vector2)
	GetDirection() Vector2
	SetDirection(d Vector2)
	GetVelocity() Vector2
	SetVelocity(v Vector2)
	UpdateVelocity(dt float64)
	GetAcceleration() Vector2
	SetAcceleration(a Vector2)
}

type GObject struct {
	CollisionRadius float64
	Center          Vector2
	Velocity        Vector2
	Acceleration    Vector2
	Direction       Vector2
}

func (g *GObject) Radius() float64 {
	return g.CollisionRadius
}

func (g *GObject) GetCenter() Vector2 {
	return g.Center
}

func (g *GObject) SetCenter(p Vector2) {
	g.Center = p
}

func (g *GObject) GetDirection() Vector2 {
	return g.Direction
}

func (g *GObject) SetDirection(d Vector2) {
	g.Direction = d
}

func (g *GObject) GetVelocity() Vector2 {
	return g.Velocity
}

func (g *GObject) GetAcceleration() Vector2 {
	return g.Acceleration
}

func (g *GObject) UpdateVelocity(dt float64) {
	g.Velocity = g.Velocity.Add(g.Acceleration.MulScalar(dt))
}

func (g *GObject) SetVelocity(a Vector2) {
	g.Velocity = a
}

func (g *GObject) SetAcceleration(a Vector2) {
	g.Acceleration = a
}

type Crate struct {
	Object

	ability Ability
}

type Bullet struct {
	GObject
}

type BouncingBall struct {
	Object
}
