package game

import "time"

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

	w := i.World

	target := w.NearestPC(i.X, i.Y)

	if target == nil {
		return
	}

	distance_map := w.DistanceMap(target.GetX(), target.GetY())

	max_moves := i.MovesLeft

	for n := 0; n < max_moves; n++ {

		best_dx := 0
		best_dy := 0
		best_tar_dist := NO_PATH

		for _, neigh := range w.Neighbours(i.X, i.Y) {
			if distance_map[neigh.X][neigh.Y] < best_tar_dist {
				best_dx = neigh.X - i.X
				best_dy = neigh.Y - i.Y
				best_tar_dist = distance_map[neigh.X][neigh.Y]
			}
		}

		if best_dx != 0 || best_dy != 0 {

			success := i.TryMove(best_dx, best_dy)
			w.Draw()
			time.Sleep(50 * time.Millisecond)

			// It's odd for this not to succeed, but...

			if success == false {
				return
			}
		}
	}
}
