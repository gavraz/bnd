package game

type Player struct {
	DynamicObject

	hp        int
	primary   Ability
	secondary Ability
}

func NewPlayer(object DynamicObject, hp int) *Player {
	return &Player{
		DynamicObject: object,
		hp:            hp,
	}
}

func (p *Player) SwapAbilities() {
	p.primary, p.secondary = p.secondary, p.primary
}
