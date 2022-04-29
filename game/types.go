package game

import (
	"bnd/engine"
	"fmt"
)

type ObjType int

const (
	Player ObjType = iota
	Crate
	Wall
	Melee
	Fart
)

func ObjectType(object engine.Object) ObjType {
	switch object.(type) {
	case *player:
		return Player
	case *crate:
		return Crate
	case *wall:
		return Wall
	case *meleeObject, *meleeParticle:
		return Melee
	case *fartObject:
		return Fart
	default:
		fmt.Printf("could not infer game object type, received: %T\n", object)
		panic("unknown object type ")
	}
}
