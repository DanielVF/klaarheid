package game

import (
	"math/rand"
	"time"

	electron "../electronbridge_golib"
)

type Area struct {
	World		*World
	WorldX		int
	WorldY		int
	Selection	*Object
	Objects		[][][]*Object			// 3d array for: X, Y, and any number of objects at that coordinate.
}

// So the main thing to note is that objects have their own X and Y variables,
// and it's important that those should match the Area's idea of them.

func NewArea(world *World, worldx, worldy int) *Area {

	self := Area{
		World: world,
		WorldX: worldx,
		WorldY: worldy,
	}

	self.Objects = make([][][]*Object, AREA_WIDTH)

	for i := 0; i < AREA_WIDTH; i++ {
		self.Objects[i] = make([][]*Object, AREA_HEIGHT)
	}

	self.AddRandomly("Tree", VEG_FACTION, 100)
	self.AddRandomly("Bush", VEG_FACTION, 100)
	self.AddRandomly("Rock", MINERAL_FACTION, 75)
	self.AddRandomly("Grass", VEG_FACTION, 400)
	self.AddRandomly("Sheep", BEAST_FACTION, 10)

	return &self
}

// ---------------------------------------------------------

func (self *Area) Blocked(x, y int) bool {
	for _, object := range self.Objects[x][y] {
		if object.Passable == false {
			return true
		}
	}
	return false
}

func (self *Area) Empty(x, y int) bool {
	return len(self.Objects[x][y]) == 0
}

func (self *Area) DeleteObject(o *Object) {
	for i := 0; i < len(self.Objects[o.X][o.Y]); i++ {
		if self.Objects[o.X][o.Y][i] == o {
			self.Objects[o.X][o.Y] = append(self.Objects[o.X][o.Y][:i], self.Objects[o.X][o.Y][i+1:]...)
			break
		}
	}
}

func (self *Area) Move(o *Object, tar_x, tar_y int) bool {		// Return success

	if inbounds(tar_x, tar_y) == false {
		return false
	}

	// No blocked check.

	self.DeleteObject(o)		// Makes use of o.X and o.Y

	o.X = tar_x
	o.Y = tar_y

	self.AddObject(o)

	return true
}

func (self *Area) Draw() {

	MAIN_WINDOW.Clear()

	for x := 0; x < AREA_WIDTH; x++ {
		for y := 0; y < AREA_HEIGHT; y++ {

			var object_drawn *Object
			for i := 0; i < len(self.Objects[x][y]); i++ {

				// FIXME: can we do better?

				if object_drawn == nil || (object_drawn.Passable && self.Objects[x][y][i].Passable == false) {
					object_drawn = self.Objects[x][y][i]
					object_drawn.Draw()
				}
			}
		}
	}

	if (self.Selection != nil) {
		MAIN_WINDOW.SetHighlight(self.Selection.X, self.Selection.Y)
		s := self.Selection.SelectionString()
		self.WriteSelection(s)
	}

	MAIN_WINDOW.Flip()
}

func (self *Area) AddObject(o *Object) {
	self.Objects[o.X][o.Y] = append(self.Objects[o.X][o.Y], o)
}

func (self *Area) AddRandomly(classname, faction string, count int) {

	for n := 0; n < count; n++ {

		x, y := rand.Intn(AREA_WIDTH), rand.Intn(AREA_HEIGHT)

		if self.Empty(x, y) {
			new_object := NewObject(classname, self, x, y, faction)
			self.AddObject(new_object)
		}
	}
}

func (self *Area) HandleMouse() bool {				// Return true if selection changed.

	original_selection := self.Selection

	for {
		click, err := electron.GetMousedown()
		if err != nil || inbounds(click.X, click.Y) == false {
			break
		}

		if self.Empty(click.X, click.Y) == false {
			self.Selection = self.Objects[click.X][click.Y][len(self.Objects[click.X][click.Y]) - 1]		// FIXME: do better?
		} else {
			self.Selection = nil
		}
	}

	if self.Selection == original_selection {
		return false
	} else {
		return true
	}
}

func (self *Area) Play() {

	COMBAT_LOG.Printf("Camera now at area [%d,%d]", self.WorldX, self.WorldY)

	self.Selection = nil

	for {

		var all_objects []*Object

		for x := 0; x < AREA_WIDTH; x++ {
			for y := 0; y < AREA_HEIGHT; y++ {
				for i := 0; i < len(self.Objects[x][y]); i++ {
					all_objects = append(all_objects, self.Objects[x][y][i])
				}
			}
		}

		for _, object := range all_objects {
			if object.AIFunc != nil {
				f := object.AIFunc
				f(object)
			}
		}

		self.Draw()

		for n := 0; n < 500; n += 50 {
			selection_changed := self.HandleMouse()
			if selection_changed {
				self.Draw()
			}
			time.Sleep(50 * time.Millisecond)
		}

		// Sanity check during development...
		// Do all objects have the X,Y coordinates that the Area believes they do?

		for x := 0; x < AREA_WIDTH; x++ {
			for y := 0; y < AREA_HEIGHT; y++ {
				for i := 0; i < len(self.Objects[x][y]); i++ {
					if self.Objects[x][y][i].X != x || self.Objects[x][y][i].Y != y {
						electron.Alertf("Warning: x/y mismatch from object to area")
					}
				}
			}
		}
	}
}

func (self *Area) WriteSelection(s string) {

	x := 0

	for _, char := range s {
		MAIN_WINDOW.Set(x, AREA_HEIGHT + 1, string(char), "w")
		x++
	}
}
