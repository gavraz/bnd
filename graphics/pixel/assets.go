package pixel

import (
	"github.com/faiface/pixel"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"io/ioutil"
	"os"
)

var (
	fontFace         font.Face
	backgroundImage  pixel.Picture
	backgroundSprite *pixel.Sprite
)

func loadAssets() {
	fontFace, _ = loadTTF("graphics/pixel/assets/fonts/Mario-Kart-DS.ttf", 72)
	backgroundImage, _ = loadPicture("graphics/pixel/assets/images/img.png")
	backgroundSprite = pixel.NewSprite(backgroundImage, backgroundImage.Bounds())
}

// Loads a picture into a usable pixel.Picture format
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {}(file)
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

// Loads a .ttf font file into a usable font.Face format
func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {}(file)

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	fnt, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(fnt, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

// Loads a .otf font file into a usable font.Face format
func loadOTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {}(file)

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	fnt, err := opentype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: 0,
	})
}
