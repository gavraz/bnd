package game

type ObjType int

const (
	Player ObjType = iota
	Crate
	Wall
)

func ObjectType(object Object) ObjType {
	switch object.(type) {
	case *player:
		return Player
	case *crate:
		return Crate
	case *wall:
		return Wall
	default:
		panic("unknown object type")
	}
}
