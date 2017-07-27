package game

import "fmt"

type Thinger interface {
	Draw()
	GetX()					int
	GetY()					int
	GetHP()					int
	SelectionString()		string
	GetFaction()			string
	GetClass()				string
}

// The base Thing object should implement minimal satisfying methods for Thinger.
// Thus, any type that embeds Thing is automatically a Thinger.

type Thing struct {
	Area				*Area
	X					int
	Y					int
	HP					int
	Char				byte
	Colour				byte
	Faction				string
	Class				string
}

func (self *Thing) GetX() int {
	return self.X
}

func (self *Thing) GetY() int {
	return self.Y
}

func (self *Thing) GetHP() int {
	return self.HP
}

func (self *Thing) GetFaction() string {
	return self.Faction
}

func (self *Thing) GetClass() string {
	return self.Class
}

func (self *Thing) SelectionString() string {
	return fmt.Sprintf("%s (hp: %d)", self.Class, self.HP)
}

func (self *Thing) Draw() {
	MAIN_WINDOW.Set(self.X, self.Y, self.Char, self.Colour)
}

// Other useful methods used by types that "inherit" from Thing...

func (self *Thing) MoveIfNotBlocked(x, y int) bool {

	tar_x := self.X + x
	tar_y := self.Y + y

	if inbounds(tar_x, tar_y) && self.Area.Blocked(tar_x, tar_y) == false {
		self.X = tar_x
		self.Y = tar_y
		return true
	}

	return false
}
