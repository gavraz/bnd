package engine

import (
	"fmt"
	"time"
)

type CollisionTypes int

const (
	Rectangle CollisionTypes = iota
	Circle
)

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
	ApplyFriction(dt float64)
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

type GameObject struct {
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
	Friction      float64
}

func (g *GameObject) GetIsPassthrough() bool {
	return g.IsPassthrough
}

func (g *GameObject) GetCollisionType() CollisionTypes {
	return g.CollisionType
}

func (g *GameObject) GetCenter() Vector2 {
	return g.Center
}

func (g *GameObject) GetHeight() float64 {
	return g.Height
}

func (g *GameObject) GetWidth() float64 {
	return g.Width
}

func (g *GameObject) SetCenter(p Vector2) {
	g.Center = p
}

func (g *GameObject) GetDirection() Vector2 {
	return g.Direction
}

func (g *GameObject) SetDirection(d Vector2) {
	g.Direction = d.Normalize()
}

func (g *GameObject) GetVelocity() Vector2 {
	return g.Velocity
}

func (g *GameObject) GetAcceleration() Vector2 {
	return g.Acceleration
}

func (g *GameObject) ForEachChild(do func(child Object)) {
	for _, child := range g.GetChildren() {
		do(child)
		child.ForEachChild(do)
	}
}

type Updater interface {
	Update(dt float64)
}

func (g *GameObject) Update(dt float64) {
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
		if !child.IsAlive() {
			g.RemoveChild(child)
			continue
		}

		if u, ok := child.(Updater); ok {
			u.Update(dt)
		}

		child.SetCenter(child.GetCenter().Sub(prevCenter).Add(g.Center))
	}
}

func (g *GameObject) UpdateTimeAlive(dt float64) {
	g.TimeAlive += dt
}

func (g *GameObject) SetVelocity(a Vector2) {
	g.Velocity = a
}

func (g *GameObject) SetAcceleration(a Vector2) {
	g.Acceleration = a
}

func (g *GameObject) GetBaseSpeed() float64 {
	return g.BaseSpeed
}

func (g *GameObject) SetBaseSpeed(s float64) {
	g.BaseSpeed = s
}

func (g *GameObject) GetMass() float64 {
	return g.Mass
}

func (g *GameObject) ApplyFriction(dt float64) {
	g.Velocity = g.Velocity.MulScalar(1 - g.Friction*dt)
}

func (g *GameObject) GetChildren() []DynamicObject {
	return g.ChildObjects
}

func (g *GameObject) SetChildren(children []DynamicObject) {
	g.ChildObjects = children
}

func (g *GameObject) AddChild(child DynamicObject) {
	g.ChildObjects = append(g.ChildObjects, child)
}

func (g *GameObject) SetParent(parent DynamicObject) {
	g.ParentObject = parent
}

func (g *GameObject) GetAppliedForce() Vector2 {
	return g.AppliedForce
}

func (g *GameObject) AddForce(force Vector2) {
	g.AppliedForce = g.AppliedForce.Add(force)
}

func (g *GameObject) GetParent() DynamicObject {
	return g.ParentObject
}

func (g *GameObject) IsAlive() bool {
	if g.Until.IsZero() {
		return true
	}
	return time.Now().Before(g.Until)
}

func (g *GameObject) RemoveChild(child DynamicObject) {
	for i, c := range g.ChildObjects {
		if c == child {
			g.ChildObjects = append(g.ChildObjects[:i], g.ChildObjects[i+1:]...)
			return
		}
	}
}

func (g *GameObject) RemoveParent() {
	g.ParentObject = nil
}

func (g *GameObject) GetHit() {
	if time.Now().After(g.HitCooldown) {
		fmt.Println("Hit!")
		g.HitCooldown = time.Now().Add(1000 * time.Millisecond)
	}
}
