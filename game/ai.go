package game

import (
	"math/rand"
)

var AI_Lookup = map[string]func(*Object){
	"PlantGrow": PlantGrow,
	"GrassWalk": GrassWalk,
}

func GrassWalk(self *Object) {
	vec := random_direction()

	tar_x := self.X + vec.Dx
	tar_y := self.Y + vec.Dy

	self.BlockableMove(tar_x, tar_y)

	for _, target := range self.Area.Objects[self.X][self.Y] {
		if target.Class == "Grass" {
			self.Attack(target)
			break
		}
	}
}

func PlantGrow(self *Object) {
	for _, neigh := range neighbours(self.X, self.Y) {
		if self.Area.Empty(neigh.X, neigh.Y) {
			if rand.Intn(1000) == 0 {
				self.Area.AddObject(NewObject(self.Class, self.Area, neigh.X, neigh.Y, self.Faction))
			}
		}
	}
}
