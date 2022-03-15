package pixel

import (
	"github.com/faiface/pixel"
	"golang.org/x/image/font"
)

var fontFace font.Face
var backgroundImage pixel.Picture
var backgroundSprite *pixel.Sprite

func loadAssets() {
	fontFace, _ = loadTTF("graphics/pixel/assets/fonts/Mario-Kart-DS.ttf", 72)
	backgroundImage, _ = loadPicture("graphics/pixel/assets/images/img.png")
	backgroundSprite = pixel.NewSprite(backgroundImage, backgroundImage.Bounds())
}
