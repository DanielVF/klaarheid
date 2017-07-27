package game

type Tree struct {
	Thing
}

func NewTree(area *Area, x, y int, faction string) *Tree {

	ret := Tree{
		Thing: Thing{
			Area: area,
			X: x,
			Y: y,
			HP: 10,
			Char: 'Y',
			Colour: 'G',
			Faction: faction,
			Class: "Tree",
		},
	}

	return &ret
}
