package pixel

import (
	"fmt"
	"github.com/faiface/pixel"
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

func (h *Handler) HandleInput(menuHandler menuHandler) {
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
