package game

type Mobber interface {
	Reset()
	TryMove(x, y int)
	Key(key string)
	AI()
}

// The base Mob object should implement minimal satisfying methods for Mobber

type Mob struct {
	Thing
	MovesLeft			int
	ActionsLeft			int
	Moves				int
	Actions				int
}

func (s *Mob) Reset() {
	s.MovesLeft = s.Moves
	s.ActionsLeft = s.Actions
}

func (s *Mob) TryMove(x, y int) {

	if s.MovesLeft <= 0 {
		return
	}

	success := s.MoveIfNotBlocked(x, y)

	if success {
		s.MovesLeft -= 1
	}
}

func (s *Mob) Key(key string) {
	return
}

func (s *Mob) AI() {
	return
}
