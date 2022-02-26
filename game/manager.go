package game

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

func (m *Manager) Update() { // TODO: param of time
	for _, g := range m.objects {
		g.UpdateVelocity()
		g.SetPoint(Point{
			X: g.GetPoint().X + g.GetVelocity().X,
			Y: g.GetPoint().Y + g.GetVelocity().Y,
		})
		g.SetVelocity(Vector2{X: g.GetVelocity().X * 0.95, Y: g.GetVelocity().Y * 0.95})
	}
}

func (m *Manager) Objects() map[string]Object {
	return m.objects // TODO safety?
}
