package game

type Player struct {
	Object

	hp        int
	primary   Ability
	secondary Ability
}

func NewPlayer(object Object, hp int) *Player {
	return &Player{
		Object: object,
		hp:     hp,
	}
}

func (p *Player) SwapAbilities() {
	p.primary, p.secondary = p.secondary, p.primary
}

func (p *Player) HP() int {
	return p.hp
}

func (m *Manager) MovePlayer(dirX float64, dirY float64, dt float64) {
	vec := Vector2{X: 0, Y: 0}
	playerObj := m.Objects()["current-player"]
	curSpeed := playerObj.GetBaseSpeed()
	vec.Y += dirY * curSpeed * dt
	vec.X += dirX * curSpeed * dt
	playerObj.SetAcceleration(vec)
}
