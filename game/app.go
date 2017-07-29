package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
	electron "../electronbridge_golib"
)

var MAIN_WINDOW *electron.GridWindow
var COMBAT_LOG *electron.TextWindow

var BASE_CLASSES = make(map[string]*Object)

func LoadClasses() {

	var base_classes []*Object

	infile, err := os.Open("game/classes.json")
	if err != nil {
		panic("LoadClasses: " + err.Error())
	}

	dec := json.NewDecoder(infile)

	err = dec.Decode(&base_classes)
	if err != nil {
		panic("LoadClasses: " + err.Error())
	}

	for _, base := range base_classes {

		if base.AI != "" {
			f, ok := AI_Lookup[base.AI]
			if !ok {
				panic(fmt.Sprintf("LoadClasses: Unknown AI '%s'", base.AI))
			} else {
				base.AIFunc = f
			}
		}

		BASE_CLASSES[base.Class] = base
	}
}

func App() {

	rand.Seed(time.Now().UTC().UnixNano())

	COMBAT_LOG = electron.NewTextWindow("Combat Log", "pages/log.html", 600, 400, true)
	MAIN_WINDOW = electron.NewGridWindow("Area", "pages/grid.html", AREA_WIDTH, AREA_HEIGHT + 2, 12, 20, 100, false)

	LoadClasses()

	world := NewWorld(1, 1)
	world.Play()
}
