package game

import (
	"fmt"
	"os"
)

func logf(format_string string, args ...interface{}) {
	s := fmt.Sprintf(format_string, args...)
	if len(s) == 0 {
		return
	}
	if s[len(s) - 1] != '\n' {
		s += "\n"
	}
	fmt.Fprintf(os.Stderr, s)
}

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
