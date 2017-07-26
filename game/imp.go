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

func (i *Imp) AI() {

	target := i.World.NearestPC(i.X, i.Y)

	if target == nil {
		return
	}

	i.PathTowards(target.GetX(), target.GetY())

	var args []interface{}

	args = append(args, i.X)
	args = append(args, i.Y)
	args = append(args, target.GetX())
	args = append(args, target.GetY())

	i.World.Window.Special("laser", args)
}
