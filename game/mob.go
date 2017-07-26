package game

import "fmt"

type Mobber interface {
	Reset()
	TryMove(x, y int)
	Key(key string)
	AI()
}

// The base Mob object should implement minimal satisfying methods for Mobber.
// Thus, any type that embeds Mob is automatically a Mobber.

type Mob struct {
	Thing
	MovesLeft			int
	ActionsLeft			int
	Moves				int
	Actions				int
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

func (s *Mob) TryMove(x, y int) {

	if s.MovesLeft <= 0 {
		return
	}

	success := s.MoveIfNotBlocked(x, y)

	if success {
		s.MovesLeft -= 1
	}
}

func (s *Mob) Key(key string) {
	return
}

func (s *Mob) AI() {
	return
}
