package game

type player struct {
	Object

	hp        int
	primary   Ability
	secondary Ability
}

func NewPlayer(object Object, hp int) *player {
	return &player{
		Object: object,
		hp:     hp,
	}
}

func (p *player) SwapAbilities() {
	p.primary, p.secondary = p.secondary, p.primary
}

type crate struct {
	Object

	ability Ability
}

type wall struct {
	Object
}

type Bullet struct {
	Object
}
