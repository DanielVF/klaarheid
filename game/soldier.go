package game

import "fmt"

type Soldier struct {
	Actor
}

func NewSoldier(w *World, x, y int, faction string) *Soldier {
	ret := Soldier{
		Actor: Actor{
			Thing: Thing{
				World: w,
				X: x,
				Y: y,
				HP: 4,
				Char: '@',
				Colour: 'g',
				Faction: "player",
			},
			Moves: 6,
			Actions: 1,
		},
	}
	return &ret
}

func (s *Soldier) Act(key string) {
	switch key {
	case FIRE_KEY: s.Fire()
	}
}

func (s *Soldier) Fire() {
	s.ActionsLeft -= 1
	return
}

func (s *Soldier) SelectionString() string {
	m := maybe_plural("move", s.MovesLeft)
	a := maybe_plural("action", s.ActionsLeft)
	return fmt.Sprintf("Soldier (%d hp, %d %s, %d %s)", s.HP, s.MovesLeft, m, s.ActionsLeft, a)
}
