package game

type Thing struct {
	World				*World
	X					int
	Y					int
	HP					int
	Char				byte
	Colour				byte
}

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

func (t *Thing) GetX() int {
	return t.X
}

func (t *Thing) GetY() int {
	return t.Y
}

func (t *Thing) Draw() {
	t.World.Window.Set(t.X, t.Y, t.Char, t.Colour)
}
