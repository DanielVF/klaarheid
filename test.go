package main

import (
	"time"

	engine "./goroguego"
)

func main() {
	w := engine.NewWindow("World", "renderer.html", 50, 30, 15, 20, 100, true)
	w.Clear('.')
	w.Set(5, 5, '@')
	w.Set(49, 10, '"')
	w.Flip()

	for {
		time.Sleep(1 * time.Second)
	}
}
