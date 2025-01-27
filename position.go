package svgpath

import (
	"fmt"
	"math"
)

/*Position - concepually the same as a 2D Vector*/
type Position struct {
	X float64
	Y float64
}

func (p *Position) Distance(o *Position) float64 {
	dY := o.Y - p.Y
	dX := o.X - p.X
	return math.Sqrt(math.Pow(dY, 2) + math.Pow(dX, 2))
}

func (p *Position) String() string {
	return fmt.Sprintf("x=%f y=%f", p.X, p.Y)
}
