package game

type Mobber interface {
	TryMove(x, y int)				bool
	AI()
}

// The base Mob object should implement minimal satisfying methods for Mobber.
// Thus, any type that embeds Mob is automatically a Mobber.

type Mob struct {
	Thing
}

func (self *Mob) TryMove(dx, dy int) bool {

	tar_x := self.X + dx
	tar_y := self.Y + dy

	if inbounds(tar_x, tar_y) && self.Area.Blocked(tar_x, tar_y) == false {
		self.X = tar_x
		self.Y = tar_y
		return true
	}

	return false
}

func (self *Mob) AI() {
	return
}
