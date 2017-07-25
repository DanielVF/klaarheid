package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type Object struct {

	// Things we get from the class json...

	Class			string		`json:"class"`
	Char			string		`json:"char"`
	Colour			string		`json:"colour"`
	Weapon			string		`json:"weapon"`
	Faction			string		`json:"faction"`
	HP				int			`json:"hp"`
	Moves			int			`json:"moves"`
	Actions			int			`json:"actions"`

	// Things that must be set at object creation...

	World			*World
	X				int
	Y				int

	// Other...

	MovesLeft		int
	ActionsLeft		int
}

func (o *Object) SelectionString() string {
	return fmt.Sprintf("- %s (%dhp), %s, %dm, %da", o.Class, o.HP, o.Weapon, o.Moves, o.Actions)
}

func (o *Object) TryMove(x, y int) {

	if o.Moves <= 0 {
		return
	}

	tar_x := o.X + x
	tar_y := o.Y + y

	if o.World.InBounds(tar_x, tar_y) && o.World.Blocked(tar_x, tar_y) == false {
		o.X = tar_x
		o.Y = tar_y
		o.Moves -= 1
	}
}

func (o *Object) Reset() {
	o.MovesLeft = o.Moves
	o.ActionsLeft = o.Actions
}

// -------------------------------------------------------------------

type World struct {
	window		*electron.Window
	width		int
	height		int
	selection	*Object
	objects		[]*Object
}

func (w *World) InBounds(x, y int) bool {
	if x >= 0 && x < w.width && y >= 0 && y < w.height {
		return true
	} else {
		return false
	}
}

func (w *World) Blocked(x, y int) bool {
	for _, object := range w.objects {
		if object.X == x && object.Y == y {
			return true
		}
	}
	return false
}

func (w *World) Draw() {

	w.window.Clear()

	for _, object := range w.objects {
		w.window.Set(object.X, object.Y, object.Char[0], object.Colour[0])		// char and colour are strings here, but the engine wants bytes
		if object == w.selection {
			w.window.SetHighlight(object.X, object.Y)
		}
	}

	if (w.selection != nil) {
		s := w.selection.SelectionString()
		w.WriteSelection(s)
	}

	w.window.Flip()
}

func (w *World) AddObject(object *Object) {
	w.objects = append(w.objects, object)
}

func (w *World) Game() {
	w.MakeLevel()
	w.PlayLevel()
}

func (w *World) MakeLevel() {
	w.objects = nil
	w.selection = nil

	w.AddObject(object_from_name("soldier", w, 1, 1))
	w.AddObject(object_from_name("soldier", w, 2, 2))
	w.AddObject(object_from_name("imp", w, WORLD_WIDTH - 2, WORLD_HEIGHT - 2))
}

func (w *World) PlayerTurn() {
	for {

		// Deal with mouse events...

		for {
			click, err := electron.GetMousedown()
			if err != nil {
				break
			}
			w.selection = nil
			for _, object := range w.objects {
				if object.X == click.X && object.Y == click.Y {
					w.selection = object
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
			w.selection = nil
		}

		if key == "Tab" {
			w.Tab()
		}

		if w.selection != nil && w.selection.Faction == "good" && key != "" {
			if key == "w" { w.selection.TryMove( 0, -1) }
			if key == "a" { w.selection.TryMove(-1,  0) }
			if key == "s" { w.selection.TryMove( 0,  1) }
			if key == "d" { w.selection.TryMove( 1,  0) }
		}

		w.Draw()

		time.Sleep(20 * time.Millisecond)
	}
}

func (w *World) ComputerTurn() {
	return
}


func (w *World) PlayLevel() {
	for {
		w.PlayerTurn()
		w.ComputerTurn()
	}
}

func (w *World) WriteSelection(s string) {
	for x := 0; x < len(s); x++ {
		w.window.Set(x, w.height + 1, s[x], 'w')
	}
}

func (w *World) Tab() {

	if w.selection == nil || w.selection.Faction != "good" {
		for _, object := range w.objects {
			if object.Faction == "good" {
				w.selection = object
				return
			}
		}
		return
	}

	index := -1

	for i, object := range w.objects {
		if object == w.selection {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	for _, object := range w.objects[index + 1:] {
		if object.Faction == "good" {
			w.selection = object
			return
		}
	}

	for _, object := range w.objects[:index] {
		if object.Faction == "good" {
			w.selection = object
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

func object_from_name(name string, world *World, x, y int) *Object {
	filename := fmt.Sprintf("game/classes/%s.json", name)

	j, err := ioutil.ReadFile(filename)
	if err != nil {
		log(err.Error())
	}

	var new_object Object

	err = json.Unmarshal(j, &new_object)
	if err != nil {
		log(err.Error())
	}

	new_object.X = x
	new_object.Y = y

	new_object.World = world
	return &new_object
}

func App() {

	world := World{
		window: electron.NewWindow("World", "pages/grid.html", WORLD_WIDTH, WORLD_HEIGHT + 2, 15, 20, 100, true),
		width: WORLD_WIDTH,
		height: WORLD_HEIGHT,
	}

	world.Game()
}
