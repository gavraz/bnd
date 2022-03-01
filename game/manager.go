package game

import (
	"fmt"
	"math"
)

const (
	playerVelocityDecay = 4.0
)

type Manager struct {
	objects map[string]Object
}

func NewManager() *Manager {
	return &Manager{
		objects: make(map[string]Object),
	}
}

func (m *Manager) SetDirection(name string, direction Vector2) {
	m.objects[name].SetDirection(direction)
}

func (m *Manager) Add(name string, object Object) {
	m.objects[name] = object
}

func (m *Manager) collidesWith(obj Object) Object {
	for _, other := range m.objects {
		if other == obj {
			continue
		}
		if collider := m.checkCollision(obj, other); collider != nil {
			return collider
		}
	}
	return nil
}

func (m *Manager) InitGame() {
	m.Add("current-player", NewPlayer(&GObject{
		Center: Vector2{
			X: 0,
			Y: 0,
		},
		BaseSpeed:     3,
		CollisionType: Circle,
		Width:         0.05,
		Height:        0.05,
		Mass:          100,
	}, 100))
	m.Add("crate", &Crate{
		Object: &GObject{
			Center: Vector2{
				X: -0.2,
				Y: -0.2,
			},
			CollisionType: Rectangle,
			Width:         0.1,
			Height:        0.1,
			Mass:          1,
		},
	})
	m.Add("crate2", &Crate{
		Object: &GObject{
			Center: Vector2{
				X: -0.3,
				Y: -0.3,
			},
			CollisionType: Rectangle,
			Width:         0.1,
			Height:        0.1,
			Mass:          1,
		},
	})

}

func (m *Manager) checkCollision(obj, other Object) Object {
	if obj.GetCollisionType() == Circle && other.GetCollisionType() == Circle {
		p1 := obj.GetCenter()
		p2 := other.GetCenter()
		r1 := obj.GetWidth() / 2
		r2 := other.GetWidth() / 2
		dist := p1.Distance(p2)
		if dist <= r1+r2 {
			return other
		}
	}
	if obj.GetCollisionType() == Circle && other.GetCollisionType() == Rectangle || obj.GetCollisionType() == Rectangle && other.GetCollisionType() == Circle {
		circle, rectangle := obj, other
		if obj.GetCollisionType() == Rectangle {
			circle, rectangle = other, obj
		}
		cx := circle.GetCenter().X
		cy := circle.GetCenter().Y
		rx := rectangle.GetCenter().X
		ry := rectangle.GetCenter().Y
		rw := rectangle.GetWidth() / 2
		rh := rectangle.GetHeight() / 2

		testX := cx
		testY := cy

		if cx < rx-rw {
			testX = rx - rw // left edge
		} else if cx > rx+rw {
			testX = rx + rw // right edge
		}
		if cy < ry-rh {
			testY = ry - rh // bottom edge
		} else if cy > ry+rh {
			testY = ry + rh // top edge
		}

		distX := cx - testX
		distY := cy - testY
		distanceSquared := (distX * distX) + (distY * distY)
		radiusSquared := (circle.GetWidth() / 2) * (circle.GetWidth() / 2)

		if distanceSquared <= radiusSquared {
			return other
		}
	}
	if obj.GetCollisionType() == Rectangle && other.GetCollisionType() == Rectangle {
		p1 := obj.GetCenter()
		p2 := other.GetCenter()
		w1 := obj.GetWidth() / 2
		h1 := obj.GetHeight() / 2
		w2 := other.GetWidth() / 2
		h2 := other.GetHeight() / 2
		if p1.X-w1 < p2.X+w2 && p1.X+w1 > p2.X-w2 && p1.Y-h1 < p2.Y+h2 && p1.Y+h1 > p2.Y-h2 {
			return other
		}
	}
	return nil
}

func (m *Manager) fixIntersection(obj Object, other Object) {
	if obj.GetCollisionType() == Circle && other.GetCollisionType() == Circle {
		p1 := obj.GetCenter()
		p2 := other.GetCenter()
		r1 := obj.GetWidth() / 2
		r2 := other.GetWidth() / 2
		penetrationDepth := r1 + r2 - p1.Distance(p2)
		direction := p1.Sub(p2).Normalize()
		obj.SetCenter(p1.Add(direction.MulScalar(penetrationDepth)))
	}

	if obj.GetCollisionType() == Rectangle && other.GetCollisionType() == Rectangle {
		p1 := obj.GetCenter()
		p2 := other.GetCenter()
		w1 := obj.GetWidth() / 2
		h1 := obj.GetHeight() / 2
		w2 := other.GetWidth() / 2
		h2 := other.GetHeight() / 2
		xDistSquared := (p1.X - p2.X) * (p1.X - p2.X)
		yDistSquared := (p1.Y - p2.Y) * (p1.Y - p2.Y)
		if xDistSquared > yDistSquared {
			if p1.X > p2.X {
				obj.SetCenter(Vector2{X: other.GetCenter().X + w2 + w1, Y: obj.GetCenter().Y})
			} else {
				obj.SetCenter(Vector2{X: other.GetCenter().X - w2 - w1, Y: obj.GetCenter().Y})
			}
		} else {
			if p1.Y > p2.Y {
				obj.SetCenter(Vector2{X: obj.GetCenter().X, Y: other.GetCenter().Y + h2 + h1})
			} else {
				obj.SetCenter(Vector2{X: obj.GetCenter().X, Y: other.GetCenter().Y - h2 - h1})
			}
		}
	}
	if obj.GetCollisionType() == Circle && other.GetCollisionType() == Rectangle {
		// https://stackoverflow.com/questions/45370692/circle-rectangle-collision-response
		circle := obj
		rect := other
		if obj.GetCollisionType() == Rectangle {
			circle, rect = rect, circle
		}

		NearestX := math.Max(rect.GetCenter().X-rect.GetWidth()/2, math.Min(circle.GetCenter().X, rect.GetCenter().X+rect.GetWidth()/2))
		NearestY := math.Max(rect.GetCenter().Y-rect.GetHeight()/2, math.Min(circle.GetCenter().Y, rect.GetCenter().Y+rect.GetHeight()/2))
		dist := Vector2{X: circle.GetCenter().X - NearestX, Y: circle.GetCenter().Y - NearestY}

		penetrationDepth := circle.GetWidth()/2 - dist.Length()
		penetrationVector := dist.Normalize().MulScalar(penetrationDepth)
		circle.SetCenter(circle.GetCenter().Add(penetrationVector.MulScalar(2)))
	}
}

func (m *Manager) Update(dt float64) {
	for _, g := range m.objects {
		g.UpdateVelocity(dt)
		if collider := m.collidesWith(g); collider != nil {
			fmt.Println("Collision detected: ", g.GetCenter(), collider.GetCenter())
			// Momentum and impulse calculations
			// https://www.gamasutra.com/view/feature/3015/pool_hall_lessons_fast_accurate_.php?page=3
			r1 := g.GetCenter()
			r2 := collider.GetCenter()
			n := r1.Sub(r2).Normalize()
			v1 := g.GetVelocity()
			v2 := collider.GetVelocity()
			a1 := v1.Dot(n)
			a2 := v2.Dot(n)
			// Using the optimized version,
			// optimizedP =  2(a1 - a2)
			//              -----------
			//                m1 + m2
			optimizedP := (2.0 * (a1 - a2)) / (g.GetMass() + collider.GetMass())
			u1 := v1.Sub(n.MulScalar(optimizedP * collider.GetMass()))
			u2 := v2.Add(n.MulScalar(optimizedP * g.GetMass()))
			g.SetVelocity(u1)
			collider.SetVelocity(u2)

			m.fixIntersection(g, collider)
		}
	}
	for _, g := range m.objects {
		g.SetCenter(g.GetCenter().Add(g.GetVelocity()))
		g.SetVelocity(Vector2{
			X: g.GetVelocity().X - g.GetVelocity().X*playerVelocityDecay*dt,
			Y: g.GetVelocity().Y - g.GetVelocity().Y*playerVelocityDecay*dt,
		})
	}
}

func (m *Manager) Objects() map[string]Object {
	return m.objects // TODO safety?
}
