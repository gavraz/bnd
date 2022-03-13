package game

type ObjType int

const (
	Player ObjType = iota
	Crate
	Wall
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
	case *fart:
		return Fart
	default:
		panic("unknown object type")
	}
}
