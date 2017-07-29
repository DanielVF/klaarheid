package game

import (
	"fmt"
)

type Object struct {
	Class				string			`json:"class"`
	Char				string			`json:"char"`
	Colour				string			`json:"colour"`
	AI					string			`json:"ai"`
	HP					int				`json:"hp"`

	Area				*Area
	X					int
	Y					int
	Faction				string
	AIFunc				func(*Object)
}

func NewObject(class string, area *Area, x, y int, faction string) *Object {

	o := new(Object)

	base_class, ok := BASE_CLASSES[class]

	if !ok {
		panic(fmt.Sprintf("Class %s not known", class))
	}

	*o = *base_class

	o.Area = area
	o.X = x
	o.Y = y
	o.Faction = faction

	return o
}

func (self *Object) SelectionString() string {
	return fmt.Sprintf("%s / %s (hp: %d)", self.Class, self.Faction, self.HP)
}

func (self *Object) Draw() {

	var char byte
	var colour byte

	// The bridge library wants bytes not strings for colour and char.

	if self.Colour != "" {
		colour = self.Colour[0]
	} else {
		colour = 'w'
	}

	if self.Char != "" {
		char = self.Char[0]
	} else {
		char = '?'
	}

	MAIN_WINDOW.Set(self.X, self.Y, char, colour)
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
