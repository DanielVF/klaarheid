package game

type BigPlant struct {
	Thing
}

func NewTree(area *Area, x, y int, faction string) *BigPlant {

	ret := BigPlant{
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

func NewBush(area *Area, x, y int, faction string) *BigPlant {

	ret := BigPlant{
		Thing: Thing{
			Area: area,
			X: x,
			Y: y,
			HP: 10,
			Char: '*',
			Colour: 'g',
			Faction: faction,
			Class: "Bush",
		},
	}

	return &ret
}
