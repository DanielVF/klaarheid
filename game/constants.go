package game

const (
	AREA_WIDTH = 72
	AREA_HEIGHT = 32
)

const (
	UP = 0
	RIGHT = 1
	DOWN = 2
	LEFT = 3
)

const (
	VEG_FACTION = "Vegetation"
	ORC_FACTION = "Orcs"
)

type Point struct {
	X			int
	Y			int
}

type Vector struct {
	Dx			int
	Dy			int
}
