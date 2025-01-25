package svgpath

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Segment interface {
	PositionAt(float64) *Position
}

var validSegmentCommands []string = []string{
	"m", "M", "l", "L", "v", "V", "h", "H", "z", "Z", "c", "C", "q", "Q", "t", "T", "s", "S", "a", "A"}

func ParseSegmentsFromSVG(svgData string) []*Segment {
	if svgData == "" {
		return nil
	}

	// Replace spaces with commas
	svgData = strings.ReplaceAll(svgData, " ", ",")

	// Add pipes before command characters
	for _, command := range validSegmentCommands {
		svgData = strings.ReplaceAll(svgData, command, "|"+command)
	}

	// Split into segments
	segmentsData := strings.Split(svgData, "|")
	var segments []*Segment

	var cpx, cpy float64
	re := regexp.MustCompile(`[-+]?(?:\d*\.\d+|\d+)(?:[eE][-+]?\d+)?`)

	for _, segmentStr := range segmentsData {
		if segmentStr == "" {
			continue
		}

		cmdChar := segmentStr[:1]

		//Parse the segments coordinates (the float numbers following a command)
		segmentPtsStr := segmentStr[1:] // Remove the command character
		segmentCoordinatesStr := re.FindAllString(segmentPtsStr, -1)

		var coords []float64
		for _, coordinateStr := range segmentCoordinatesStr {
			val, err := strconv.ParseFloat(coordinateStr, 64)
			if err != nil {
				coords = append(coords, 0)
			} else {
				coords = append(coords, val)
			}
		}

		for len(coords) > 0 {

			var points []float64
			startX, startY := cpx, cpy
			var cmd string

			switch cmdChar {
			case "l":
				cpx += coords[0]
				cpy += coords[1]
				points = append(points, cpx, cpy)
				coords = coords[2:]
				cmd = "L"
			case "L":
				cpx = coords[0]
				cpy = coords[1]
				points = append(points, cpx, cpy)
				coords = coords[2:]
			case "m":
				cpx += coords[0]
				cpy += coords[1]
				points = append(points, cpx, cpy)
				cmd = "M"
				cmdChar = "l" // subsequent points in this segment are treated as relative lineTo
				coords = coords[2:]
			case "M":
				cpx = coords[0]
				cpy = coords[1]
				points = append(points, cpx, cpy)
				cmd = "M"
				cmdChar = "L" // subsequent points in this segment are treated as relative lineTo
				coords = coords[2:]
			case "h":
				cpx += coords[0]
				points = append(points, cpx, cpy)
				coords = coords[1:]
				cmd = "L"
			case "H":
				cpx = coords[0]
				points = append(points, cpx, cpy)
				coords = coords[1:]
				cmd = "L"
			case "v":
				cpy += coords[0]
				points = append(points, cpx, cpy)
				coords = coords[1:]
				cmd = "L"
			case "V":
				cpy = coords[0]
				points = append(points, cpx, cpy)
				coords = coords[1:]
				cmd = "L"
			case "C":
				points = append(points, coords[0], coords[1], coords[2], coords[3])
				cpx = coords[4]
				cpy = coords[5]
				points = append(points, cpx, cpy)
				coords = coords[6:]
			case "c":
				points = append(points, cpx+coords[0], cpy+coords[1], cpx+coords[2], cpy+coords[3])
				cpx += coords[4]
				cpy += coords[5]
				points = append(points, cpx, cpy)
				coords = coords[6:]
				cmd = "C"
			case "Q":
				points = append(points, coords[0], coords[1])
				cpx = coords[2]
				cpy = coords[3]
				points = append(points, cpx, cpy)
				coords = coords[4:]
			case "q":
				points = append(points, cpx+coords[0], cpy+coords[1])
				cpx += coords[2]
				cpy += coords[3]
				points = append(points, cpx, cpy)
				coords = coords[4:]
				cmd = "Q"
			case "T":
				// Smooth quadratic Bézier curve
				prevCmd := getLastCommand(segments)
				var ctrlX, ctrlY float64
				if prevCmd == "Q" {
					lastSeg := segments[len(segments)-1]
					ctrlX = 2*cpx - lastSeg.Points[0]
					ctrlY = 2*cpy - lastSeg.Points[1]
				} else {
					ctrlX = cpx
					ctrlY = cpy
				}
				cpx = coords[0]
				cpy = coords[1]
				points = append(points, ctrlX, ctrlY, cpx, cpy)
				coords = coords[2:]
				cmd = "Q"
			case "t":
				// Smooth quadratic Bézier curve (relative)
				prevCmd := getLastCommand(segments)
				var ctrlX, ctrlY float64
				if prevCmd == "Q" {
					lastSeg := segments[len(segments)-1]
					ctrlX = 2*cpx - lastSeg.Points[0]
					ctrlY = 2*cpy - lastSeg.Points[1]
				} else {
					ctrlX = cpx
					ctrlY = cpy
				}
				cpx += coords[0]
				cpy += coords[1]
				points = append(points, ctrlX, ctrlY, cpx, cpy)
				coords = coords[2:]
				cmd = "Q"
			case "A", "a":
				panic("elliptical arc not implemented")
			case "z", "Z":
				segments = append(segments, &Segment{
					Command: "z",
					Points:  nil,
					Start:   nil,
					Length:  0,
				})
				coords = nil
			default:
				coords = nil
			}

			command := cmdChar
			if cmd != "" {
				command = cmd
			}

			start := &Position{X: startX, Y: startY}
			length := calcLength(start, command, points)

			//For cubic bezier curves, create a lookup table
			//converting t values to evenly distanced spacial values
			var cubicBezierLookupTable *LookupTable
			if command == "C" {
				cubicBezierLookupTable = generateCubicBezierLookupTable(start, length, points)
			}

			segments = append(segments, &Segment{
				Command:                command,
				Points:                 points,
				Start:                  &Position{X: startX, Y: startY},
				Length:                 length,
				CubicBezierLookupTable: cubicBezierLookupTable,
			})

		}

		if cmdChar == "Z" || cmdChar == "z" {
			segments = append(segments, &Segment{
				Command: "z",
			})
		}
	}

	return segments
}

func getLastCommand(segments []*Segment) string {
	if len(segments) == 0 {
		return ""
	}
	return segments[len(segments)-1].Command
}

func convertEndpointToCenterParameterization(x1, y1, x2, y2, fa, fs, rx, ry, psi float64) []float64 {
	// Placeholder function for elliptical arc conversion
	return []float64{x1, y1, x2, y2}
}

func calcLength(start *Position, cmd string, points []float64) float64 {
	switch cmd {
	case "L":
		// Line length
		return getLineLength(start, &Position{X: points[0], Y: points[1]})
	case "C":
		// Cubic Bézier curve length
		return getCubicArcLength(
			[]float64{start.X, points[0], points[2], points[4]},
			[]float64{start.Y, points[1], points[3], points[5]},
			1.0,
		)
	case "Q":
		// Quadratic Bézier curve length
		return getQuadraticArcLength(
			[]float64{start.X, points[0], points[2]},
			[]float64{start.Y, points[1], points[3]},
			1.0,
		)
	case "A":
		// Elliptical arc length (approximated)
		return getEllipticalArcLength(points)
	}

	return 0
}

func getPointOnEllipticalArc(c, r *Position, theta, psi float64) *Position {
	cosPsi := math.Cos(psi)
	sinPsi := math.Sin(psi)
	pt := &Position{
		X: r.X * math.Cos(theta),
		Y: r.Y * math.Sin(theta),
	}
	return &Position{
		X: c.X + (pt.X*cosPsi - pt.Y*sinPsi),
		Y: c.Y + (pt.X*sinPsi + pt.Y*cosPsi),
	}
}

func GetPositionAtLength(l, pathLength float64, segments []*Segment) (*Position, error) {

	if len(segments) == 0 {
		return nil, fmt.Errorf("Cannot calculate position on path with no segments")
	}

	if pathLength <= l {
		//return last sements end coordinates
		lastSegment := segments[len(segments)-1]
		return &Position{
			X: lastSegment.Points[len(lastSegment.Points)-2],
			Y: lastSegment.Points[len(lastSegment.Points)-1],
		}, nil
	}

	reachedSegmentIdx := 0
	for i, s := range segments {
		if l >= s.Length {
			l -= s.Length
			reachedSegmentIdx = i + 1
		} else {
			break
		}
	}

	s := segments[reachedSegmentIdx]
	ps := s.Points

	if l < 0.01 {
		return &Position{
			X: s.Start.X,
			Y: s.Start.Y,
		}, nil
	}

	switch s.Command {
	case "L":
		p2 := &Position{
			X: ps[0],
			Y: ps[1],
		}
		return getPointOnLine(l, s.Start, p2, nil), nil
	case "C":
		return getPointOnCubicBezier(
			s.CubicBezierLookupTable.GetClosestT(l),
			s.Start,
			&Position{
				X: ps[0],
				Y: ps[1],
			},
			&Position{
				X: ps[2],
				Y: ps[3],
			},
			&Position{
				ps[4],
				ps[5],
			},
		), nil
	case "Q":
		return getPointOnQuadraticBezier(
			l,
			s.Start,
			&Position{
				X: ps[0],
				Y: ps[1],
			},
			&Position{
				X: ps[2],
				Y: ps[3],
			}), nil
	case "A":
		c := &Position{
			X: ps[0],
			Y: ps[1],
		}
		r := &Position{
			X: ps[2],
			Y: ps[3],
		}
		theta := ps[4]
		dTheta := ps[5]
		psi := ps[6]
		theta += (dTheta * l) / s.Length
		return getPointOnEllipticalArc(c, r, theta, psi), nil
	}

	return nil, fmt.Errorf("no case match for command %s when calculating point on path", s.Command)

}
