package game

func initialise_object(o *Object, area *Area, x, y int, faction string) {
	o.Area = area
	o.X = x
	o.Y = y
	o.Faction = faction
}

func NewTree(area *Area, x, y int, faction string) *Object {

	ret := Object{
		HP: 10,
		Char: 'Y',
		Colour: 'G',
		Class: "Tree",
	}

	initialise_object(&ret, area, x, y, faction)
	return &ret
}

func NewBush(area *Area, x, y int, faction string) *Object {

	ret := Object{
		HP: 10,
		Char: '*',
		Colour: 'g',
		Class: "Bush",
	}

	initialise_object(&ret, area, x, y, faction)
	return &ret
}

func NewOrc(area *Area, x, y int, faction string) *Object {

	ret := Object{
		HP: 4,
		Char: 'o',
		Colour: 'r',
		Class: "Orc",

		AI: RandomWalk,
	}

	initialise_object(&ret, area, x, y, faction)
	return &ret
}
