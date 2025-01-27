package svgpath

import (
	"fmt"
	"math"
)

type Segment interface {
	PositionAt(float64) *Position
	Length() float64
	StartPosition() *Position
	EndPosition() *Position
}

type Path struct {
	segments []Segment
	length   float64
}

func (p *Path) calculateLength() float64 {
	length := 0.0
	for _, s := range p.segments {
		length += s.Length()
	}
	return length
}

func (p *Path) GetPositionAtLength(l float64) *Position {

	if len(p.segments) == 0 {
		return nil
	}

	if p.length <= l {
		//return last sement's end coordinates
		return p.segments[len(p.segments)-1].EndPosition()
	}

	reachedSegmentIdx := 0
	for i, s := range p.segments {
		if l >= s.Length() {
			l -= s.Length()
			reachedSegmentIdx = i + 1
		} else {
			break
		}
	}

	s := p.segments[reachedSegmentIdx]

	if l < 0.01 {
		return s.StartPosition()
	}

	return s.PositionAt(l)

}

func (p *Path) String() string {
	return fmt.Sprintf("Path(%s)[l=%f]", p.segments, p.length)
}

type Empty struct{}

func NewEmpty() *Empty {
	return &Empty{}
}

func (e *Empty) PositionAt(t float64) *Position {
	return nil
}

func (e *Empty) Length() float64 {
	return 0
}

func (e *Empty) StartPosition() *Position {
	return nil
}

func (e *Empty) EndPosition() *Position {
	return nil
}

func (e *Empty) String() string {
	return "Empty"
}

type Move struct{}

func NewMove() *Move {
	return &Move{}
}

func (m *Move) PositionAt(t float64) *Position {
	return nil
}

func (m *Move) Length() float64 {
	return 0
}

func (m *Move) StartPosition() *Position {
	return nil
}

func (m *Move) EndPosition() *Position {
	return nil
}

func (m *Move) String() string {
	return "Move"
}

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

	lineLength := getLineLength(l.start, l.end)

	if lineLength < math.Pow(1, -10) {
		return l.start
	}

	if l.start.X == l.end.X {
		if l.end.Y < l.start.Y {
			t = -t
		}
		return &Position{
			X: l.start.X,
			Y: t,
		}
	}

	m := (l.end.Y - l.start.Y) / (l.end.X - l.start.X)
	mult := 1.0
	if l.end.X < l.start.X {
		mult = -1.0
	}
	run := math.Sqrt((t*t)/(1+m*m)) * mult
	rise := m * run

	if math.Abs(l.start.Y-l.start.Y-m*(l.start.X-l.start.X)) < math.Pow(1, -10) {
		return &Position{
			X: l.start.X + run,
			Y: l.start.Y + rise,
		}
	}

	u := ((l.start.X-l.start.X)*(l.end.X-l.start.X) + (l.start.Y-l.start.Y)*(l.end.Y-l.start.Y)) / (lineLength * lineLength)
	i := &Position{
		X: l.start.X + u*(l.end.X-l.start.X),
		Y: l.start.Y + u*(l.end.Y-l.start.Y),
	}
	pRise := getLineLength(l.start, i)
	pRun := math.Sqrt(t*t - pRise*pRise)
	adjustedRun := math.Sqrt((pRun*pRun)/(1+m*m)) * mult
	adjustedRise := m * adjustedRun

	return &Position{
		X: i.X + adjustedRun,
		Y: i.Y + adjustedRise,
	}

}

func (l *Line) Length() float64 {
	return l.length
}

func (l *Line) StartPosition() *Position {
	return l.start
}

func (l *Line) EndPosition() *Position {
	return l.end
}

func (l *Line) String() string {
	return fmt.Sprintf("Line(%s %s)[l=%f]", l.start, l.end, l.length)
}

/*EllipticalArc to be implemented*/
type EllipticalArc struct {
	start  *Position
	end    *Position
	c      *Position
	r      *Position
	theta  float64
	psi    float64
	length float64
}

func (e *EllipticalArc) PositionAt(t float64) *Position {
	cosPsi := math.Cos(e.psi)
	sinPsi := math.Sin(e.psi)
	pt := &Position{
		X: e.r.X * math.Cos(e.theta),
		Y: e.r.Y * math.Sin(e.theta),
	}
	return &Position{
		X: e.c.X + (pt.X*cosPsi - pt.Y*sinPsi),
		Y: e.c.Y + (pt.X*sinPsi + pt.Y*cosPsi),
	}
}

func (e *EllipticalArc) calculateLength() float64 {
	return 0.0
}

func (e *EllipticalArc) Length() float64 {
	return 0.0
}

func (e *EllipticalArc) StartPosition() *Position {
	return nil
}

func (e *EllipticalArc) EndPosition() *Position {
	return nil
}

func (e *EllipticalArc) String() string {
	return fmt.Sprintf("EllipticalArc(%s %s)[l=%f]", "", "", e.length)
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

func (q *QuadraticBezier) Length() float64 {
	return q.length
}

func (q *QuadraticBezier) StartPosition() *Position {
	return q.start
}

func (q *QuadraticBezier) EndPosition() *Position {
	return q.b
}

func (q *QuadraticBezier) String() string {
	return fmt.Sprintf("QuadraticBezier(%s %s %s)[l=%f]", q.start, q.b, q.c, q.length)
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
	bez.lookupTable = bez.generateLookupTable()
	bez.length = bez.lookupTable.distValues[len(bez.lookupTable.distValues)-1]
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

func (c *CubicBezier) PositionAt(dist float64) *Position {
	return c.PositionAtT(c.ApproximateT(dist))
}

func (c *CubicBezier) PositionAtT(t float64) *Position {
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
	tValueSampleStep := 1.0 / 2000 //Using 500 samples
	tValues := []float64{}
	distValues := []float64{}
	distAcc := 0.0
	previousPoint := c.start
	for t := 0.0; t <= 1.0; t += tValueSampleStep {
		tValues = append(tValues, t)

		pos := c.PositionAtT(t)
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

func (c *CubicBezier) Length() float64 {
	return c.length
}

func (c *CubicBezier) StartPosition() *Position {
	return c.start
}

func (c *CubicBezier) EndPosition() *Position {
	return c.d
}

func (c *CubicBezier) String() string {
	return fmt.Sprintf("CubicBezier(%s %s %s %s)[l=%f]", c.start, c.b, c.c, c.d, c.length)
}
