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

type CollisionTypes int

const (
	Rectangle CollisionTypes = iota
	Circle
)

type Vector2 = v.Vector2

type Object interface {
	GetCollisionType() CollisionTypes
	GetCenter() Vector2
	SetCenter(p Vector2)
	GetWidth() float64
	GetHeight() float64
}

type StaticObject interface {
	Object
}

type DynamicObject interface {
	Object
	GetVelocity() Vector2
	SetVelocity(v Vector2)
	GetAcceleration() Vector2
	SetAcceleration(a Vector2)
	UpdateVelocity(dt float64)
	MoveObject()
	GetBaseSpeed() float64
	SetBaseSpeed(s float64)
	GetMass() float64
	ApplyFriction(friction, dt float64)
}

type GObject struct {
	CollisionType CollisionTypes
	Width         float64
	Height        float64
	Center        Vector2
	Velocity      Vector2
	Acceleration  Vector2
	Direction     Vector2
	BaseSpeed     float64
	Mass          float64
}

func (g *GObject) GetCollisionType() CollisionTypes {
	return g.CollisionType
}

func (g *GObject) GetCenter() Vector2 {
	return g.Center
}

func (g *GObject) GetHeight() float64 {
	return g.Height
}

func (g *GObject) GetWidth() float64 {
	return g.Width
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

func (g *GObject) MoveObject() {
	g.Center = g.Center.Add(g.Velocity)
}

func (g *GObject) SetVelocity(a Vector2) {
	g.Velocity = a
}

func (g *GObject) SetAcceleration(a Vector2) {
	g.Acceleration = a
}

func (g *GObject) GetBaseSpeed() float64 {
	return g.BaseSpeed
}

func (g *GObject) SetBaseSpeed(s float64) {
	g.BaseSpeed = s
}

func (g *GObject) GetMass() float64 {
	return g.Mass
}

func (g *GObject) ApplyFriction(friction, dt float64) {
	g.Velocity = g.Velocity.MulScalar(1 - friction*dt)
}

type Crate struct {
	DynamicObject

	ability Ability
}

type Wall struct {
	StaticObject
}

type Bullet struct {
	DynamicObject
}
