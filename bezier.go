package svgpath

import (
	"math"
)

type Line struct {
	start *Position
	end   *Position

	length float64
}

func NewLine(start, end *Position) *Line {
	l := &Line{start, end, 0}
	l.length = getLineLength(start, end)
	return l
}

func getLineLength(start, end *Position) float64 {
	return math.Sqrt((end.X-start.X)*(end.X-start.X) + (end.Y-start.Y)*(end.Y-start.Y))
}

func (l *Line) PositionAt(t float64) *Position {

	if from == nil {
		from = p1
	}

	lineLength := getLineLength(p1, p2)

	if lineLength < math.Pow(1, -10) {
		return p1
	}

	if p1.X == p2.X {
		if p2.Y < p1.Y {
			dist = -dist
		}
		return &Position{
			X: from.X,
			Y: dist,
		}
	}

	m := (p2.Y - p1.Y) / (p2.X - p1.X)
	mult := 1.0
	if p2.X < p1.X {
		mult = -1.0
	}
	run := math.Sqrt((dist*dist)/(1+m*m)) * mult
	rise := m * run

	if math.Abs(from.Y-p1.Y-m*(from.X-p1.X)) < math.Pow(1, -10) {
		return &Position{
			X: from.X + run,
			Y: from.Y + rise,
		}
	}

	u := ((from.X-p1.X)*(p2.X-p1.X) + (from.Y-p1.Y)*(p2.Y-p1.Y)) / (lineLength * lineLength)
	i := &Position{
		X: p1.X + u*(p2.X-p1.X),
		Y: p1.Y + u*(p2.Y-p1.Y),
	}
	pRise := getLineLength(from, i)
	pRun := math.Sqrt(dist*dist - pRise*pRise)
	adjustedRun := math.Sqrt((pRun*pRun)/(1+m*m)) * mult
	adjustedRise := m * adjustedRun

	return &Position{
		X: i.X + adjustedRun,
		Y: i.Y + adjustedRise,
	}

}

type EllipticalArc struct {
	length float64
}

func (e *EllipticalArc) calculateLength() float64 {
	return 0.0
}

type QuadraticBezier struct {
	start  *Position
	b      *Position
	c      *Position
	length float64
}

func NewQuadraticBezier(start, b, c *Position) *QuadraticBezier {
	q := &QuadraticBezier{start, b, c, 0}
	q.length = q.DistanceAt(1.0)
	return q
}

/*
Get the length of the curve/arc from its start to t
where t is a value between 0 and 1 denoting a position along the path.
This requires calculation as t values are not linear proportinal to distance along the curve.
See: https://acegikmo.medium.com/the-ever-so-lovely-b%C3%A9zier-curve-eb27514da3bf
*/
func (q *QuadraticBezier) DistanceAt(t float64) float64 {
	ax := q.start.X - 2*q.b.X + q.c.X
	ay := q.start.Y - 2*q.b.Y + q.c.Y
	bx := 2*q.b.X - 2*q.start.X
	by := 2*q.b.Y - 2*q.start.Y

	ca := 4 * (ax*ax + ay*ay)
	cb := 4 * (ax*bx + ay*by)
	cc := bx*bx + by*by

	if ca == 0 {
		return (t * math.Sqrt(math.Pow(q.c.X-q.start.X, 2)+math.Pow(q.c.Y-q.start.Y, 2)))
	}
	b := cb / (2 * ca)
	c := cc / ca
	u := t + b
	k := c - b*b

	uuk := 0.0
	if u*u+k > 0 {
		math.Sqrt(u*u + k)
	}

	bbk := 0.0
	if b*b+k > 0 {
		math.Sqrt(b*b + k)
	}

	term := 0.0
	if b+math.Sqrt(b*b+k) != 0 {
		term = k * math.Log(math.Abs((u+uuk)/(b+bbk)))
	}

	return (math.Sqrt(ca) / 2) * (u*uuk - b*bbk + term)
}

func (q *QuadraticBezier) PositionAt(t float64) *Position {
	qb1 := t * t
	qb2 := 2 * t * (1 - t)
	qb3 := (1 - t) * (1 - t)

	return &Position{
		X: q.c.X*qb1 + q.b.X*qb2 + q.start.X*qb3,
		Y: q.c.Y*qb1 + q.b.Y*qb2 + q.start.Y*qb3,
	}
}

type LookupTable struct {
	distValues []float64
	tValues    []float64
}

type CubicBezier struct {
	start *Position
	b     *Position
	c     *Position
	d     *Position

	length float64

	lookupTable *LookupTable
}

func NewCubicBezier(start, b, c, d *Position) *CubicBezier {
	bez := &CubicBezier{start, b, c, d, 0, nil}
	bez.length = bez.DistanceAt(1.0)
	bez.lookupTable = bez.generateLookupTable()
	return bez
}

func (c *CubicBezier) ApproximateT(dist float64) float64 {
	for idx, lookupDist := range c.lookupTable.distValues {
		if dist < lookupDist {
			if idx == 0 {
				return 0
			} else {
				diffLookupDist := lookupDist - c.lookupTable.distValues[idx-1]
				diffRelativeDist := dist - c.lookupTable.distValues[idx-1]
				diffLookupT := c.lookupTable.tValues[idx] - c.lookupTable.tValues[idx-1]
				return c.lookupTable.tValues[idx-1] + (diffLookupT*diffRelativeDist)/diffLookupDist
			}
		}
	}
	return 1.0
}

func (c *CubicBezier) PositionAt(t float64) *Position {
	cb1 := t * t * t
	cb2 := 3 * t * t * (1 - t)
	cb3 := 3 * t * (1 - t) * (1 - t)
	cb4 := (1 - t) * (1 - t) * (1 - t)

	return &Position{
		X: c.d.X*cb1 + c.c.X*cb2 + c.b.X*cb3 + c.start.X*cb4,
		Y: c.d.Y*cb1 + c.c.Y*cb2 + c.b.Y*cb3 + c.start.Y*cb4,
	}
}

func (c *CubicBezier) generateLookupTable() *LookupTable {
	tValueSampleStep := 1.0 / 500 //Using 500 samples
	tValues := []float64{}
	distValues := []float64{}
	distAcc := 0.0
	previousPoint := c.start
	for t := 0.0; t <= 1.0; t += tValueSampleStep {
		tValues = append(tValues, t)

		pos := c.PositionAt(t)
		p2pDist := previousPoint.Distance(pos)
		distAcc += p2pDist
		previousPoint = pos
		distValues = append(distValues, distAcc)
	}

	return &LookupTable{
		tValues:    tValues,
		distValues: distValues,
	}
}
