package game

import "fmt"

type Soldier struct {
	Mob
}

func NewSoldier(w *World, x, y int, faction string) *Soldier {
	ret := Soldier{
		Mob: Mob{
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

func (s *Soldier) SelectionString() string {
	return fmt.Sprintf("Soldier (hp: %d, moves: %d, actions: %d)", s.HP, s.MovesLeft, s.ActionsLeft)
}

func (s *Soldier) Key(key string) {
	switch key {
	case FIRE_KEY: s.PlayerFire()
	}
}

func (s *Soldier) PlayerFire() {
	s.ActionsLeft -= 1
	return
}
