package game

func maybe_plural(s string, n int) string {
	if n != 1 {
		s += "s"
	}
	return s
}
