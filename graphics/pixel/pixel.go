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
	win *pixelgl.Window
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Init(cfg pixelgl.WindowConfig) {
	var err error
	h.win, err = pixelgl.NewWindow(cfg)
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
	return h.win.Bounds().W()
}

func (h *Handler) h() float64 {
	return h.win.Bounds().H()
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

func (h *Handler) HandleInput(env Objecter, dt float64) {
	playerObj := env.Objects()["current-player"]
	vec := game.Vector2{X: 0, Y: 0}
	if h.win.Pressed(pixelgl.KeyS) || h.win.Pressed(pixelgl.KeyDown) {
		vec.Y += (-playerObj.GetBaseSpeed() * dt)
	}
	if h.win.Pressed(pixelgl.KeyW) || h.win.Pressed(pixelgl.KeyUp) {
		vec.Y += (playerObj.GetBaseSpeed() * dt)
	}
	if h.win.Pressed(pixelgl.KeyA) || h.win.Pressed(pixelgl.KeyLeft) {
		vec.X += (-playerObj.GetBaseSpeed() * dt)
	}
	if h.win.Pressed(pixelgl.KeyD) || h.win.Pressed(pixelgl.KeyRight) {
		vec.X += (playerObj.GetBaseSpeed() * dt)
	}
	playerObj.SetAcceleration(vec)
}

func (h *Handler) DrawGame(env Objecter) {
	objects := env.Objects()
	h.win.Clear(colornames.Black) // TODO decide color

	var sidePadding = h.w() * 0.02
	var bottomPadding = h.h() * 0.15
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

func (h *Handler) toLocalSpace(v game.Vector2) game.Vector2 {
	return game.Vector2{X: v.X/h.w()*2 - 1, Y: v.Y/h.h()*2 - 1}
}

func (h *Handler) toGlobalSpace(v game.Vector2) game.Vector2 {
	return game.Vector2{X: (v.X + 1) * h.w() / 2, Y: (v.Y + 1) * h.h() / 2}
}

func (h *Handler) toGlobalUnits(v game.Vector2) game.Vector2 {
	return game.Vector2{X: v.X * h.w() / 2, Y: v.Y * h.h() / 2}
}

func (h *Handler) drawGameObject(object game.Object) {
	switch object.(type) {
	case *game.Player:
		playerCenter := h.toGlobalSpace(object.GetCenter())
		playerSize := h.toGlobalUnits(game.Vector2{X: object.GetWidth(), Y: object.GetHeight()})
		imd := imdraw.New(nil)
		imd.Color = colornames.Orange
		imd.Push(pixel.V(playerCenter.X, playerCenter.Y))
		imd.Ellipse(pixel.Vec{X: playerSize.X / 2, Y: playerSize.Y / 2}, 0)
		imd.Draw(h.win)
	case *game.Crate:
		crateCenter := h.toGlobalSpace(object.GetCenter())
		crateSize := h.toGlobalUnits(game.Vector2{X: object.GetWidth(), Y: object.GetHeight()})
		imd := imdraw.New(nil)
		imd.Color = colornames.Cyan
		imd.Push(pixel.V(crateCenter.X-crateSize.X/2, crateCenter.Y-crateSize.Y/2))
		imd.Push(pixel.V(crateCenter.X+crateSize.X/2, crateCenter.Y+crateSize.Y/2))
		imd.Rectangle(0)
		imd.Draw(h.win)
	case *game.BouncingBall:
		ballCenter := h.toGlobalSpace(object.GetCenter())
		ballSize := h.toGlobalUnits(game.Vector2{X: object.GetWidth(), Y: object.GetHeight()})
		imd := imdraw.New(nil)
		imd.Color = colornames.Green
		imd.Push(pixel.V(ballCenter.X, ballCenter.Y))
		imd.Ellipse(pixel.Vec{X: ballSize.X / 2, Y: ballSize.Y / 2}, 0)
		imd.Draw(h.win)
	case *game.Wall:
		wallCenter := h.toGlobalSpace(object.GetCenter())
		wallSize := h.toGlobalUnits(game.Vector2{X: object.GetWidth(), Y: object.GetHeight()})
		imd := imdraw.New(nil)
		imd.Color = colornames.Darkblue
		imd.Push(pixel.V(wallCenter.X-wallSize.X/2, wallCenter.Y-wallSize.Y/2))
		imd.Push(pixel.V(wallCenter.X+wallSize.X/2, wallCenter.Y+wallSize.Y/2))
		imd.Rectangle(0)
		imd.Draw(h.win)

	default:
		fmt.Println("Unknown object type")
	}

}
