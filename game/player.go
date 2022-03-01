package game

import (
	"fmt"
	"time"
)

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

func (p *Player) CastAbility(ab Ability) {
	if ab == Speed {
		fmt.Println(p.GetBaseSpeed(), abilitySettings["Acceleration"][ab])
		p.SetBaseSpeed(p.GetBaseSpeed() * abilitySettings["Acceleration"][ab])
		fmt.Println(p.GetBaseSpeed())
		//Maybe change it later into a function which is called based on a timer
		go func() {
			time.Sleep(time.Duration(abilitySettings["Duration"][ab]) * time.Second)
			p.SetBaseSpeed(p.GetBaseSpeed() / abilitySettings["Acceleration"][ab])
			fmt.Println(p.GetBaseSpeed())
		}()
	}
}

func (p *Player) SwapAbilities() {
	p.primary, p.secondary = p.secondary, p.primary
}

func (p *Player) HP() int {
	return p.hp
}
