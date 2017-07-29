package game

var AI_Lookup = map[string]func(*Object){
	"RandomWalk": RandomWalk,
}

func RandomWalk(self *Object) {
	vec := random_direction()
	self.TryMove(vec.Dx, vec.Dy)
}
