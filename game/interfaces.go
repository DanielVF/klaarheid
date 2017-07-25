package game

type Exister interface {
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

type Keyer interface {
	Key(key string)
}
