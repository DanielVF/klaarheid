package game

type Imp struct {
	Mob
}

func NewImp(w *World, x, y int, faction string) *Imp {
	ret := Imp{
		Mob: Mob{
			Thing: Thing{
				World: w,
				X: x,
				Y: y,
				HP: 4,
				Char: 'i',
				Colour: 'r',
				Faction: faction,
				Class: "Imp",
			},
			Moves: 10,
			Actions: 1,
		},
	}
	return &ret
}
