package svgpath

import (
	"regexp"
	"strconv"
	"strings"
)

var validSegmentCommands []string = []string{
	"m", "M", "l", "L", "v", "V", "h", "H", "z", "Z", "c", "C", "q", "Q", "t", "T", "s", "S", "a", "A"}

func PathFromSVG(svgData string) *Path {
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
	var segments []Segment

	var end = &Position{}
	re := regexp.MustCompile(`[-+]?(?:\d*\.\d+|\d+)(?:[eE][-+]?\d+)?`)

	for _, segmentStr := range segmentsData {

		if segmentStr == "" {
			continue
		}

		nextCmd := segmentStr[:1]

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

			start := end

			switch nextCmd {
			case "l": //line
				end = &Position{end.X + coords[0], end.Y + coords[1]}
				segments = append(segments, NewLine(start, end))
				coords = coords[2:]
			case "L": //line
				end = &Position{coords[0], coords[1]}
				segments = append(segments, NewLine(start, end))
				coords = coords[2:]
			case "m": //move
				end = &Position{end.X + coords[0], end.Y + coords[1]}
				//points = append(points, cpx, cpy)
				//thisCmd = "M"
				nextCmd = "l"                          // subsequent points in this segment are treated as relative lineTo
				segments = append(segments, NewMove()) //TODO
				coords = coords[2:]
			case "M": //move
				end = &Position{coords[0], coords[1]}
				//points = append(points, cpx, cpy)
				//thisCmd = "M"
				nextCmd = "L"                          // subsequent points in this segment are treated as relative lineTo
				segments = append(segments, NewMove()) //TODO
				coords = coords[2:]
			case "h": //horizontal line
				end = &Position{end.X + coords[0], end.Y}
				segments = append(segments, NewLine(start, end))
				coords = coords[1:]
			case "H": //horizontal line
				end = &Position{coords[0], end.Y}
				segments = append(segments, NewLine(start, end))
				coords = coords[1:]
			case "v": //vertical line
				end = &Position{end.X, end.Y + coords[0]}
				segments = append(segments, NewLine(start, end))
				coords = coords[1:]
			case "V": //vertical line
				end = &Position{end.X, coords[0]}
				segments = append(segments, NewLine(start, end))
				coords = coords[1:]
			case "C": //cubic bezier
				end = &Position{coords[4], coords[5]}
				segments = append(segments, NewCubicBezier(
					start,
					&Position{coords[0], coords[1]},
					&Position{coords[2], coords[3]},
					end,
				))
				coords = coords[6:]
			case "c": //cubic bezier
				end = &Position{end.X + coords[4], end.Y + coords[5]}
				segments = append(segments, NewCubicBezier(
					start,
					&Position{start.X + coords[0], start.Y + coords[1]},
					&Position{start.X + coords[2], start.Y + coords[3]},
					end,
				))
				coords = coords[6:]
			case "Q": //quadratic bezier
				end = &Position{coords[2], coords[3]}
				segments = append(segments, NewQuadraticBezier(
					start,
					&Position{coords[0], coords[1]},
					end,
				))
				coords = coords[4:]
			case "q": //quadratic bezier
				end = &Position{end.X + coords[2], end.Y + coords[3]}
				segments = append(segments, NewQuadraticBezier(
					start,
					&Position{start.X + coords[0], start.Y + coords[1]},
					end,
				))
				coords = coords[4:]
			case "T": // Smooth quadratic Bézier curve
				lastSeg, isLastSeqQuadratic := getLastSegement(segments).(*QuadraticBezier)
				ctrl := start
				if isLastSeqQuadratic {
					ctrl = &Position{
						2*start.X - lastSeg.b.X,
						2*start.Y - lastSeg.b.Y,
					}
				}
				end = &Position{coords[0], coords[1]}
				segments = append(segments, NewQuadraticBezier(
					start,
					ctrl,
					end,
				))
				coords = coords[2:]
			case "t": // Smooth quadratic Bézier curve (relative)
				ctrl := start
				lastSeg, isLastSeqQuadratic := getLastSegement(segments).(*QuadraticBezier)
				if isLastSeqQuadratic {
					ctrl = &Position{
						2*start.X - lastSeg.b.X,
						2*start.Y - lastSeg.b.Y,
					}
				}
				end = &Position{
					start.X + coords[0],
					start.Y + coords[1],
				}
				segments = append(segments, NewQuadraticBezier(
					start,
					ctrl,
					end,
				))
				coords = coords[2:]
			case "A", "a": //elliptical arc
				panic("elliptical arc not implemented")
			case "z", "Z": //close path
				segments = append(segments, NewEmpty())
				coords = nil
			default:
				coords = nil
			}

		}

		//Do we need this?
		if nextCmd == "Z" || nextCmd == "z" {
			segments = append(segments, NewEmpty())
		}
	}

	p := &Path{
		segments: segments,
		length:   0.0,
	}
	p.length = p.calculateLength()
	return p
}

func getLastSegement(segments []Segment) Segment {
	if len(segments) == 0 {
		return nil
	}
	return segments[len(segments)-1]
}
