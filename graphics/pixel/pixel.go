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

type Environmenter interface {
	ForEachGameObject(do func(object game.Object))
	HP() int
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

func (h *Handler) DrawMenu(c choicer) {
	h.win.Clear(colornames.Skyblue)

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	v := game.Vector2{0.0, 0.0}
	v = h.toGlobalSpace(v)
	basicTxt := text.New(pixel.V(v.X, v.Y), basicAtlas)
	current := c.CurrentChoice()
	for i, item := range c.Choices() {
		if i == current {
			basicTxt.Color = colornames.Red
		} else {
			basicTxt.Color = colornames.White
		}
		_, _ = fmt.Fprintln(basicTxt, item)
	}
	basicTxt.Orig = basicTxt.Orig.Add(pixel.V(0.0, basicTxt.Bounds().H()/2))
	basicTxt.Draw(h.win, pixel.IM.Moved(pixel.V(-basicTxt.Bounds().W()/2, basicTxt.Bounds().H()/2)).Scaled(basicTxt.Orig, 2.0))
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
	MovePlayer(direction game.Direction)
}

func (h *Handler) HandleInput(m Mover) {
	var dir game.Direction
	if h.win.Pressed(pixelgl.KeyW) || h.win.Pressed(pixelgl.KeyUp) {
		dir.Up()
	}
	if h.win.Pressed(pixelgl.KeyS) || h.win.Pressed(pixelgl.KeyDown) {
		dir.Down()
	}
	if h.win.Pressed(pixelgl.KeyA) || h.win.Pressed(pixelgl.KeyLeft) {
		dir.Left()
	}
	if h.win.Pressed(pixelgl.KeyD) || h.win.Pressed(pixelgl.KeyRight) {
		dir.Right()
	}
	m.MovePlayer(dir)
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

func (h *Handler) toLocalSpace(v game.Vector2) game.Vector2 {
	var aspectRatio = h.w() / h.h()
	return game.Vector2{X: v.X/h.w()*2 - 1*aspectRatio, Y: v.Y/h.h()*2 - 1}
}

// Converts a point from local space to global space (i.e. screen space)
func (h *Handler) toGlobalSpace(v game.Vector2) game.Vector2 {
	var aspectRatio = h.w() / h.h()
	return game.Vector2{X: (v.X + aspectRatio) * h.w() / 2 / aspectRatio, Y: (v.Y + 1) * h.h() / 2}
}

func (h *Handler) toGlobalUnits(v game.Vector2) game.Vector2 {
	var aspectRatio = h.w() / h.h()
	return game.Vector2{X: v.X * h.w() / 2 / aspectRatio, Y: v.Y * h.h() / 2}
}

func (h *Handler) drawGameObject(obj game.Object) {
	switch game.ObjectType(obj) {
	case game.Player:
		playerCenter := h.toGlobalSpace(obj.GetCenter())
		playerSize := h.toGlobalUnits(game.Vector2{X: obj.GetWidth(), Y: obj.GetHeight()})
		imd := imdraw.New(nil)
		imd.Color = colornames.Orange

		switch obj.GetCollisionType() {
		case game.Circle:
			imd.Push(pixel.V(playerCenter.X, playerCenter.Y))
			imd.Ellipse(pixel.Vec{X: playerSize.X / 2, Y: playerSize.Y / 2}, 0)
		case game.Rectangle:
			imd.Push(pixel.V(playerCenter.X-playerSize.X/2, playerCenter.Y-playerSize.Y/2))
			imd.Push(pixel.V(playerCenter.X+playerSize.X/2, playerCenter.Y+playerSize.Y/2))
			imd.Rectangle(0)
		}

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
	default:
		fmt.Println("pixel: drawing unimplemented for type")
	}
}
