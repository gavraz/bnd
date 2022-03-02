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
