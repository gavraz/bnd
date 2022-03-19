package game

import (
	//v "bnd/vector_pointers"
	v "bnd/vector"
	"fmt"
	"time"
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
	GetBaseSpeed() float64
	SetBaseSpeed(s float64)
	GetMass() float64
	ApplyFriction(friction, dt float64)
	SetDirection(dir Vector2)
	GetDirection() Vector2
	GetAppliedForce() Vector2
	AddForce(force Vector2)
	UpdateTimeAlive(dt float64)
	ForEachChild(do func(child Object))
	GetChildren() []DynamicObject
	SetChildren(children []DynamicObject)
	AddChild(child DynamicObject)
	SetParent(parent DynamicObject)
	GetParent() DynamicObject
	RemoveParent()
	RemoveChild(child DynamicObject)
	IsAlive() bool
	GetHit()
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
	Until         time.Time
	HitCooldown   time.Time
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

func (g *GObject) ForEachChild(do func(child Object)) {
	for _, child := range g.GetChildren() {
		do(child)
		child.ForEachChild(do)
	}
}

func (g *GObject) Update(dt float64) {
	g.Acceleration = g.AppliedForce.DivScalar(g.Mass)
	g.Velocity = g.Velocity.Add(g.Acceleration.MulScalar(dt))
	prevCenter := g.Center
	g.Center = g.Center.Add(g.Velocity.MulScalar(dt))
	if g.Velocity.Length() != 0 {
		g.SetDirection(g.Velocity)
	}
	g.AppliedForce = Vector2{}
	for _, child := range g.GetChildren() {
		child.UpdateTimeAlive(dt)
		if ObjectType(child) == Melee {
			child.(*meleeObject).update(dt)
		}
		child.SetCenter(child.GetCenter().Sub(prevCenter).Add(g.Center))
		if !child.IsAlive() {
			g.RemoveChild(child)
			continue
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

func (g *GObject) IsAlive() bool {
	if g.Until.IsZero() {
		return true
	}
	return time.Now().Before(g.Until)
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

func (g *GObject) GetHit() {
	if time.Now().After(g.HitCooldown) {
		fmt.Println("Hit!")
		g.HitCooldown = time.Now().Add(1000 * time.Millisecond)
	}
}
