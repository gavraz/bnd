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

type Objecter interface {
	Objects() map[string]game.Object
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

func (h *Handler) w() float64 {
	return h.cfg.Bounds.W()
}

func (h *Handler) h() float64 {
	return h.cfg.Bounds.H()
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

func (h *Handler) HandleInput(env Objecter) {
	playerObj := env.Objects()["current-player"]
	if h.win.Pressed(pixelgl.KeyS) || h.win.Pressed(pixelgl.KeyDown) {
		playerObj.SetAcceleration(game.Vector2{X: 0, Y: -0.3})
	} else if h.win.Pressed(pixelgl.KeyW) || h.win.Pressed(pixelgl.KeyUp) {
		playerObj.SetAcceleration(game.Vector2{X: 0, Y: 0.3})
	} else if h.win.Pressed(pixelgl.KeyA) || h.win.Pressed(pixelgl.KeyLeft) {
		playerObj.SetAcceleration(game.Vector2{X: -0.3, Y: 0})
	} else if h.win.Pressed(pixelgl.KeyD) || h.win.Pressed(pixelgl.KeyRight) {
		playerObj.SetAcceleration(game.Vector2{X: 0.3, Y: 0})
	} else {
		playerObj.SetAcceleration(game.Vector2{X: 0, Y: 0})
	}
}

func (h *Handler) DrawGame(env Objecter) {
	objects := env.Objects()
	h.win.Clear(colornames.Black) // TODO decide color

	var sidePadding = h.cfg.Bounds.W() * 0.02
	var bottomPadding = h.cfg.Bounds.H() * 0.15
	border := imdraw.New(nil)
	border.Color = pixel.RGB(255, 255, 255)
	border.Push(pixel.V(sidePadding, bottomPadding+sidePadding))
	border.Push(pixel.V(h.w()-sidePadding, h.h()-sidePadding))
	border.Rectangle(1)
	border.Draw(h.win)

	for _, o := range objects {
		h.drawGameObject(o)
	}

	h.drawGameBottomPanel()
}

func (h *Handler) drawGameBottomPanel() {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(h.w()*0.02, h.h()*0.1), basicAtlas)
	_, _ = fmt.Fprintln(basicTxt, "HP: 100")
	basicTxt.Draw(h.win, pixel.IM.Scaled(basicTxt.Orig, 2.5))
}

func (h *Handler) drawGameObject(object game.Object) {
	switch object.(type) {
	case *game.Player:
		imd := imdraw.New(nil)
		imd.Color = colornames.Orange
		imd.Push(pixel.V(object.GetPoint().X, object.GetPoint().Y))
		imd.Circle(h.w()*0.01, 0)
		imd.Draw(h.win)
	case *game.Crate:
	default:
		fmt.Println("Unknown object type")
	}
}
