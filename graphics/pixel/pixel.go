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

type Mover interface {
	MovePlayer(dirX float64, dirY float64, dt float64)
}

func (h *Handler) HandleInput(m Mover, dt float64) {
	var x, y float64
	if h.win.Pressed(pixelgl.KeyS) || h.win.Pressed(pixelgl.KeyDown) {
		y = -1.0
	}
	if h.win.Pressed(pixelgl.KeyW) || h.win.Pressed(pixelgl.KeyUp) {
		y = 1.0
	}
	if h.win.Pressed(pixelgl.KeyA) || h.win.Pressed(pixelgl.KeyLeft) {
		x = -1.0
	}
	if h.win.Pressed(pixelgl.KeyD) || h.win.Pressed(pixelgl.KeyRight) {
		x = 1.0
	}
	m.MovePlayer(x, y, dt)
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

	player := (objects["current-player"]).(*game.Player)
	h.drawGameBottomPanel(player.HP())
}

func (h *Handler) drawGameBottomPanel(hp int) {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(h.w()*0.02, h.h()*0.1), basicAtlas)
	_, _ = fmt.Fprintln(basicTxt, "HP: ", hp)
	basicTxt.Draw(h.win, pixel.IM.Scaled(basicTxt.Orig, 2.5))
}

func (h *Handler) drawGameObject(object game.Object) {
	switch object.(type) {
	case *game.Player:
		imd := imdraw.New(nil)
		imd.Color = colornames.Orange
		imd.Push(pixel.V(float64(object.GetCenter().X), float64(object.GetCenter().Y)))
		imd.Circle(object.GetWidth()/2, 0)
		imd.Draw(h.win)
	case *game.Crate:
		imd := imdraw.New(nil)
		imd.Color = colornames.Cyan
		imd.Push(pixel.V(float64(object.GetCenter().X-object.GetWidth()/2), float64(object.GetCenter().Y-object.GetHeight()/2)))
		imd.Push(pixel.V(float64(object.GetCenter().X+object.GetWidth()/2), float64(object.GetCenter().Y+object.GetHeight()/2)))
		imd.Rectangle(0)
		imd.Draw(h.win)
	case *game.BouncingBall:
		imd := imdraw.New(nil)
		imd.Color = colornames.Green
		imd.Push(pixel.V(float64(object.GetCenter().X), float64(object.GetCenter().Y)))
		imd.Circle(object.GetWidth()/2, 0)
		imd.Draw(h.win)

	default:
		fmt.Println("Unknown object type")
	}
}
