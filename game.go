package main

import (
	"fmt"
	// "os"
	"time"
	engine "./goroguego"
)

// var logfile, _ = os.Create("gamelog.txt")

const (
	WORLD_WIDTH = 50
	WORLD_HEIGHT = 28
)

// -----------------------------------------------------------------------------------

type Object interface {
	PC()									bool
	HP()									int
	X()										int
	Y()										int
	String()								string
	Draw(win *engine.Window)
	Highlight(win *engine.Window)
}

type Moveable interface {
	TryMove(x, y int, world *World)			bool
}

// -----------------------------------------------------------------------------------

type Body struct {
	world		*World
	char		byte
	colour		byte
	x			int
	y			int
	hp			int
	pc			bool
}

func (u *Body) PC() bool {
	return u.pc
}

func (u *Body) HP() int {
	return u.hp
}

func (u *Body) X() int {
	return u.x
}

func (u *Body) Y() int {
	return u.y
}

func (u *Body) String() string {
	return "something"
}

func (u *Body) Draw(win *engine.Window) {
	win.Set(u.x, u.y, u.char, u.colour)
}

func (u *Body) Highlight(win *engine.Window) {
	win.SetHighlight(u.x, u.y)
}

// -----------------------------------------------------------------------------------

type Unit struct {
	Body
	class		string
	weapon		string
}

func (u *Unit) String() string {
	return fmt.Sprintf("- %s (%c) - %dhp - %s", u.class, u.char, u.hp, u.weapon)
}

func (u *Unit) TryMove(x, y int, world *World) bool {

	tar_x := u.x + x
	tar_y := u.y + y

	if world.InBounds(tar_x, tar_y) {
		u.x = tar_x
		u.y = tar_y
		return true
	}

	return false
}

// -----------------------------------------------------------------------------------

type World struct {
	window		*engine.Window
	width		int
	height		int
	selection	Object
	objects		[]Object
}

func (w *World) InBounds(x, y int) bool {
	if x >= 0 && x < w.width && y >= 0 && y < w.height {
		return true
	} else {
		return false
	}
}

func (w *World) Draw() {

	w.window.Clear()

	for _, object := range w.objects {
		object.Draw(w.window)
		if object == w.selection {
			object.Highlight(w.window)
		}
	}

	if (w.selection != nil) {
		s := w.selection.String()
		w.WriteSelection(s)
	}

	w.window.Flip()
}

func (w *World) AddObject(object Object) {
	w.objects = append(w.objects, object)
}

func (w *World) Start() {

	soldier := Unit{
		Body: Body{
			world: w,
			char: '@',
			colour: 'g',
			x: 5,
			y: 5,
			hp: 4,
			pc: true,
		},
		class: "soldier",
		weapon: "rifle",
	}

	w.AddObject(&soldier)
	w.Play()
}

func (w *World) Play() {

	for {

		// Deal with mouse events...

		for {
			click, err := engine.GetMousedown()
			if err != nil {
				break
			}
			w.selection = nil
			for _, object := range w.objects {
				if object.X() == click.X && object.Y() == click.Y {
					w.selection = object
				}
			}
		}

		// Deal with key events...

		var key = ""

		// For now, we just skip all but the last keypress on the queue...

		for {
			nextkey, err := engine.GetKeypress()
			if err != nil {
				break
			}
			key = nextkey
		}

		if key == "Escape" {
			w.selection = nil
		}

		if w.selection != nil && w.selection.PC() == true {

			mover, ok := w.selection.(Moveable)

			if ok {
				if key == "w" { mover.TryMove( 0, -1, w) }
				if key == "a" { mover.TryMove(-1,  0, w) }
				if key == "s" { mover.TryMove( 0,  1, w) }
				if key == "d" { mover.TryMove( 1,  0, w) }
			}
		}

		w.Draw()

		time.Sleep(50 * time.Millisecond)
	}
}

func (w *World) WriteSelection(s string) {
	for x := 0; x < len(s); x++ {
		w.window.Set(x, w.height + 1, s[x], 'w')
	}
}

func main() {

	world := World{
		window: engine.NewWindow("World", "renderer.html", WORLD_WIDTH, WORLD_HEIGHT + 2, 15, 20, 100, true),
		width: WORLD_WIDTH,
		height: WORLD_HEIGHT,
	}

	world.Start()
}
