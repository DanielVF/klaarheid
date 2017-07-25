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
				Faction: faction,
			},
			Moves: 6,
			Actions: 1,
		},
	}
	return &ret
}

func (s *Soldier) Key(key string) {
	switch key {
	case FIRE_KEY: s.Fire()
	}
}

func (s *Soldier) Fire() {
	s.ActionsLeft -= 1
	return
}

func (s *Soldier) SelectionString() string {
	return fmt.Sprintf("Soldier (hp: %d, moves: %d, actions: %d)", s.HP, s.MovesLeft, s.ActionsLeft)
}
