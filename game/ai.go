package game

func RandomWalk(self *Object) {
	vec := random_direction()
	self.TryMove(vec.Dx, vec.Dy)
}
