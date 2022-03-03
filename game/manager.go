package game

import (
	"fmt"
)

const (
	playerVelocityDecay = 3.0
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

func (m *Manager) AddDynamicObj(name string, object DynamicObject) {
	m.dynamicObjects[name] = object
}

func (m *Manager) AddStaticObj(name string, object StaticObject) {
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
	return m.dynamicObjects["current-player"].(*Player).hp
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
	m.AddDynamicObj("current-player", NewPlayer(&GObject{
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
	m.AddDynamicObj("crate2", &Crate{
		DynamicObject: &GObject{
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

func (m *Manager) DynamicObjects() map[string]DynamicObject {
	return m.dynamicObjects // TODO safety?
}

func (m *Manager) StaticObjects() map[string]StaticObject {
	return m.staticObjects
}

func (m *Manager) MovePlayer(dir Direction, dt float64) {
	playerObj := m.DynamicObjects()["current-player"]
	curSpeed := playerObj.GetBaseSpeed()
	playerObj.SetAcceleration(dirToVec2(dir).MulScalar(curSpeed * dt))
}
