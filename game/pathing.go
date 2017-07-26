package game

const NO_PATH = 999999

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

	for x := 0; x < w.Width; x++ {
		for y := 0; y < w.Height; y++ {
			ret[x][y] = NO_PATH
		}
	}

	if w.InBounds(x, y) == false {
		return ret
	}

	ret[x][y] = 0

	var seeds []Point
	var next_seeds []Point

	next_seeds = append(next_seeds, Point{x, y})

	dist := 0

	for {

		dist++

		seeds = next_seeds
		next_seeds = nil

		for _, seed := range(seeds) {
			for _, neigh := range w.Neighbours(seed.X, seed.Y) {
				if ret[neigh.X][neigh.Y] == NO_PATH && blocks[neigh.X][neigh.Y] == false {
					ret[neigh.X][neigh.Y] = dist
					next_seeds = append(next_seeds, Point{neigh.X, neigh.Y})
				}
			}
		}

		if len(next_seeds) == 0 {
			return ret
		}
	}
}

func (w *World) NearestPC(i, j int) Thinger {

	distances := w.DistanceMap(i, j)

	best_dist := NO_PATH
	var best_object Thinger = nil

	for _, object := range w.Objects {

		if object.IsPlayerControlled() == false {
			continue
		}

		x := object.GetX()
		y := object.GetY()

		if x == i && y == j {
			return object
		}

		for _, neigh := range w.Neighbours(x, y) {
			if distances[neigh.X][neigh.Y] < best_dist {
				best_object = object
				best_dist = distances[neigh.X][neigh.Y]
			}
		}
	}

	return best_object
}
