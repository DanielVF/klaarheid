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

const (
	ESC_KEY = "Escape"
	TAB_KEY = "Tab"

	LEFT_KEY = "a"
	RIGHT_KEY = "d"
	UP_KEY = "w"
	DOWN_KEY = "s"

	FIRE_KEY = "f"
	TURN_END_KEY = "t"
)

// -------------------------------------------------------------------

type World struct {
	Window		*electron.Window
	Width		int
	Height		int
	Selection	Thinger
	Objects		[]Thinger
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

func (w *World) AddObject(object Thinger) {
	w.Objects = append(w.Objects, object)
}

func (w *World) Game() {
	w.MakeLevel()
	w.PlayLevel()
}

func (w *World) MakeLevel() {
	w.Objects = nil
	w.Selection = nil

	w.AddObject(NewSoldier(w, 1, 1, "player"))
	w.AddObject(NewSoldier(w, 2, 2, "player"))

	w.AddObject(NewImp(w, WORLD_WIDTH - 2, WORLD_HEIGHT - 2, "demons"))
}

func (w *World) PlayerTurn() {

	for _, object := range w.Objects {
		if object.IsPlayerControlled() {
			if r, ok := object.(Mobber); ok {
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

		if key == ESC_KEY {
			w.Selection = nil
		}

		if key == TAB_KEY {
			w.Tab()
		}

		if key == TURN_END_KEY {
			return
		}

		if w.Selection != nil && w.Selection.IsPlayerControlled() && key != "" {

			if tm, ok := w.Selection.(Mobber); ok {
				if key == UP_KEY { tm.TryMove( 0, -1) }
				if key == LEFT_KEY { tm.TryMove(-1,  0) }
				if key == DOWN_KEY { tm.TryMove( 0,  1) }
				if key == RIGHT_KEY { tm.TryMove( 1,  0) }
			} else {
				log("Player controlled unit was not a Mobber")
			}

			if ks, ok := w.Selection.(Mobber); ok {
				ks.Key(key)
			} else {
				log("Player controlled unit was not a Mobber")
			}
		}

		w.Draw()

		time.Sleep(20 * time.Millisecond)
	}
}

func (w *World) ComputerTurn() {

	for _, object := range w.Objects {
		if object.IsPlayerControlled() == false {
			if r, ok := object.(Mobber); ok {
				r.Reset()
			}
		}
	}

	for _, object := range w.Objects {
		if object.IsPlayerControlled() == false {
			if a, ok := object.(Mobber); ok {
				a.AI()
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
