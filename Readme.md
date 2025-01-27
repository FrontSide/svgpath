svgpath (Go)
============

A Go Library for dealing with paths which can be parsed from SVG Path elements.

A path consistes of a number of segments, each of which may be one of the following:
    - straight line
    - elliptical arc (not yet supported)
    - quadratic bézier curve
    - cubic bézier curve 

The library supports:
    - calculating the length of the path up until a given spacial distance value
    - calculating the point on a path at a given distance value 

Example
```
path := PathFromSVG("m 1633.8176,1077.4212 c 0,0 18.4277,-511.56464 -14.7423,-535.31585 -32.2488,-23.0917 -318.7995,-9.13506 -318.7995,-9.13506 H 830.3691 c 0,0 -182.43438,-54.81052 -189.80546,-129.7182 -3.97298,-40.37463 -16.58496,-164.43147 -16.58496,-164.43147 0,0 -60.81148,-89.5238 -180.59162,-95.00485 C 323.6069,138.33472 -4.4064367,134.68069 -4.4064367,134.68069")

path.length //2391.198781219853
path.GetPositionAtLength(0) //&Position{X: 1633.8176, Y: 1077.4212}
path.GetPositionAtLength(91.9549) //&Position{X: 1636.4624661312378, Y: 985.5047820629673}
```
