package game

const (
	WORLD_WIDTH = 72
	WORLD_HEIGHT = 32
)

const (
	ESC_KEY = "Escape"
	TAB_KEY = "Tab"

	LEFT_KEY = "a"
	RIGHT_KEY = "d"
	UP_KEY = "w"
	DOWN_KEY = "s"

	FIRE_KEY = "f"
	TURN_END_KEY = "t"
)

const (
	PLAYER_FACTION = "Army of Light"
	DEMON_FACTION = "Demonic Horde"
	INANIMATE_FACTION = "Inanimate Objects"
)

type Point struct {
	X			int
	Y			int
}
