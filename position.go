package svgpath

import (
	"fmt"
	"math"
)

type Radians float64

func RadiansFromDegrees(degrees float64) Radians {
	return Radians(math.Mod(degrees*math.Pi/180, 2*math.Pi))
}

func (r Radians) ToDegrees() float64 {
	return float64(r) * 180 / math.Pi
}

func (r Radians) AddDegrees(degrees float64) Radians {
	return RadiansFromDegrees(r.ToDegrees() + degrees)
}

/*Position - concepually the same as a 2D Vector*/
type Position struct {
	X float64
	Y float64
}

func XUnitVector() *Position {
	return &Position{1, 0}
}

func YUnitVector() *Position {
	return &Position{0, 1}
}

/*Trajectory returns the angle in radiants of the vector from Position p to Position o*/
func (p *Position) Trajectory(o *Position) Radians {
	return Radians(math.Atan2(o.Y-p.Y, o.X-p.X))
}

func (p *Position) Distance(o *Position) float64 {
	dY := o.Y - p.Y
	dX := o.X - p.X
	return math.Sqrt(math.Pow(dY, 2) + math.Pow(dX, 2))
}

func (p *Position) String() string {
	return fmt.Sprintf("x=%f y=%f", p.X, p.Y)
}

/*Find a position given the origin position p, an angle (trajectory) and a distance away from p*/
func (p *Position) GetOffsetPosition(trajectory Radians, distance float64) *Position {
	return &Position{
		X: p.X + distance*math.Cos(float64(trajectory)),
		Y: p.Y + distance*math.Sin(float64(trajectory)),
	}
}

func (p *Position) GetAdded(o *Position) *Position {
	return &Position{
		X: p.X + o.X,
		Y: p.Y + o.Y,
	}
}

func (p *Position) GetSubtracted(o *Position) *Position {
	return &Position{
		X: p.X - o.X,
		Y: p.Y - o.Y,
	}
}

/*
Rotating a position means rotating the vector pointing to it,
with the origin x:0,y:0 being the anchor point.
*/
func (p *Position) GetRotatedPosition(angle Radians) *Position {
	return &Position{
		X: p.X*math.Cos(float64(angle)) - p.Y*math.Sin(float64(angle)),
		Y: p.X*math.Sin(float64(angle)) + p.Y*math.Cos(float64(angle)),
	}
}
