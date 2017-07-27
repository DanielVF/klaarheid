package game

// ---------------------------------------------------------

type World struct {
	Width		int
	Height		int
	ViewX		int
	ViewY		int
	Areas		[][]*Area
}

func NewWorld(width, height int) *World {

	self := World{
		Width: width,
		Height: height,
		ViewX: 0,
		ViewY: 0,
	}

	self.Areas = make([][]*Area, width)

	for x := 0; x < width; x++ {

		self.Areas[x] = make([]*Area, height)

		for y := 0; y < height; y++ {

			self.Areas[x][y] = NewArea(&self, x, y)
		}
	}

	return &self
}

func (self *World) Play() {
	for {
		self.Areas[self.ViewX][self.ViewY].Play()
	}
}
