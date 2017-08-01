package game

import (
	"fmt"
	"reflect"
)

type Object struct {
	Class				string			`json:"class"`
	Char				string			`json:"char"`
	Colour				string			`json:"colour"`
	AI					string			`json:"ai"`
	HP					int				`json:"hp"`
	Damage				int				`json:"damage"`
	Passable			bool			`json:"passable"`

	Area				*Area
	X					int				// Object must not set X or Y itself
	Y					int				// but rather call Area.Move()
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
	if o.AIClass != nil {
		ai := reflect.New(reflect.TypeOf(o.AIClass)).Elem().Interface()
		o.AIClass = ai.(Ai)
	}

	return o
}

func (self *Object) SelectionString() string {
	return fmt.Sprintf("%s / %s (hp: %d) [%d, %d]", self.Class, self.Faction, self.HP, self.X, self.Y)
}

func (self *Object) Draw() {
	MAIN_WINDOW.Set(self.X, self.Y, self.Char, self.Colour)
}

func (self *Object) BlockableMove(tar_x, tar_y int) bool {

	if inbounds(tar_x, tar_y) && self.Area.Blocked(tar_x, tar_y) == false {
		success := self.Area.Move(self, tar_x, tar_y)
		return success
	}

	return false
}

func (self *Object) Destroy() {
	self.Area.DeleteObject(self)
}

func (self *Object) SufferAttack(source *Object) {		// FIXME: return something or other
	self.HP -= source.Damage
	if self.HP <= 0 {
		self.Destroy()
	}
	COMBAT_LOG.Printf("%s [%d, %d] attacks %s for %d", source.Class, source.X, source.Y, self.Class, source.Damage)
}

func (self *Object) Attack(target *Object) {
	target.SufferAttack(self)
}
