package game

type Manager struct {
	objects []Object
}

func (m *Manager) Add(object Object) {
	m.objects = append(m.objects, object)
}

func (m *Manager) Update() { // TODO: param of time
	for _, g := range m.objects {
		g.SetPoint(Point{
			X: g.GetPoint().X + g.GetDirection().DX*g.GetVelocity(),
			Y: g.GetPoint().Y + g.GetDirection().DY*g.GetVelocity(),
		})
	}
}

func (m *Manager) Objects() []Object {
	return m.objects // TODO safety?
}
