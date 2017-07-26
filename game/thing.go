package game

import "fmt"

type Thinger interface {
	Draw()
	GetX()					int
	GetY()					int
	IsPlayerControlled()	bool
	SelectionString()		string
	GetClass()				string
}

// The base Thing object should implement minimal satisfying methods for Thinger.
// Thus, any type that embeds Thing is automatically a Thinger.

type Thing struct {
	World				*World
	X					int
	Y					int
	HP					int
	Char				byte
	Colour				byte
	Faction				string
	Class				string
}

func (t *Thing) GetX() int {
	return t.X
}

func (t *Thing) GetY() int {
	return t.Y
}

func (t *Thing) GetHP() int {
	return t.HP
}

func (t *Thing) GetClass() string {
	return t.Class
}

func (t *Thing) IsPlayerControlled() bool {
	return t.Faction == PLAYER_FACTION
}

func (t *Thing) SelectionString() string {
	return fmt.Sprintf("%s (%d HP)", t.Class, t.HP)
}

func (t *Thing) Draw() {
	t.World.Window.Set(t.X, t.Y, t.Char, t.Colour)
}

// Other useful methods used by types that "inherit" from Thing...

func (t *Thing) MoveIfNotBlocked(x, y int) bool {

	tar_x := t.X + x
	tar_y := t.Y + y

	if t.World.InBounds(tar_x, tar_y) && t.World.Blocked(tar_x, tar_y) == false {
		t.X = tar_x
		t.Y = tar_y
		return true
	}

	return false
}
