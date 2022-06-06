package game

import (
	"bnd/engine"
	"math/rand"
)

type Manager struct {
	env *engine.Environment
}

func NewManager() *Manager {
	return &Manager{
		env: engine.NewEnvironment(),
	}
}

func (m *Manager) HP() int {
	return m.env.ObjectByName("current-player").(*player).hp
}

func (m *Manager) removeDynamicObject(object engine.DynamicObject) {
	m.env.RemoveDynamicObject(object)
}

func (m *Manager) removeStaticObject(object engine.StaticObject) {
	m.env.RemoveStaticObject(object)
}

func (m *Manager) InitGame() {
	m.env.AddDynamicObject(&player{
		DynamicObject: engine.NewDynamicObject(engine.GameObjectConf{
			Name: "current-player",
			Center: engine.Vector2{
				X: 0,
				Y: 0,
			},
			BaseSpeed:     3,
			CollisionType: engine.Circle,
			Width:         0.05,
			Height:        0.05,
			Mass:          1,
			Direction:     engine.Vector2{Y: 1},
			Friction:      4.0,
		}),
		hp: 100,
	})

	m.env.AddDynamicObject(&player{
		DynamicObject: engine.NewDynamicObject(engine.GameObjectConf{
			Name: "enemy-player",
			Center: engine.Vector2{
				X: 0.2,
				Y: 0.3,
			},
			BaseSpeed:     1,
			CollisionType: engine.Circle,
			Width:         0.2,
			Height:        0.2,
			Mass:          1,
			Friction:      4.0,
		}),
		hp: 100,
	})
	for i := 0; i < 10; i++ {
		crateName := "Crate" + string(i)
		pos := engine.Vector2{rand.Float64()*1.8 - 0.9, rand.Float64()*1.5 - 0.6}
		crateSize := 0.05
		m.env.AddStaticObject(&crate{
			StaticObject: engine.NewStaticObject(engine.GameObjectConf{
				Name:          crateName,
				Center:        pos,
				Width:         crateSize,
				Height:        crateSize,
				IsPassthrough: true,
			}),
			removeCrate: m.removeStaticObject})
	}

	m.env.AddStaticObject(&wall{
		StaticObject: engine.NewStaticObject(engine.GameObjectConf{
			Name: "wall-bottom",
			Center: engine.Vector2{
				X: 0,
				Y: -0.83,
			},
			CollisionType: engine.Rectangle,
			Width:         1.92,
			Height:        0.34,
		}),
	})
	m.env.AddStaticObject(&wall{
		StaticObject: engine.NewStaticObject(engine.GameObjectConf{
			Name: "wall-left",
			Center: engine.Vector2{
				X: -0.98,
				Y: 0,
			},
			CollisionType: engine.Rectangle,
			Width:         0.04,
			Height:        2,
		}),
	})
	m.env.AddStaticObject(&wall{
		StaticObject: engine.NewStaticObject(engine.GameObjectConf{
			Name: "wall-right",
			Center: engine.Vector2{
				X: 0.98,
				Y: 0,
			},
			CollisionType: engine.Rectangle,
			Width:         0.04,
			Height:        2,
		}),
	})
	m.env.AddStaticObject(&wall{
		StaticObject: engine.NewStaticObject(engine.GameObjectConf{
			Name: "wall-top",
			Center: engine.Vector2{
				X: 0,
				Y: 0.98,
			},
			CollisionType: engine.Rectangle,
			Width:         1.92,
			Height:        0.04,
		}),
	})

}

func (m *Manager) MovePlayer(dir Direction) {
	m.env.ObjectByName("current-player")
	playerObj := m.env.ObjectByName("current-player").(*player)
	curSpeed := playerObj.GetBaseSpeed()
	playerObj.AddForce(dir.v.MulScalar(curSpeed))
}

func (m *Manager) ResetGame() {
	m.clearGameData()
	m.InitGame()
}

func (m *Manager) clearGameData() {
	m.env.ClearGameData()
}

func (m *Manager) Fart() {
	player := m.env.ObjectByName("current-player").(*player)
	fart := newFartObject(player, 0.5, 0.5)
	player.AddChild(fart)
}

func (m *Manager) Melee() {
	player := m.env.ObjectByName("current-player").(*player)
	addMeleeObject(player)
}

func (m *Manager) Update(dt float64) {
	m.env.Update(dt)
}

func (m *Manager) ForEachGameObject(do func(object engine.Object)) {
	m.env.ForEachGameObject(do)
}
