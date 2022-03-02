package game

type objectType int

const (
	TypePlayer objectType = iota
	TypeCrate
	TypeWall
)

func ObjectType(object Object) objectType {
	switch object.(type) {
	case *Player:
		return TypePlayer
	case *Crate:
		return TypeCrate
	case *Wall:
		return TypeWall
	default:
		panic("unknown object type")
	}
}
