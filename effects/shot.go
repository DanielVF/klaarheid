package effects

import (
	"time"
	electron "../electrongrid"
)

func Shot(win *electron.Window, x1, y1, x2, y2 int) {

	points := line(x1, y1, x2, y2)

	original_contents := make(map[Point]electron.Spot)

	for _, point := range points {
		original_contents[point] = win.Get(point.X, point.Y)
	}

	null_point := Point{-1, -1}

	needs_fix := null_point

	for _, point := range points {

		win.Set(point.X, point.Y, '*', 'y')

		if needs_fix != null_point {
			win.SetPointSpot(electron.Point(needs_fix), original_contents[needs_fix])
		}

		needs_fix = point

		win.Flip()
		time.Sleep(20 * time.Millisecond)
	}

	if needs_fix != null_point {
		win.SetPointSpot(electron.Point(needs_fix), original_contents[needs_fix])
	}
}
