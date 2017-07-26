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
