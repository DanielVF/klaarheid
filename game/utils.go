package game

const NO_PATH = 999999

func make_2d_bool_array(width, height int) [][]bool {
	ret := make([][]bool, width)
	for x := 0; x < width; x++ {
		ret[x] = make([]bool, height)
	}
	return ret
}

func make_2d_int_array(width, height int) [][]int {
	ret := make([][]int, width)
	for x := 0; x < width; x++ {
		ret[x] = make([]int, height)
	}
	return ret
}

func (w *World) BlockMap() [][]bool {

	ret := make_2d_bool_array(w.Width, w.Height)

	for _, object := range w.Objects {

		x := object.GetX()
		y := object.GetY()

		if x >= 0 && x < w.Width && y >= 0 && y < w.Height {
			ret[x][y] = true
		}
	}

	return ret
}

func (w *World) DistanceMap(x, y int) [][]int {

	blocks := w.BlockMap()

	ret := make_2d_int_array(w.Width, w.Height)
	done := make_2d_bool_array(w.Width, w.Height)

	for x := 0; x < w.Width; x++ {
		for y := 0; y < w.Height; y++ {
			ret[x][y] = NO_PATH
		}
	}

	if w.InBounds(x, y) == false {
		return ret
	}

	ret[x][y] = 0
	done[x][y] = true

	dist := 0

	for {

		dist++

		totally_done := true

		for x := 0; x < w.Width; x++ {

			for y := 0; y < w.Height; y++ {

				if done[x][y] || blocks[x][y] {
					continue
				}

				for _, neigh := range w.Neighbours(x, y) {
					if ret[neigh.X][neigh.Y] == dist - 1 {
						ret[x][y] = dist
						done[x][y] = true
						totally_done = false
					}
				}
			}
		}

		if totally_done {	// i.e. no updates were made this loop
			break
		}
	}

	return ret
}

