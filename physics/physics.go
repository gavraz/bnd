package physics

import "bnd/game"

func Move(object game.Object) {
	object.SetPoint(game.Point{
		X: object.GetPoint().X + object.GetDirection().DX*object.GetVelocity(),
		Y: object.GetPoint().Y + object.GetDirection().DY*object.GetVelocity(),
	})
}
