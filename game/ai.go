package game

import (
	"math/rand"
)

var AI_Lookup = map[string]func(*Object){
	"TileGrow": TileGrow,
	"RandomWalk": RandomWalk,
}

func RandomWalk(self *Object) {
	vec := random_direction()
	self.TryMove(vec.Dx, vec.Dy)
}

func TileGrow(self *Object) {
	for _, neigh := range neighbours(self.X, self.Y) {
		if self.Area.Tiles[neigh.X][neigh.Y] == nil {
			if rand.Intn(1000) == 0 {
				self.Area.Tiles[neigh.X][neigh.Y] = NewObject(self.Class, self.Area, neigh.X, neigh.Y, self.Faction)
			}
		}
	}
}
