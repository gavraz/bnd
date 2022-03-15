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
	GetIsPassthrough() bool
	GetVelocity() Vector2
	SetVelocity(v Vector2)
	GetAcceleration() Vector2
	SetAcceleration(a Vector2)
	Update(dt float64)
	MoveObject(dt float64)
	GetBaseSpeed() float64
	SetBaseSpeed(s float64)
	GetMass() float64
	ApplyFriction(friction, dt float64)
	SetDirection(dir Vector2)
	GetDirection() Vector2
	GetAppliedForce() Vector2
	AddForce(force Vector2)
	isDead() bool
	UpdateTimeAlive(dt float64)
	GetChildren() []DynamicObject
	SetChildren(children []DynamicObject)
	AddChild(child DynamicObject)
	SetParent(parent DynamicObject)
	GetParent() DynamicObject
	RemoveParent()
	RemoveChild(child DynamicObject)
}

type GObject struct {
	CollisionType CollisionTypes
	ChildObjects  []DynamicObject
	ParentObject  DynamicObject
	AppliedForce  Vector2
	Width         float64
	Height        float64
	Center        Vector2
	Velocity      Vector2
	Acceleration  Vector2
	Direction     Vector2
	TimeAlive     float64
	TimeToLive    float64
	BaseSpeed     float64
	Mass          float64
	IsPassthrough bool
}

func (g *GObject) GetIsPassthrough() bool {
	return g.IsPassthrough
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
	g.Direction = d.Normalize()
}

func (g *GObject) GetVelocity() Vector2 {
	return g.Velocity
}

func (g *GObject) GetAcceleration() Vector2 {
	return g.Acceleration
}

func (g *GObject) Update(dt float64) {
	g.Acceleration = g.AppliedForce.DivScalar(g.Mass)
	g.Velocity = g.Velocity.Add(g.Acceleration.MulScalar(dt))
	g.Center = g.Center.Add(g.Velocity.MulScalar(dt))
	g.SetDirection(g.Velocity)
	g.AppliedForce = Vector2{}

	if g.GetChildren() != nil {
		for _, child := range g.GetChildren() {
			if ObjectType(child) == MeleeObject {
				child.(*meleeObject).UpdateMelee(dt)
			} else {
				child.SetCenter(g.GetCenter()) // Only applies for fart atm
			}
			child.UpdateTimeAlive(dt)
			if child.isDead() {
				g.RemoveChild(child)
				continue
			}

		}
	}
}

func (g *GObject) MoveObject(dt float64) {
	//g.Center = g.Center.Add(g.Velocity.MulScalar(dt).Add(g.Acceleration.MulScalar(dt * dt / 2)))
	moveVector := g.Velocity.MulScalar(dt)
	prevCenter := g.Center
	g.Center = g.Center.Add(moveVector)
	children := g.GetChildren()
	if children != nil {
		for _, child := range g.GetChildren() {
			child.SetCenter(child.GetCenter().Sub(prevCenter).Add(g.Center))
		}
	}
}

func (g *GObject) UpdateTimeAlive(dt float64) {
	g.TimeAlive += dt
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

func (g *GObject) GetChildren() []DynamicObject {
	return g.ChildObjects
}

func (g *GObject) SetChildren(children []DynamicObject) {
	g.ChildObjects = children
}

func (g *GObject) AddChild(child DynamicObject) {
	g.ChildObjects = append(g.ChildObjects, child)
}

func (g *GObject) SetParent(parent DynamicObject) {
	g.ParentObject = parent
}

func (g *GObject) GetAppliedForce() Vector2 {
	return g.AppliedForce
}

func (g *GObject) AddForce(force Vector2) {
	g.AppliedForce = g.AppliedForce.Add(force)
}

func (g *GObject) GetParent() DynamicObject {
	return g.ParentObject
}

func (g *GObject) isDead() bool {
	return g.TimeAlive > g.TimeToLive
}

func (g *GObject) RemoveChild(child DynamicObject) {
	for i, c := range g.ChildObjects {
		if c == child {
			g.ChildObjects = append(g.ChildObjects[:i], g.ChildObjects[i+1:]...)
			return
		}
	}
}

func (g *GObject) RemoveParent() {
	g.ParentObject = nil
}
