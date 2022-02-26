package pixel

import (
	"bnd/game"
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type choicer interface {
	CurrentChoice() int
	Choices() []string
}

type menuHandler interface {
	NextChoice()
	PrevChoice()
	Choose()
	GoBack()
}

type Handler struct {
	cfg pixelgl.WindowConfig
	win *pixelgl.Window
}

func New(cfg pixelgl.WindowConfig) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) Init() {
	var err error
	h.win, err = pixelgl.NewWindow(h.cfg)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) Closed() bool {
	return h.win.Closed()
}

func (h *Handler) Update() {
	h.win.Update()
}

func (h *Handler) DrawMenu(c choicer) {
	h.win.Clear(colornames.Skyblue)

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(100, 500), basicAtlas)

	current := c.CurrentChoice()
	for i, item := range c.Choices() {
		if i == current {
			basicTxt.Color = colornames.Red
		} else {
			basicTxt.Color = colornames.White
		}
		_, _ = fmt.Fprintln(basicTxt, item)
		basicTxt.Draw(h.win, pixel.IM.Scaled(basicTxt.Orig, 4))
	}
}

func (h *Handler) HandleMenuInput(menuHandler menuHandler) {
	if h.win.JustPressed(pixelgl.KeyS) || h.win.JustPressed(pixelgl.KeyDown) {
		menuHandler.NextChoice()
	}
	if h.win.JustPressed(pixelgl.KeyW) || h.win.JustPressed(pixelgl.KeyUp) {
		menuHandler.PrevChoice()
	}
	if h.win.JustPressed(pixelgl.KeyEnter) {
		menuHandler.Choose()
	}
	if h.win.JustPressed(pixelgl.KeyEscape) || h.win.JustPressed(pixelgl.KeyBackspace) {
		menuHandler.GoBack()
	}
}

type SetDirectioner interface {
	SetDirection(name string, direction game.Direction)
}

func (h *Handler) HandleInput(directioner SetDirectioner) {
	if h.win.Pressed(pixelgl.KeyS) || h.win.Pressed(pixelgl.KeyDown) {
		directioner.SetDirection("current-player", game.Direction{DX: 0, DY: -1})
	}
	if h.win.Pressed(pixelgl.KeyW) || h.win.Pressed(pixelgl.KeyUp) {
		directioner.SetDirection("current-player", game.Direction{DX: 0, DY: 1})
	}
	if h.win.Pressed(pixelgl.KeyA) || h.win.Pressed(pixelgl.KeyLeft) {
		directioner.SetDirection("current-player", game.Direction{DX: -1, DY: 0})
	}
	if h.win.Pressed(pixelgl.KeyD) || h.win.Pressed(pixelgl.KeyRight) {
		directioner.SetDirection("current-player", game.Direction{DX: 1, DY: 0})
	}
}

func (h *Handler) DrawGame(objects map[string]game.Object) {
	h.win.Clear(colornames.Black) // TODO decide color

	var width = h.cfg.Bounds.W() * 1.0
	var height = h.cfg.Bounds.H() * 1.0
	var sidePadding = h.cfg.Bounds.W() * 0.02
	var bottomPadding = h.cfg.Bounds.H() * 0.15
	border := imdraw.New(nil)
	border.Color = pixel.RGB(255, 255, 255)
	border.Push(pixel.V(sidePadding, bottomPadding+sidePadding))
	border.Push(pixel.V(width-sidePadding, height-sidePadding))
	border.Rectangle(1)
	border.Draw(h.win)

	for _, o := range objects {
		h.drawGameObject(o)
	}
}

func (h *Handler) drawGameObject(object game.Object) {
	switch object.(type) {
	case *game.Player:
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(1, 0, 0)
		imd.Push(pixel.V(float64(object.GetPoint().X), float64(object.GetPoint().Y)))
		imd.Color = pixel.RGB(0, 1, 0)
		imd.Push(pixel.V(float64(object.GetPoint().X+50), float64(object.GetPoint().Y)))
		imd.Color = pixel.RGB(0, 0, 1)
		imd.Push(pixel.V(float64(object.GetPoint().X+50), float64(object.GetPoint().Y+50)))
		imd.Color = pixel.RGB(0, 1, 1)
		imd.Push(pixel.V(float64(object.GetPoint().X), float64(object.GetPoint().Y+50)))

		imd.Polygon(0)

		imd.Draw(h.win)
	case *game.Crate:
	default:
		fmt.Println("Unknown object type")
	}
}
