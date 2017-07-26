package effects

import (
	"time"
	electron "../../electronbridge"
)

func Shot(win *electron.Window, x1, y1, x2, y2 int) {

	points := line(x1, y1, x2, y2)

	// Reverse the slice if needed to make the start be x1, y1:

	if points[0] != (Point{x1, y1}) && points[len(points) - 1] == (Point{x1, y1}) {
		for left, right := 0, len(points) - 1; left < right; left, right = left + 1, right - 1 {
			points[left], points[right] = points[right], points[left]
		}
	}

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
