package main

import (
	"time"

	engine "./goroguego"
)

const (
	WIDTH = 50
	HEIGHT = 30
)

func main() {
	w := engine.NewWindow("World", "renderer.html", WIDTH, HEIGHT, 15, 20, 100, true)

	player_x := 5
	player_y := 5

	w.Clear()
	w.Set(player_x, player_y, '@', 'g')
	w.Flip()

	for {

		key := engine.GetKeypress()
		if key == "" {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if key == "a" && player_x > 0 {
			player_x -= 1
		}

		if key == "d" && player_x < WIDTH - 1 {
			player_x += 1
		}

		if key == "w" && player_y > 0 {
			player_y -= 1
		}

		if key == "s" && player_y < HEIGHT - 1 {
			player_y += 1
		}

		w.Clear()
		w.Set(player_x, player_y, '@', 'g')
		w.Flip()
	}
}
