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
	w.Clear()

	player_x := 5
	player_y := 5

	for {
		player_x += 1

		if player_x >= WIDTH {
			player_x = 0
		}

		w.Clear()
		w.Set(player_x, player_y, '@', 'g')
		w.Flip()

		time.Sleep(150 * time.Millisecond)
	}
}
