package game

import (
	"fmt"
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
		p1 := obj.GetCenter()
		p2 := other.GetCenter()
		r1 := obj.Radius()
		r2 := other.Radius()
		dist := p1.Distance(p2)
		if dist <= r1+r2 {
			return other
		}
	}
	return nil
}

func (m *Manager) Update(dt float64) { // TODO: param of time
	//l := list.New()
	for _, g := range m.objects {
		g.UpdateVelocity(dt)
		if collider := m.collidesWith(g); collider != nil {
			//fmt.Println("Collision detected: ", g, collider)
			fmt.Println("Collision detected: ", g.GetCenter(), collider.GetCenter())
			// Momentum and impulse calculations
			// http://www.gamasutra.com/view/feature/131790/pool_hall_lessons_fast_accurate_.php?page=3
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
			optimizedP := (2.0 * (a1 - a2)) / 2
			u1 := v1.Sub(n.MulScalar(optimizedP))
			u2 := v2.Add(n.MulScalar(optimizedP))
			g.SetVelocity(u1)
			collider.SetVelocity(u2)
			g.SetCenter(collider.GetCenter().Add(r1.Sub(r2).Normalize().MulScalar(g.Radius() + collider.Radius() + 1)))
			collider.SetCenter(g.GetCenter().Add(r2.Sub(r1).Normalize().MulScalar(g.Radius() + collider.Radius() + 1)))

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
