package game

import "fmt"

const (
	playerVelocityDecay = 4.0
)

type Manager struct {
	dynamicObjects map[string]DynamicObject
	staticObjects  map[string]StaticObject
}

func NewManager() *Manager {
	return &Manager{
		dynamicObjects: make(map[string]DynamicObject),
		staticObjects:  make(map[string]StaticObject),
	}
}

func (m *Manager) AddDynamicObject(name string, object DynamicObject) {
	m.dynamicObjects[name] = object
}

func (m *Manager) AddStaticObject(name string, object StaticObject) {
	m.staticObjects[name] = object
}

func (m *Manager) ForEachGameObject(do func(object Object)) {
	for _, obj := range m.dynamicObjects {
		do(obj)
	}
	for _, obj := range m.staticObjects {
		do(obj)
	}
}

func (m *Manager) HP() int {
	return m.dynamicObjects["current-player"].(*player).hp
}

func (m *Manager) resolveDynamicCollisions(obj DynamicObject) Object {
	for _, other := range m.dynamicObjects {
		if other == obj {
			continue
		}
		if collider := CheckDynamicCollision(obj, other); collider != nil {
			return collider
		}
	}
	return nil
}

func (m *Manager) resolveStaticCollisions(obj DynamicObject) Object {
	for _, other := range m.staticObjects {
		if collider := CheckStaticCollision(obj, other); collider != nil {
			return collider
		}
	}
	return nil
}

func (m *Manager) InitGame() {
	m.AddDynamicObject("current-player", NewPlayer(&GObject{
		Center: Vector2{
			X: 0,
			Y: 0,
		},
		BaseSpeed:     3,
		CollisionType: Circle,
		Width:         0.05,
		Height:        0.05,
		Mass:          2,
	}, 100))
	m.AddStaticObject("enemy-player", NewPlayer(&GObject{
		Center: Vector2{
			X: 0.2,
			Y: 0.3,
		},
		BaseSpeed:     3,
		CollisionType: Circle,
		Width:         0.2,
		Height:        0.2,
	}, 100))
	m.AddDynamicObject("crate", &crate{
		DynamicObject: &GObject{
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
	m.AddDynamicObject("crate2", &crate{
		DynamicObject: &GObject{
			Center: Vector2{
				X: -0.6,
				Y: -0.5,
			},
			CollisionType: Rectangle,
			Width:         0.2,
			Height:        0.2,
			Mass:          10,
		},
	})
	m.AddStaticObject("wall-bottom", &wall{
		StaticObject: &GObject{
			Center: Vector2{
				X: 0,
				Y: -0.83,
			},
			CollisionType: Rectangle,
			Width:         1.92,
			Height:        0.34,
		},
	})
	m.AddStaticObject("wall-left", &wall{
		StaticObject: &GObject{
			Center: Vector2{
				X: -0.98,
				Y: 0,
			},
			CollisionType: Rectangle,
			Width:         0.04,
			Height:        2,
		},
	})
	m.AddStaticObject("wall-right", &wall{
		StaticObject: &GObject{
			Center: Vector2{
				X: 0.98,
				Y: 0,
			},
			CollisionType: Rectangle,
			Width:         0.04,
			Height:        2,
		},
	})
	m.AddStaticObject("wall-top", &wall{
		StaticObject: &GObject{
			Center: Vector2{
				X: 0,
				Y: 0.98,
			},
			CollisionType: Rectangle,
			Width:         1.92,
			Height:        0.04,
		},
	})

}

func (m *Manager) Update(dt float64) {
	for _, obj := range m.dynamicObjects {
		obj.UpdateVelocity(dt)
		obj.ApplyFriction(playerVelocityDecay, dt)
		obj.MoveObject()

		if collider := m.resolveDynamicCollisions(obj); collider != nil {
			fmt.Println("Dynamic Collision detected: ", obj.GetCenter(), collider.GetCenter())
		}
		if collider := m.resolveStaticCollisions(obj); collider != nil {
			fmt.Println("Static Collision detected: ", obj.GetCenter(), collider.GetCenter())
		}
	}
}

func (m *Manager) MovePlayer(dir Direction, dt float64) {
	playerObj := m.dynamicObjects["current-player"]
	curSpeed := playerObj.GetBaseSpeed()
	playerObj.SetAcceleration(dirToVec2(dir).MulScalar(curSpeed * dt))
}
