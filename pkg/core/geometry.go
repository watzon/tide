package core

type Point struct {
	X, Y int
}

type Size struct {
	Width, Height int
}

type Rect struct {
	Min, Max Point
}

func NewRect(x, y, width, height int) Rect {
	return Rect{
		Min: Point{X: x, Y: y},
		Max: Point{X: x + width, Y: y + height},
	}
}

func (r Rect) Size() Size {
	return Size{
		Width:  r.Max.X - r.Min.X,
		Height: r.Max.Y - r.Min.Y,
	}
}

func (r Rect) Contains(p Point) bool {
	return p.X >= r.Min.X && p.X < r.Max.X &&
		p.Y >= r.Min.Y && p.Y < r.Max.Y
}
