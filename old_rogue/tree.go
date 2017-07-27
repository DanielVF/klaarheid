package game

type Tree struct {
	Thing
}

func NewTree(w *World, x, y int, faction string) *Tree {

	ret := Tree{
		Thing: Thing{
			World: w,
			X: x,
			Y: y,
			HP: 10,
			Char: 'T',
			Colour: 'G',
			Faction: faction,
			Class: "Tree",
		},
	}

	return &ret
}
