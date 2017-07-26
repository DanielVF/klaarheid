package game

import (
	"fmt"
	"os"
	"time"

	// effects "./effects"
	electron "../electronbridge"
)

const (
	WORLD_WIDTH = 72
	WORLD_HEIGHT = 32
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

const (
	PLAYER_FACTION = "Army of Light"
	DEMON_FACTION = "Demonic Horde"
)

// -------------------------------------------------------------------

type Point struct {
	X			int
	Y			int
}

type World struct {
	Window		*electron.Window
	Width		int
	Height		int
	Selection	Thinger
	Objects		[]Thinger
}

func (w *World) Neighbours(x, y int) []Point {
	var ret []Point
	if w.InBounds(x - 1, y) { ret = append(ret, Point{x - 1, y}) }
	if w.InBounds(x + 1, y) { ret = append(ret, Point{x + 1, y}) }
	if w.InBounds(x, y - 1) { ret = append(ret, Point{x, y - 1}) }
	if w.InBounds(x, y + 1) { ret = append(ret, Point{x, y + 1}) }
	return ret
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

	w.AddObject(NewSoldier(w, 1, 1, PLAYER_FACTION))
	w.AddObject(NewSoldier(w, 2, 2, PLAYER_FACTION))

	w.AddObject(NewImp(w, WORLD_WIDTH - 2, WORLD_HEIGHT - 2, DEMON_FACTION))
}

func (w *World) PlayerTurn() {

	for _, object := range w.Objects {
		if object.IsPlayerControlled() {
			if r, ok := object.(Mobber); ok {
				r.Reset()
			}
		}
	}

	w.Draw()

	for {

		must_draw := false

		// Deal with mouse events...

		for {
			click, err := electron.GetMousedown()
			if err != nil {
				break
			}
			must_draw = true
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
			must_draw = true
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

		if must_draw {
			w.Draw()
		}

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

func log(format_string string, args ...interface{}) {
	s := fmt.Sprintf(format_string, args...)
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
		Window: electron.NewWindow("World", "pages/grid.html", WORLD_WIDTH, WORLD_HEIGHT + 2, 12, 20, 100, true),
		Width: WORLD_WIDTH,
		Height: WORLD_HEIGHT,
	}

	world.Game()
}
