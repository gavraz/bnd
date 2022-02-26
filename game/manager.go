package game

type Manager struct {
	objects map[string]Object
}

func NewManager() *Manager {
	return &Manager{
		objects: make(map[string]Object),
	}
}

func (m *Manager) SetDirection(name string, direction Direction) {
	m.objects[name].SetDirection(direction)
}

func (m *Manager) Add(name string, object Object) {
	m.objects[name] = object
}

func (m *Manager) Update() { // TODO: param of time
	for _, g := range m.objects {
		g.SetPoint(Point{
			X: g.GetPoint().X + g.GetDirection().DX*g.GetVelocity(),
			Y: g.GetPoint().Y + g.GetDirection().DY*g.GetVelocity(),
		})
	}
}

func (m *Manager) Objects() map[string]Object {
	return m.objects // TODO safety?
}
