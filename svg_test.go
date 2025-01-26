package svgpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePath(t *testing.T) {

	path := PathFromSVG("m 1633.8176,1077.4212 c 0,0 18.4277,-511.56464 -14.7423,-535.31585 -32.2488,-23.0917 -318.7995,-9.13506 -318.7995,-9.13506 H 830.3691 c 0,0 -182.43438,-54.81052 -189.80546,-129.7182 -3.97298,-40.37463 -16.58496,-164.43147 -16.58496,-164.43147 0,0 -60.81148,-89.5238 -180.59162,-95.00485 C 323.6069,138.33472 -4.4064367,134.68069 -4.4064367,134.68069")

	assert.Equal(t, 2389.771398799874, path.length)

	t.Log(path)

	assert.Equal(t, &Position{1633.8176, 1077.4212}, path.GetPositionAtLength(0))
	assert.Equal(t, &Position{1633.8176, 1077.4212}, path.GetPositionAtLength(1))
	assert.Equal(t, &Position{0, 985}, path.GetPositionAtLength(91.9549))
	assert.Equal(t, &Position{0, 985}, path.GetPositionAtLength(95.8549))
	assert.Equal(t, &Position{0, 985}, path.GetPositionAtLength(2380))

	//There is a bug here, the movement along the path is not consistent, there's some back and forth jumoing
	/*it looks as though the problem is with the beziercurve position calculation. potentially incorrect coefficients or some other bug
	p1, err := GetPositionAtLength(91.9549, totalLength, segments)
	assert.Nil(t, err)

	assert.Equal(t, []float64{0, 985}, []float64{p1.X, p1.Y}) //pY should be 985 is 983

	p2, err := GetPositionAtLength(94.8549, totalLength, segments)
	assert.Nil(t, err)
	assert.Equal(t, []float64{0, 982}, []float64{p2.X, p2.Y}) // pY should be 982 is 984
	diffX, diffY := p2.X-p1.X, p2.Y-p1.Y
	assert.Equal(t, 0, diffX)
	assert.Equal(t, 0, diffY)

	p3, err := GetPositionAtLength(97.7549, totalLength, segments)
	assert.Equal(t, []float64{0, 0}, []float64{p3.X, p3.Y})
	assert.Nil(t, err)
	diffX, diffY = p3.X-p2.X, p3.Y-p2.Y
	assert.Equal(t, 0, diffX)
	assert.Equal(t, 0, diffY)

	p4, err := GetPositionAtLength(110, totalLength, segments)
	assert.Equal(t, []float64{0, 0}, []float64{p4.X, p4.Y})
	assert.Nil(t, err)

	diffX, diffY = p4.X-p3.X, p4.Y-p3.Y
	assert.Equal(t, 0, diffX)
	assert.Equal(t, 0, diffY)*/
}
