package game

import (
	"math/rand"
	"time"

	electron "../electronbridge_golib"
)

// ---------------------------------------------------------

type Area struct {
	World		*World
	X			int
	Y			int
	Selection	Thinger
	Objects		[]Thinger
}

func NewArea(world *World, x, y int) *Area {

	self := Area{
		World: world,
		X: x,
		Y: y,
	}

	for n := 0; n < 100; n++ {

		x := rand.Intn(AREA_WIDTH)
		y := rand.Intn(AREA_HEIGHT)

		if self.Blocked(x, y) == false {
			self.AddObject(NewTree(&self, x, y, VEG_FACTION))
		}
	}

	return &self
}

// ---------------------------------------------------------

func (self *Area) Blocked(x, y int) bool {
	for _, object := range self.Objects {
		if object.GetX() == x && object.GetY() == y {
			return true
		}
	}
	return false
}

func (self *Area) Draw() {

	MAIN_WINDOW.Clear()

	for _, object := range self.Objects {

		object.Draw()

		if object == self.Selection {
			MAIN_WINDOW.SetHighlight(object.GetX(), object.GetY())
		}
	}

	if (self.Selection != nil) {
		s := self.Selection.SelectionString()
		self.WriteSelection(s)
	}

	MAIN_WINDOW.Flip()
}

func (self *Area) AddObject(object Thinger) {
	self.Objects = append(self.Objects, object)
}

func (self *Area) Play() {

	COMBAT_LOG.Printf("Camera now at area [%d,%d]", self.X, self.Y)

	self.Selection = nil

	for {

		// Mouse events...

		for {
			click, err := electron.GetMousedown()
			if err != nil {
				break
			}
			self.Selection = nil
			for _, object := range self.Objects {
				if object.GetX() == click.X && object.GetY() == click.Y {
					self.Selection = object
				}
			}
		}

		for _, object := range self.Objects {
			if m, ok := object.(Mobber); ok {
				m.AI()
			}
		}

		self.Draw()
		time.Sleep(100 * time.Millisecond)
	}
}

func (self *Area) WriteSelection(s string) {
	for x := 0; x < len(s); x++ {
		MAIN_WINDOW.Set(x, AREA_HEIGHT + 1, s[x], 'w')
	}
}
