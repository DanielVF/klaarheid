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

type Unit struct {
	world		*World
	char		byte
	colour		byte
	class		string
	weapon		string
	x			int
	y			int
	hp			int
	pc			bool
}

func (u *Unit) String() string {
	return fmt.Sprintf("- %s (%dhp), %s", u.class, u.hp, u.weapon)
}

func (u *Unit) TryMove(x, y int) {

	tar_x := u.x + x
	tar_y := u.y + y

	if u.world.InBounds(tar_x, tar_y) {
		u.x = tar_x
		u.y = tar_y
	}
}

type World struct {
	window		*engine.Window
	width		int
	height		int
	selection	*Unit
	units		[]*Unit
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

	for _, unit := range w.units {
		w.window.Set(unit.x, unit.y, unit.char, unit.colour)
		if unit == w.selection {
			w.window.SetHighlight(unit.x, unit.y)
		}
	}

	if (w.selection != nil) {
		s := w.selection.String()
		w.WriteSelection(s)
	}

	w.window.Flip()
}

func (w *World) AddUnit(unit *Unit) {
	w.units = append(w.units, unit)
}

func (w *World) Start() {

	soldier := Unit{
		world: w,
		char: '@',
		colour: 'g',
		class: "soldier",
		weapon: "rifle",
		x: 5,
		y: 5,
		hp: 4,
		pc: true,
	}

	w.AddUnit(&soldier)
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
			for _, unit := range w.units {
				if unit.x == click.X && unit.y == click.Y {
					w.selection = unit
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

		if w.selection != nil {
			if key == "w" { w.selection.TryMove( 0, -1) }
			if key == "a" { w.selection.TryMove(-1,  0) }
			if key == "s" { w.selection.TryMove( 0,  1) }
			if key == "d" { w.selection.TryMove( 1,  0) }
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
