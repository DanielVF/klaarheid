package game

import "fmt"

type Soldier struct {
	Thing
	MovesLeft			int
	ActionsLeft			int
}

func NewSoldier(w *World, x, y int) *Soldier {
	ret := Soldier{Thing: Thing{World: w, X: x, Y: y, HP: 4, Char: '@', Colour: 'g'}}
	return &ret
}

func (s *Soldier) IsPlayerControlled() bool {
	return true
}

func (s *Soldier) SelectionString() string {
	m := maybe_plural("move", s.MovesLeft)
	a := maybe_plural("action", s.ActionsLeft)
	return fmt.Sprintf("Soldier (%d hp, %d %s, %d %s)", s.HP, s.MovesLeft, m, s.ActionsLeft, a)
}

func (s *Soldier) TryMove(x, y int) {

	if s.MovesLeft <= 0 {
		return
	}

	success := s.MoveIfNotBlocked(x, y)

	if success {
		s.MovesLeft -= 1
	}
}

func (s *Soldier) Reset() {
	s.MovesLeft = 6
	s.ActionsLeft = 1
}
