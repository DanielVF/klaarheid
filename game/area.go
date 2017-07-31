package game

import (
	"math/rand"
	"time"

	electron "../electronbridge_golib"
)

type Area struct {
	World		*World
	X			int
	Y			int
	Selection	*Object
	Objects		[]*Object
	Tiles		[][]*Object
}

func NewArea(world *World, x, y int) *Area {

	self := Area{
		World: world,
		X: x,
		Y: y,
	}

	self.Tiles = make([][]*Object, AREA_WIDTH)

	for i := 0; i < AREA_WIDTH; i++ {
		self.Tiles[i] = make([]*Object, AREA_HEIGHT)
	}

	self.AddRandomly("Tree", VEG_FACTION, 100)
	self.AddRandomly("Bush", VEG_FACTION, 100)
	self.AddRandomly("Orc", ORC_FACTION, 10)

	self.AddTileRandomly("Grass", VEG_FACTION, 400)

	return &self
}

// ---------------------------------------------------------

func (self *Area) Blocked(x, y int) bool {
	for _, object := range self.Objects {
		if object.X == x && object.Y == y {
			return true
		}
	}
	return false
}

func (self *Area) Draw() {

	MAIN_WINDOW.Clear()

	for x := 0; x < AREA_WIDTH; x++ {
		for y := 0; y < AREA_HEIGHT; y++ {
			if self.Tiles[x][y] != nil {
				self.Tiles[x][y].Draw()
			}
		}
	}

	for _, object := range self.Objects {
		object.Draw()
	}

	if (self.Selection != nil) {
		MAIN_WINDOW.SetHighlight(self.Selection.X, self.Selection.Y)
		s := self.Selection.SelectionString()
		self.WriteSelection(s)
	}

	MAIN_WINDOW.Flip()
}

func (self *Area) AddRandomly(classname, faction string, count int) {

	for n := 0; n < count; n++ {

		x, y := rand.Intn(AREA_WIDTH), rand.Intn(AREA_HEIGHT)

		if self.Blocked(x, y) == false {
			new_object := NewObject(classname, self, x, y, faction)
			self.Objects = append(self.Objects, new_object)
		}
	}
}

func (self *Area) AddTileRandomly(classname, faction string, count int) {

	for n := 0; n < count; n++ {

		x, y := rand.Intn(AREA_WIDTH), rand.Intn(AREA_HEIGHT)

		if self.Tiles[x][y] == nil {
			new_object := NewObject(classname, self, x, y, faction)
			self.Tiles[x][y] = new_object
		}
	}
}

func (self *Area) HandleMouse() bool {				// Return true if selection changed.

	original_selection := self.Selection

	for {
		click, err := electron.GetMousedown()
		if err != nil {
			break
		}

		// Start by selecting the tile (possibly nil) but replace that selection if there's a better object here.

		self.Selection = self.Tiles[click.X][click.Y]

		for _, object := range self.Objects {
			if object.X == click.X && object.Y == click.Y {
				self.Selection = object
			}
		}
	}

	if self.Selection == original_selection {
		return false
	} else {
		return true
	}
}

func (self *Area) Play() {

	COMBAT_LOG.Printf("Camera now at area [%d,%d]", self.X, self.Y)

	self.Selection = nil

	for {
		for _, object := range self.Objects {
			if object.AIFunc != nil {
				f := object.AIFunc
				f(object)
			}
		}

		for x := 0; x < AREA_WIDTH; x++ {
			for y := 0; y < AREA_HEIGHT; y++ {
				if self.Tiles[x][y] != nil {
					if self.Tiles[x][y].AIFunc != nil {
						f := self.Tiles[x][y].AIFunc
						f(self.Tiles[x][y])
					}
				}
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
	}
}

func (self *Area) WriteSelection(s string) {

	x := 0

	for _, char := range s {
		MAIN_WINDOW.Set(x, AREA_HEIGHT + 1, string(char), "w")
		x++
	}
}
