package game

import (
	"fmt"
)

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
		if obj.GetChildren() != nil {
			for _, child := range obj.GetChildren() {
				do(child)
			}
		}
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
		for _, child := range obj.GetChildren() {
			if Object(other) == child {
				goto IsChild
			}
		}
		for _, child := range other.GetChildren() {
			if Object(obj) == child {
				goto IsChild
			}
		}
		if other == obj {
			continue
		}
		if collider := CheckDynamicCollision(obj, other); collider != nil {
			return collider
		}
	IsChild:
	}
	return nil
}

func (m *Manager) resolveStaticCollisions(obj DynamicObject) Object {
	for _, other := range m.staticObjects {
		for _, child := range obj.GetChildren() {
			if Object(other) == child {
				goto IsChild
			}
		}
		if collider := CheckStaticCollision(obj, other); collider != nil {
			return collider
		}
	IsChild:
	}
	return nil
}

func (m *Manager) InitGame() {
	m.AddDynamicObject("current-player", NewPlayer(&GObject{
		Center: Vector2{
			X: 0,
			Y: 0,
		},
		BaseSpeed:     5,
		CollisionType: Circle,
		Width:         0.05,
		Height:        0.05,
		Mass:          1,
	}, 100))
	m.AddStaticObject("enemy-player", NewPlayer(&GObject{
		Center: Vector2{
			X: 0.2,
			Y: 0.3,
		},
		BaseSpeed:     1,
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
			Mass:          2,
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
		if obj.GetParent() != nil {
			continue
		}
		obj.ApplyFriction(playerVelocityDecay, dt)
		obj.Update(dt)

		if collider := m.resolveDynamicCollisions(obj); collider != nil {
			fmt.Println("Dynamic Collision detected: ", obj.GetCenter(), collider.GetCenter())
		}
		if collider := m.resolveStaticCollisions(obj); collider != nil {
			fmt.Println("Static Collision detected: ", obj.GetCenter(), collider.GetCenter())
		}
	}
}

func (m *Manager) MovePlayer(dir Direction) {
	playerObj := m.dynamicObjects["current-player"]
	curSpeed := playerObj.GetBaseSpeed()
	playerObj.AddForce(dirToVec2(dir).MulScalar(curSpeed))
}

func (m *Manager) ResetGame() {
	m.clearGameData()
	m.InitGame()
}

func (m *Manager) clearGameData() {
	m.dynamicObjects = make(map[string]DynamicObject)
	m.staticObjects = make(map[string]StaticObject)
}

func (m *Manager) Fart(dt float64) {
	fart := &fart{
		DynamicObject: &GObject{
			CollisionType: Circle,
			ParentObject:  m.dynamicObjects["current-player"],
			Center:        m.dynamicObjects["current-player"].GetCenter(),
			Width:         0.5,
			Height:        0.5,
			IsPassthrough: true,
			TimeToLive:    0.2,
		},
	}
	m.dynamicObjects["current-player"].AddChild(fart)
	m.pushAwayObjects(m.dynamicObjects["current-player"], 0.3, dt)
}

func (m *Manager) pushAwayObjects(pusherObject DynamicObject, dist float64, dt float64) {
	for _, obj := range m.dynamicObjects {
		if pusherObject == obj {
			continue
		}
		if obj.GetCenter().Distance(pusherObject.GetCenter()) > dist {
			continue
		}
		pushVector := obj.GetCenter().Sub(pusherObject.GetCenter()).Normalize().DivScalar(dt)
		obj.AddForce(pushVector)
	}

}
