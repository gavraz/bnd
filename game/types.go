package game

type objectType int

const (
	TypePlayer objectType = iota
	TypeCrate
)

func ObjectType(object Object) objectType {
	switch object.(type) {
	case *Player:
		return TypePlayer
	case *Crate:
		return TypeCrate
	default:
		panic("unknown object type")
	}
	return 0
}
