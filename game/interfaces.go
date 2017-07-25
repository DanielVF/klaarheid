package game

type Object interface {
	GetX()					int
	GetY()					int
	IsPlayerControlled()	bool
	SelectionString()		string
	Draw()
}

type TryMover interface {
	TryMove(x, y int)
}

type Reseter interface {
	Reset()
}

type Keyser interface {
	Act(key string)
}
