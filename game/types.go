package game

import "fmt"

type ObjType int

const (
	Player ObjType = iota
	Crate
	Wall
	Melee
	Fart
)

func ObjectType(object Object) ObjType {
	switch object.(type) {
	case *player:
		return Player
	case *crate:
		return Crate
	case *wall:
		return Wall
	case *meleeObject:
		return Melee
	case *fart:
		return Fart
	default:
		fmt.Printf("could not infer game object type, received: %T", object)
		panic("unknown object type ")
	}
}
