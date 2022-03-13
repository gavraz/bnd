package pixel

import (
	"bnd/input"
	"github.com/faiface/pixel/pixelgl"
)

func toPixel(key input.Key) pixelgl.Button {
	switch key {
	case input.KeyW:
		return pixelgl.KeyW
	case input.KeyD:
		return pixelgl.KeyD
	case input.KeyS:
		return pixelgl.KeyS
	case input.KeyA:
		return pixelgl.KeyA
	case input.KeyEsc:
		return pixelgl.KeyEscape
	case input.KeyEnter:
		return pixelgl.KeyEnter
	case input.KeySpace:
		return pixelgl.KeySpace
	default:
		return pixelgl.KeyUnknown
	}
}
