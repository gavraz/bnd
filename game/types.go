package game

import "fmt"

type ObjType int

const (
	Player ObjType = iota
	Crate
	Wall
	MeleeObject
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
		return MeleeObject
	case *fart:
		return Fart
	default:
		fmt.Printf("%T", object)
		panic("unknown object type ")
	}
}
