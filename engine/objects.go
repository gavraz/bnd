package engine

import (
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
	SetRootParent(g DynamicObject)
	GetRootParent() DynamicObject
	onCollision(collider Object)
}

type GameObjectConf struct {
	CollisionType CollisionTypes
	Width         float64
	Height        float64
	Center        Vector2
	Velocity      Vector2
	Acceleration  Vector2
	Direction     Vector2
	BaseSpeed     float64
	Mass          float64
	IsPassthrough bool
	Friction      float64
	Until         time.Time
}

type gameObject struct {
	childObjects []DynamicObject
	parentObject DynamicObject
	rootParent   DynamicObject
	appliedForce Vector2
	timeAlive    float64

	GameObjectConf
}

func NewDynamicObject(conf GameObjectConf) DynamicObject {
	obj := &gameObject{GameObjectConf: conf}
	obj.rootParent = obj
	return obj
}

func NewStaticObject(conf GameObjectConf) StaticObject {
	obj := &gameObject{GameObjectConf: conf}
	return obj
}

func (g *gameObject) GetIsPassthrough() bool {
	return g.IsPassthrough
}

func (g *gameObject) GetCollisionType() CollisionTypes {
	return g.CollisionType
}

func (g *gameObject) GetCenter() Vector2 {
	return g.Center
}

func (g *gameObject) GetHeight() float64 {
	return g.Height
}

func (g *gameObject) GetWidth() float64 {
	return g.Width
}

func (g *gameObject) SetCenter(p Vector2) {
	g.Center = p
}

func (g *gameObject) GetDirection() Vector2 {
	return g.Direction
}

func (g *gameObject) SetDirection(d Vector2) {
	g.Direction = d.Normalize()
}

func (g *gameObject) GetVelocity() Vector2 {
	return g.Velocity
}

func (g *gameObject) GetAcceleration() Vector2 {
	return g.Acceleration
}

func (g *gameObject) ForEachChild(do func(child Object)) {
	for _, child := range g.GetChildren() {
		do(child)
		child.ForEachChild(do)
	}
}

func (g *gameObject) onCollision(collider Object) {
	if d, ok := collider.(DynamicObject); ok {
		if g.CollisionType == Circle && collider.GetCollisionType() == Circle {
			g.onDynamicCollisionCircles(d)
		} else if g.CollisionType == Rectangle && collider.GetCollisionType() == Rectangle {
			g.onDynamicCollisionRectangles(d)
		} else if g.CollisionType == Circle && collider.GetCollisionType() == Rectangle {
			g.onDynamicCollisionCircleRectangle(d)
		} else if g.CollisionType == Rectangle && collider.GetCollisionType() == Circle {
			g.onDynamicCollisionRectangleCircle(d)
		}
	} else if s, ok := collider.(StaticObject); ok {
		if g.CollisionType == Circle && collider.GetCollisionType() == Circle {
			g.onStaticCollisionCircles(s)
		} else if g.CollisionType == Rectangle && collider.GetCollisionType() == Rectangle {
			g.onStaticCollisionRectangles(s)
		} else if g.CollisionType == Circle && collider.GetCollisionType() == Rectangle {
			g.onStaticCollisionCircleRectangle(s)
		} else if g.CollisionType == Rectangle && collider.GetCollisionType() == Circle {
			g.onStaticCollisionRectangleCircle(s)
		}
	}
}

type Updater interface {
	Update(dt float64)
}

func (g *gameObject) Update(dt float64) {
	g.Acceleration = g.appliedForce.DivScalar(g.Mass)
	prevCenter := g.Center
	g.Velocity = g.Velocity.Add(g.Acceleration.MulScalar(dt))
	g.Center = g.Center.Add(g.Velocity.MulScalar(dt))
	if g.Velocity.Length() != 0 {
		g.SetDirection(g.Velocity)
	}
	g.appliedForce = Vector2{}
	for _, child := range g.GetChildren() {
		child.UpdateTimeAlive(dt)
		if !child.IsAlive() {
			g.RemoveChild(child)
			continue
		}

		if u, ok := child.(Updater); ok {
			u.Update(dt)
		}
		child.SetCenter(g.GetCenter().Add(child.GetCenter().Sub(prevCenter)))
	}
}

func (g *gameObject) UpdateTimeAlive(dt float64) {
	g.timeAlive += dt
}

func (g *gameObject) SetVelocity(a Vector2) {
	g.Velocity = a
}

func (g *gameObject) SetAcceleration(a Vector2) {
	g.Acceleration = a
}

func (g *gameObject) GetBaseSpeed() float64 {
	return g.BaseSpeed
}

func (g *gameObject) SetBaseSpeed(s float64) {
	g.BaseSpeed = s
}

func (g *gameObject) GetMass() float64 {
	return g.Mass
}

func (g *gameObject) ApplyFriction(dt float64) {
	g.Velocity = g.Velocity.MulScalar(1 - g.Friction*dt)
}

func (g *gameObject) GetChildren() []DynamicObject {
	return g.childObjects
}

func (g *gameObject) SetChildren(children []DynamicObject) {
	g.childObjects = children
}

func (g *gameObject) AddChild(child DynamicObject) {
	child.SetRootParent(g.GetRootParent())
	child.SetParent(g)
	g.childObjects = append(g.childObjects, child)
}

func (g *gameObject) SetRootParent(parent DynamicObject) {
	g.rootParent = parent
}

func (g *gameObject) GetRootParent() DynamicObject {
	return g.rootParent
}

func (g *gameObject) SetParent(parent DynamicObject) {
	g.parentObject = parent
}

func (g *gameObject) GetAppliedForce() Vector2 {
	return g.appliedForce
}

func (g *gameObject) AddForce(force Vector2) {
	g.appliedForce = g.appliedForce.Add(force)
}

func (g *gameObject) GetParent() DynamicObject {
	return g.parentObject
}

func (g *gameObject) IsAlive() bool {
	if g.Until.IsZero() {
		return true
	}
	return time.Now().Before(g.Until)
}

func (g *gameObject) RemoveChild(child DynamicObject) {
	for i, c := range g.childObjects {
		if c == child {
			g.childObjects = append(g.childObjects[:i], g.childObjects[i+1:]...)
			return
		}
	}
}

func (g *gameObject) RemoveParent() {
	g.parentObject = nil
}
