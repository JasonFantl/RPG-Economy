package economy

import (
	"fmt"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func DrawGraphs(screen *ebiten.Image) {
	graphExpectedValues(screen)
}

func graphExpectedValues(screen *ebiten.Image) {
	jump := 1

	buckets := make(map[int]int)
	// lets add the 0 index for the graph
	buckets[0] = 0
	for _, a := range actors {
		bucketIndex := int(a.expectedValues[FOOD] / float64(jump))
		buckets[bucketIndex]++
	}

	points := make([][]float64, 0)

	for k, v := range buckets {
		points = append(points, []float64{float64(k * jump), float64(v)})
	}

	drawGraph(points, jump, 1, 10.0, 400.0, 20.0, 20.0, screen)
}

// values should be in the format {{1, 4}, {5, 8}}. Note if jumpX is not 1, then the x's must be multiples of jumpX
func drawGraph(points [][]float64, jumpXAxis, jumpYAxis int, drawXOff, drawYOff, drawXZoom, drawYZoom float64, screen *ebiten.Image) {
	sort.Slice(points, func(i, j int) bool {
		return points[i][0] < points[j][0]
	})

	drawXZoom /= float64(jumpXAxis)
	drawYZoom /= float64(jumpYAxis)

	// find graph info
	minX, maxX := points[0][0], points[len(points)-1][0]
	minY, maxY := points[0][1], points[0][1]
	for _, point := range points {
		if point[1] < minY {
			minY = point[1]
		}
		if point[1] > maxY {
			maxY = point[1]
		}
	}

	// X axis
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff+drawXZoom*float64(maxX-minX+2), drawYOff, color.White)
	for i := int(minX); i <= int(maxX+1); i += jumpXAxis {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), int(drawXOff+drawXZoom*float64(i-int(minX))), int(drawYOff))
	}
	// Y axis
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff, drawYOff-drawYZoom*float64(maxY-minY+1), color.White)
	for i := int(minY); i <= int(maxY+1); i += jumpYAxis {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), int(drawXOff), int(drawYOff-drawYZoom*float64(i-int(minY))))
	}

	// histogram
	for i := 0; i < len(points); i++ {
		x := drawXOff + drawXZoom*(points[i][0]-minX)
		y := drawYOff
		w := drawXZoom * float64(jumpXAxis)
		h := -drawYZoom * (points[i][1] - minY)
		ebitenutil.DrawRect(screen, x, y, w, h, color.White)
	}
}
