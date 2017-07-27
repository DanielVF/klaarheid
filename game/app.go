package game

import (
	"math/rand"
	"time"
	electron "../electronbridge_golib"
)

var MAIN_WINDOW *electron.GridWindow
var COMBAT_LOG *electron.TextWindow

func App() {

	rand.Seed(time.Now().UTC().UnixNano())

	COMBAT_LOG = electron.NewTextWindow("Combat Log", "pages/log.html", 600, 400, true)
	MAIN_WINDOW = electron.NewGridWindow("Area", "pages/grid.html", AREA_WIDTH, AREA_HEIGHT + 2, 12, 20, 100, false)

	world := NewWorld(1, 1)
	world.Play()
}
