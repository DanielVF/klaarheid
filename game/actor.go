package game

type Actor struct {
	Thing
	MovesLeft			int
	ActionsLeft			int
	Moves				int
	Actions				int
}

func (s *Actor) TryMove(x, y int) {

	if s.MovesLeft <= 0 {
		return
	}

	success := s.MoveIfNotBlocked(x, y)

	if success {
		s.MovesLeft -= 1
	}
}

func (s *Actor) Reset() {
	s.MovesLeft = s.Moves
	s.ActionsLeft = s.Actions
}
