package game

const NO_PATH = 999999

func (self *Area) BlockMap() [][]bool {

	ret := make_2d_bool_array(AREA_WIDTH, AREA_HEIGHT)

	for _, object := range self.Objects {

		x := object.GetX()
		y := object.GetY()

		if x >= 0 && x < AREA_WIDTH && y >= 0 && y < AREA_HEIGHT {
			ret[x][y] = true
		}
	}

	return ret
}

func (self *Area) DistanceMap(x, y int) [][]int {

	blocks := self.BlockMap()

	ret := make_2d_int_array(AREA_WIDTH, AREA_HEIGHT)

	for x := 0; x < AREA_WIDTH; x++ {
		for y := 0; y < AREA_HEIGHT; y++ {
			ret[x][y] = NO_PATH
		}
	}

	if inbounds(x, y) == false {
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
			for _, neigh := range neighbours(seed.X, seed.Y) {
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

func (self *Area) NearestFactionMob(faction string, i, j int) Thinger {

	distances := self.DistanceMap(i, j)

	best_dist := NO_PATH
	var best_object Thinger = nil

	for _, object := range self.Objects {

		if object.GetFaction() != faction {
			continue
		}

		x := object.GetX()
		y := object.GetY()

		if x == i && y == j {
			return object
		}

		for _, neigh := range neighbours(x, y) {
			if distances[neigh.X][neigh.Y] < best_dist {
				best_object = object
				best_dist = distances[neigh.X][neigh.Y]
			}
		}
	}

	return best_object
}
