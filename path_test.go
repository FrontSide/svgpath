package svgpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {

	path := PathFromSVG("m 1633.8176,1077.4212 c 0,0 18.4277,-511.56464 -14.7423,-535.31585 -32.2488,-23.0917 -318.7995,-9.13506 -318.7995,-9.13506 H 830.3691 c 0,0 -182.43438,-54.81052 -189.80546,-129.7182 -3.97298,-40.37463 -16.58496,-164.43147 -16.58496,-164.43147 0,0 -60.81148,-89.5238 -180.59162,-95.00485 C 323.6069,138.33472 -4.4064367,134.68069 -4.4064367,134.68069")

	assert.Equal(t, 2391.198781219853, path.length)

	assert.Equal(t, &Position{X: 1633.8176, Y: 1077.4212}, path.GetPositionAtLength(0))
	assert.Equal(t, &Position{X: 1633.8530165129534, Y: 1076.4219182017573}, path.GetPositionAtLength(1))
	assert.Equal(t, &Position{X: 1636.4624661312378, Y: 985.5047820629673}, path.GetPositionAtLength(91.9549))
	assert.Equal(t, &Position{X: 1636.5547132869497, Y: 981.6058997892234}, path.GetPositionAtLength(95.8549))
	assert.Equal(t, &Position{X: 6.791577162469242, Y: 134.8109575568843}, path.GetPositionAtLength(2380))
}
