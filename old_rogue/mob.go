package game

import (
	"fmt"
	"math/rand"
	"time"
)

type Mobber interface {
	Reset()
	TryMove(x, y int)			bool
	PathTowards(x, y int)
	Key(key string)
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

func (s *Mob) SelectionString() string {
	if s.IsPlayerControlled() {
		return fmt.Sprintf("%s (hp: %d, moves: %d, actions: %d)", s.Class, s.HP, s.MovesLeft, s.ActionsLeft)
	} else {
		return fmt.Sprintf("%s (hp: %d)", s.Class, s.HP)
	}
}

func (s *Mob) Reset() {
	s.MovesLeft = s.Moves
	s.ActionsLeft = s.Actions
}

func (s *Mob) TryMove(x, y int) bool {

	if s.MovesLeft <= 0 {
		return false
	}

	success := s.MoveIfNotBlocked(x, y)

	if success {
		s.MovesLeft -= 1
	}

	return success
}

func (m *Mob) PathTowards(x, y int) {

	w := m.World

	distance_map := w.DistanceMap(x, y)

	max_moves := m.MovesLeft

	for n := 0; n < max_moves; n++ {

		best_dx := 0
		best_dy := 0
		best_tar_dist := NO_PATH

		for _, neigh := range w.Neighbours(m.X, m.Y) {
			if distance_map[neigh.X][neigh.Y] < best_tar_dist || (distance_map[neigh.X][neigh.Y] <= best_tar_dist && rand.Intn(2) == 0) {
				best_dx = neigh.X - m.X
				best_dy = neigh.Y - m.Y
				best_tar_dist = distance_map[neigh.X][neigh.Y]
			}
		}

		if best_dx != 0 || best_dy != 0 {

			success := m.TryMove(best_dx, best_dy)
			w.Draw()
			time.Sleep(50 * time.Millisecond)

			if success == false {
				return
			}
		}
	}
}

func (s *Mob) Key(key string) {
	return
}

func (s *Mob) AI() {
	return
}
