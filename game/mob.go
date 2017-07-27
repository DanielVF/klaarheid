package game

import (
	"fmt"
	"math/rand"
	"time"
)

type Mobber interface {
	TryMove(x, y int)			bool
	MoveIfNotBlocked(x, y int)	bool
	PathTowards(x, y int)
	AI()
}

// The base Mob object should implement minimal satisfying methods for Mobber.
// Thus, any type that embeds Mob is automatically a Mobber.

type Mob struct {
	Thing
	MovesLeft					int
	ActionsLeft					int
	Moves						int
	Actions						int
}

func (self *Mob) SelectionString() string {
	return fmt.Sprintf("%s (hp: %d, moves: %d, actions: %d)", self.Class, self.HP, self.MovesLeft, self.ActionsLeft)
}

func (self *Mob) TryMove(x, y int) bool {

	if self.MovesLeft <= 0 {
		return false
	}

	success := self.MoveIfNotBlocked(x, y)

	if success {
		self.MovesLeft -= 1
	}

	return success
}

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

func (self *Mob) PathTowards(x, y int) {

	a := self.Area

	distance_map := a.DistanceMap(x, y)

	max_moves := self.MovesLeft

	for n := 0; n < max_moves; n++ {

		best_dx := 0
		best_dy := 0
		best_tar_dist := NO_PATH

		for _, neigh := range neighbours(self.X, self.Y) {
			if distance_map[neigh.X][neigh.Y] < best_tar_dist || (distance_map[neigh.X][neigh.Y] <= best_tar_dist && rand.Intn(2) == 0) {
				best_dx = neigh.X - self.X
				best_dy = neigh.Y - self.Y
				best_tar_dist = distance_map[neigh.X][neigh.Y]
			}
		}

		if best_dx != 0 || best_dy != 0 {

			success := self.TryMove(best_dx, best_dy)
			a.Draw()
			time.Sleep(50 * time.Millisecond)

			if success == false {
				return
			}
		}
	}
}

func (self *Mob) AI() {
	return
}
