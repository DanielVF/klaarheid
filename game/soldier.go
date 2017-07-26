package game

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
				Class: "Soldier",
			},
			Moves: 6,
			Actions: 1,
		},
	}
	return &ret
}

func (s *Soldier) Key(key string) {
	switch key {
	case FIRE_KEY: s.PlayerFire()
	}
}

func (s *Soldier) PlayerFire() {
	if s.ActionsLeft <= 0 {
		return
	}
	s.ActionsLeft -= 1
	return
}
