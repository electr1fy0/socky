package game

const MinSnakeLength = 5

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Snake struct {
	Body      []Point
	Head      Point
	Tail      Point
	Direction Direction
	Score     int
}

func (s *Snake) Init() {
	s.Body = make([]Point, MinSnakeLength)
	s.Direction = Right
	s.Score = 0
}

func (s *Snake) shift(newHead Point) {
	n := len(s.Body)
	s.Tail = Point{s.Body[0].X, s.Body[0].Y}
	for i := 0; i < n-1; i++ {
		s.Body[i].X = s.Body[i+1].X
		s.Body[i].Y = s.Body[i+1].Y
	}
	s.Body[n-1] = newHead
	s.Head = newHead
}

func (s *Snake) Move() {
	var newHead Point
	switch s.Direction {
	case Up:
		newHead = Point{s.Head.X - 1, s.Head.Y}
	case Down:
		newHead = Point{s.Head.X + 1, s.Head.Y}

	case Left:
		newHead = Point{s.Head.X, s.Head.Y - 1}

	case Right:
		newHead = Point{s.Head.X, s.Head.Y + 1}
	}

	s.shift(newHead)

}
