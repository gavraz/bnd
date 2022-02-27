package game

import "fmt"

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
		x1 := obj.GetCenter().X
		y1 := obj.GetCenter().Y
		r1 := obj.Radius()
		r2 := obj.Radius()
		x2 := other.GetCenter().X
		y2 := other.GetCenter().Y

		actualDistSquared := (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
		radiusSumSquared := (r1 + r2) * (r1 + r2)
		if actualDistSquared <= radiusSumSquared {
			return other
		}
	}
	return nil
}

func (m *Manager) Update(dt float64) { // TODO: param of time
	for _, g := range m.objects {
		if collider := m.collidesWith(g); collider != nil {
			fmt.Println("Collision detected: ", g, collider)
		}

		g.UpdateVelocity(dt)
		g.SetCenter(Vector2{
			X: g.GetCenter().X + g.GetVelocity().X,
			Y: g.GetCenter().Y + g.GetVelocity().Y,
		})
		g.SetVelocity(Vector2{
			X: g.GetVelocity().X - g.GetVelocity().X*playerVelocityDecay*dt,
			Y: g.GetVelocity().Y - g.GetVelocity().Y*playerVelocityDecay*dt,
		})
	}
}

func (m *Manager) Objects() map[string]Object {
	return m.objects // TODO safety?
}
