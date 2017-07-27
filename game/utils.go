package game

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

func inbounds(x, y int) bool {
	if x >= 0 && x < AREA_WIDTH && y >= 0 && y < AREA_HEIGHT {
		return true
	} else {
		return false
	}
}

func neighbours(x, y int) []Point {
	var ret []Point
	if inbounds(x - 1, y) { ret = append(ret, Point{x - 1, y}) }
	if inbounds(x + 1, y) { ret = append(ret, Point{x + 1, y}) }
	if inbounds(x, y - 1) { ret = append(ret, Point{x, y - 1}) }
	if inbounds(x, y + 1) { ret = append(ret, Point{x, y + 1}) }
	return ret
}
