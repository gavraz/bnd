package game

import (
	"bnd/engine"
	"math"
	"time"
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

func (m *Manager) InitGame() {
	m.env.AddDynamicObject("current-player", NewPlayer(&engine.GObject{
		Center: engine.Vector2{
			X: 0,
			Y: 0,
		},
		BaseSpeed:     2,
		CollisionType: engine.Circle,
		Width:         0.05,
		Height:        0.05,
		Mass:          1,
		Direction:     engine.Vector2{Y: 1},
		Friction:      4.0,
	}, 100))
	m.env.AddStaticObject("enemy-player", NewPlayer(&engine.GObject{
		Center: engine.Vector2{
			X: 0.2,
			Y: 0.3,
		},
		BaseSpeed:     1,
		CollisionType: engine.Circle,
		Width:         0.2,
		Height:        0.2,
		Friction:      4.0,
	}, 100))
	m.env.AddDynamicObject("crate", &crate{
		DynamicObject: &engine.GObject{
			Center: engine.Vector2{
				X: -0.2,
				Y: -0.2,
			},
			CollisionType: engine.Rectangle,
			Width:         0.1,
			Height:        0.1,
			Mass:          1,
			Friction:      4.0,
		},
	})
	m.env.AddDynamicObject("crate2", &crate{
		DynamicObject: &engine.GObject{
			Center: engine.Vector2{
				X: -0.6,
				Y: -0.5,
			},
			CollisionType: engine.Rectangle,
			Width:         0.2,
			Height:        0.2,
			Mass:          2,
			Friction:      4.0,
		},
	})
	m.env.AddStaticObject("wall-bottom", &wall{
		StaticObject: &engine.GObject{
			Center: engine.Vector2{
				X: 0,
				Y: -0.83,
			},
			CollisionType: engine.Rectangle,
			Width:         1.92,
			Height:        0.34,
		},
	})
	m.env.AddStaticObject("wall-left", &wall{
		StaticObject: &engine.GObject{
			Center: engine.Vector2{
				X: -0.98,
				Y: 0,
			},
			CollisionType: engine.Rectangle,
			Width:         0.04,
			Height:        2,
		},
	})
	m.env.AddStaticObject("wall-right", &wall{
		StaticObject: &engine.GObject{
			Center: engine.Vector2{
				X: 0.98,
				Y: 0,
			},
			CollisionType: engine.Rectangle,
			Width:         0.04,
			Height:        2,
		},
	})
	m.env.AddStaticObject("wall-top", &wall{
		StaticObject: &engine.GObject{
			Center: engine.Vector2{
				X: 0,
				Y: 0.98,
			},
			CollisionType: engine.Rectangle,
			Width:         1.92,
			Height:        0.04,
		},
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

func (m *Manager) Fart(dt float64) {
	player := m.env.ObjectByName("current-player").(*player)
	fart := &fart{
		DynamicObject: &engine.GObject{
			CollisionType: engine.Circle,
			ParentObject:  player,
			Center:        player.GetCenter(),
			Width:         0.5,
			Height:        0.5,
			IsPassthrough: true,
			Until:         time.Now().Add(100 * time.Millisecond),
		},
	}
	player.AddChild(fart)
	m.pushAwayObjects(fart, dt)
}

func (m *Manager) Melee() {
	lifeTime := 0.15
	radius := 0.1
	size := 0.01
	user := m.env.ObjectByName("current-player").(*player)
	sword := newMeleeObject(&engine.GObject{
		CollisionType: engine.Circle,
		Width:         size,
		Height:        size,
		Mass:          1,
		ParentObject:  user,
		Until:         time.Now().Add(time.Duration(lifeTime*1000) * time.Millisecond),
		IsPassthrough: true,
	}, user.GetDirection(), user.GetCenter(), user.GetWidth(), math.Pi/4, lifeTime, size, radius)
	user.AddChild(sword)
}

func (m *Manager) pushAwayObjects(pusherObject engine.DynamicObject, dt float64) {
	if collider := m.env.ResolveDynamicCollisions(pusherObject); collider != nil {
		pushVector := collider.GetCenter().Sub(pusherObject.GetCenter()).Normalize().DivScalar(dt)
		collider.AddForce(pushVector)
	}
}

func (m *Manager) Update(dt float64) {
	m.env.Update(dt)
}

func (m *Manager) ForEachGameObject(do func(object engine.Object)) {
	m.env.ForEachGameObject(do)
}
