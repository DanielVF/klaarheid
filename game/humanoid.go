package game

type Humanoid struct {
	Mob
}

func (self *Humanoid) AI() {
	vec := random_direction()
	self.TryMove(vec.Dx, vec.Dy)
}

func NewOrc(area *Area, x, y int, faction string) *Humanoid {
	ret := Humanoid{
		Mob: Mob{
			Thing: Thing{
				Area: area,
				X: x,
				Y: y,
				HP: 4,
				Char: 'o',
				Colour: 'r',
				Faction: faction,
				Class: "Orc",
			},
		},
	}
	return &ret
}
