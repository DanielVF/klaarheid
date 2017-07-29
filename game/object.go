package game

import "fmt"

type Object struct {
	Area				*Area
	X					int
	Y					int
	HP					int
	Char				byte
	Colour				byte
	Faction				string
	Class				string
	AI					func(*Object)
}

func (self *Object) SelectionString() string {
	return fmt.Sprintf("%s (hp: %d)", self.Class, self.HP)
}

func (self *Object) Draw() {
	MAIN_WINDOW.Set(self.X, self.Y, self.Char, self.Colour)
}

func (self *Object) TryMove(dx, dy int) bool {

	tar_x := self.X + dx
	tar_y := self.Y + dy

	if inbounds(tar_x, tar_y) && self.Area.Blocked(tar_x, tar_y) == false {
		self.X = tar_x
		self.Y = tar_y
		return true
	}

	return false
}
