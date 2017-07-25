package game

import (
	"fmt"
	"os"
	"time"

	// effects "./effects"
	electron "../electronbridge"
)

const (
	WORLD_WIDTH = 50
	WORLD_HEIGHT = 28
)

// -------------------------------------------------------------------

type World struct {
	Window		*electron.Window
	Width		int
	Height		int
	Selection	Object
	Objects		[]Object
}

func (w *World) InBounds(x, y int) bool {
	if x >= 0 && x < w.Width && y >= 0 && y < w.Height {
		return true
	} else {
		return false
	}
}

func (w *World) Blocked(x, y int) bool {
	for _, object := range w.Objects {
		if object.GetX() == x && object.GetY() == y {
			return true
		}
	}
	return false
}

func (w *World) Draw() {

	w.Window.Clear()

	for _, object := range w.Objects {

		object.Draw()

		if object == w.Selection {
			w.Window.SetHighlight(object.GetX(), object.GetY())
		}
	}

	if (w.Selection != nil) {
		s := w.Selection.SelectionString()
		w.WriteSelection(s)
	}

	w.Window.Flip()
}

func (w *World) AddObject(object Object) {
	w.Objects = append(w.Objects, object)
}

func (w *World) Game() {
	w.MakeLevel()
	w.PlayLevel()
}

func (w *World) MakeLevel() {
	w.Objects = nil
	w.Selection = nil

	w.AddObject(NewSoldier(w, 1, 1))
	w.AddObject(NewSoldier(w, 2, 2))
}

func (w *World) PlayerTurn() {

	for _, object := range w.Objects {
		if object.IsPlayerControlled() {
			if r, ok := object.(Reseter); ok {
				r.Reset()
			}
		}
	}

	for {

		// Deal with mouse events...

		for {
			click, err := electron.GetMousedown()
			if err != nil {
				break
			}
			w.Selection = nil
			for _, object := range w.Objects {
				if object.GetX() == click.X && object.GetY() == click.Y {
					w.Selection = object
				}
			}
		}

		// Deal with key events...

		var key = ""

		// For now, we just skip all but the last keypress on the queue...

		for {
			nextkey, err := electron.GetKeypress()
			if err != nil {
				break
			}
			key = nextkey
		}

		if key == "Escape" {
			w.Selection = nil
		}

		if key == "Tab" {
			w.Tab()
		}

		if w.Selection != nil && w.Selection.IsPlayerControlled() && key != "" {
			if tm, ok := w.Selection.(TryMover); ok {
				if key == "w" { tm.TryMove( 0, -1) }
				if key == "a" { tm.TryMove(-1,  0) }
				if key == "s" { tm.TryMove( 0,  1) }
				if key == "d" { tm.TryMove( 1,  0) }
			} else {
				log("Player controlled unit was not a TryMover")
			}
		}

		w.Draw()

		time.Sleep(20 * time.Millisecond)
	}
}

func (w *World) ComputerTurn() {

	for _, object := range w.Objects {
		if object.IsPlayerControlled() == false {
			if r, ok := object.(Reseter); ok {
				r.Reset()
			}
		}
	}

	// TODO: behave
}

func (w *World) PlayLevel() {
	for {
		w.PlayerTurn()
		w.ComputerTurn()
	}
}

func (w *World) WriteSelection(s string) {
	for x := 0; x < len(s); x++ {
		w.Window.Set(x, w.Height + 1, s[x], 'w')
	}
}

func (w *World) Tab() {

	if w.Selection == nil || w.Selection.IsPlayerControlled() == false {
		for _, object := range w.Objects {
			if object.IsPlayerControlled() {
				w.Selection = object
				return
			}
		}
		return
	}

	index := -1

	for i, object := range w.Objects {
		if object == w.Selection {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	for _, object := range w.Objects[index + 1:] {
		if object.IsPlayerControlled() {
			w.Selection = object
			return
		}
	}

	for _, object := range w.Objects[:index] {
		if object.IsPlayerControlled() {
			w.Selection = object
			return
		}
	}
}

// -------------------------------------------------------------------

func log(s string) {
	if len(s) == 0 {
		return
	}
	if s[len(s) - 1] != '\n' {
		s += "\n"
	}
	fmt.Fprintf(os.Stderr, s)
}

func App() {

	world := World{
		Window: electron.NewWindow("World", "pages/grid.html", WORLD_WIDTH, WORLD_HEIGHT + 2, 15, 20, 100, true),
		Width: WORLD_WIDTH,
		Height: WORLD_HEIGHT,
	}

	world.Game()
}
