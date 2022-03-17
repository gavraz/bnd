package game

type player struct {
	DynamicObject

	hp        int
	primary   Ability
	secondary Ability
}

func NewPlayer(object DynamicObject, hp int) *player {
	return &player{
		DynamicObject: object,
		hp:            hp,
	}
}

func (p *player) SwapAbilities() {
	p.primary, p.secondary = p.secondary, p.primary
}

type crate struct {
	DynamicObject

	ability Ability
}

type wall struct {
	StaticObject
}

type fart struct {
	DynamicObject
}
