package pixel

import (
	"bnd/game"
	"bnd/input"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type choicer interface {
	CurrentChoice() int
	Choices() []string
}

type Environmenter interface {
	ForEachGameObject(do func(object game.Object))
	HP() int
}

type Handler struct {
	win *pixelgl.Window
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Init(cfg pixelgl.WindowConfig) {
	var err error
	h.win, err = pixelgl.NewWindow(cfg)
	h.win.SetSmooth(true)

	loadAssets()

	if err != nil {
		panic(err)
	}
}

func (h *Handler) ChangeResolution(width, height int) {
	h.win.SetBounds(pixel.R(0, 0, float64(width), float64(height)))
}

func (h *Handler) Closed() bool {
	return h.win.Closed()
}

func (h *Handler) Update() {

	h.win.Update()
}

func (h *Handler) w() float64 {
	return h.win.Bounds().W()
}

func (h *Handler) h() float64 {
	return h.win.Bounds().H()
}

func (h *Handler) DrawMainMenu(c choicer) {
	h.win.Clear(color.RGBA{R: 10, G: 30, B: 30, A: 255})
	backgroundSprite.Draw(h.win, pixel.IM.Moved(h.win.Bounds().Center()))
	h.drawMenuText(c, fontFace, colornames.Red, colornames.White)
}

func (h *Handler) DrawPauseMenu(c choicer) {
	imd := imdraw.New(nil)
	imd.Color = color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 150,
	}
	imd.Push(h.win.Bounds().Min)
	imd.Push(h.win.Bounds().Max)
	imd.Rectangle(0)
	imd.Draw(h.win)

	h.drawMenuText(c, fontFace, colornames.Red, colornames.White)
}

func (h *Handler) drawMenuText(c choicer, fontface font.Face, highlighted color.Color, normal color.Color) {
	atlas := text.NewAtlas(fontface, text.ASCII)
	var v game.Vector2
	v = h.toGlobalSpace(v)
	txt := text.New(pixel.V(v.X, v.Y), atlas)
	current := c.CurrentChoice()
	for i, item := range c.Choices() {
		if i == current {
			txt.Color = highlighted
		} else {
			txt.Color = normal
		}
		_, _ = fmt.Fprintln(txt, strings.ToLower(item)) // this 'Mario-Kart-DS' font is different for upper and lower
	}
	txt.Draw(h.win, pixel.IM.Moved(pixel.V(-txt.Bounds().W()/2, txt.Bounds().H()/2)).Scaled(txt.Orig, math.Max(1-txt.Bounds().W()/h.win.Bounds().W(), 0.5)))
}

func (h *Handler) Pressed(key input.Key) bool {
	return h.win.Pressed(toPixel(key))
}

func (h *Handler) JustPressed(key input.Key) bool {
	return h.win.JustPressed(toPixel(key))
}

func (h *Handler) DrawGame(env Environmenter) {
	h.win.Clear(colornames.Black) // TODO decide color

	env.ForEachGameObject(h.drawGameObject)
	h.drawGameBottomPanel(env.HP())
}

func (h *Handler) drawGameBottomPanel(hp int) {

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(h.w()*0.02, h.h()*0.1), basicAtlas)
	_, _ = fmt.Fprintln(basicTxt, "HP: ", hp)
	basicTxt.Draw(h.win, pixel.IM.Scaled(basicTxt.Orig, 2.5))
}

func (h *Handler) aspectRatio() float64 {
	return h.w() / h.h()
}

func (h *Handler) toLocalSpace(v game.Vector2) game.Vector2 {
	aspectRatio := h.aspectRatio()
	return game.Vector2{X: v.X/h.w()*2 - 1*aspectRatio, Y: v.Y/h.h()*2 - 1}
}

// Converts a point from local space to global space (i.e. screen space)
func (h *Handler) toGlobalSpace(v game.Vector2) game.Vector2 {
	aspectRatio := h.aspectRatio()
	return game.Vector2{X: (v.X + aspectRatio) * h.w() / 2 / aspectRatio, Y: (v.Y + 1) * h.h() / 2}
}

func (h *Handler) toGlobalUnits(v game.Vector2) game.Vector2 {
	aspectRatio := h.aspectRatio()
	return game.Vector2{X: v.X * h.w() / 2 / aspectRatio, Y: v.Y * h.h() / 2}
}

func (h *Handler) drawGameObject(obj game.Object) {
	switch game.ObjectType(obj) {
	case game.Player:
		playerCenter := h.toGlobalSpace(obj.GetCenter())
		playerSize := h.toGlobalUnits(game.Vector2{X: obj.GetWidth(), Y: obj.GetHeight()})
		imd := imdraw.New(nil)
		imd.Color = colornames.Orange
		imd.Push(pixel.V(playerCenter.X, playerCenter.Y))
		imd.Ellipse(pixel.Vec{X: playerSize.X / 2, Y: playerSize.Y / 2}, 0)

		imd.Draw(h.win)
	case game.Crate:
		crateCenter := h.toGlobalSpace(obj.GetCenter())
		crateSize := h.toGlobalUnits(game.Vector2{X: obj.GetWidth(), Y: obj.GetHeight()})
		imd := imdraw.New(nil)
		imd.Color = colornames.Cyan
		imd.Push(pixel.V(crateCenter.X-crateSize.X/2, crateCenter.Y-crateSize.Y/2))
		imd.Push(pixel.V(crateCenter.X+crateSize.X/2, crateCenter.Y+crateSize.Y/2))
		imd.Rectangle(0)
		imd.Draw(h.win)
	case game.Wall:
		wallCenter := h.toGlobalSpace(obj.GetCenter())
		wallSize := h.toGlobalUnits(game.Vector2{X: obj.GetWidth(), Y: obj.GetHeight()})
		imd := imdraw.New(nil)
		imd.Color = colornames.Darkblue
		imd.Push(pixel.V(wallCenter.X-wallSize.X/2, wallCenter.Y-wallSize.Y/2))
		imd.Push(pixel.V(wallCenter.X+wallSize.X/2, wallCenter.Y+wallSize.Y/2))
		imd.Rectangle(0)
		imd.Draw(h.win)
	case game.Fart:
		fartCenter := h.toGlobalSpace(obj.GetCenter())
		fartSize := h.toGlobalUnits(game.Vector2{X: obj.GetWidth(), Y: obj.GetHeight()})
		imd := imdraw.New(nil)
		imd.Color = color.RGBA{
			R: 0,
			G: 100,
			B: 0,
			A: 100,
		}
		imd.Push(pixel.V(fartCenter.X, fartCenter.Y))
		imd.Ellipse(pixel.Vec{X: fartSize.X / 2, Y: fartSize.Y / 2}, 0)
		imd.Draw(h.win)
	default:
		fmt.Println("pixel: drawing unimplemented for type")
	}
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
