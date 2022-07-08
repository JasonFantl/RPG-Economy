package economy

import (
	"fmt"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func DrawGraphs(screen *ebiten.Image) {
	graphExpectedValues(FOOD, screen)
	graphOwnership(FOOD, screen)
	graphDemand(FOOD, screen)
	graphSupply(FOOD, screen)

	graphPersonalValues(FOOD, screen)

}

func graphExpectedValues(good Good, screen *ebiten.Image) {
	jump := 1

	buckets := make(map[int]int)
	// lets add the 0 index for the graph
	buckets[0] = 0
	for _, a := range actors {
		bucketIndex := int(a.expectedValues[good] / float64(jump))
		buckets[bucketIndex]++
	}

	points := make([][]float64, 0)

	for k, v := range buckets {
		points = append(points, []float64{float64(k * jump), float64(v)})
	}

	drawGraph(points, jump, 1, 30.0, 315.0, 15.0, 15.0, fmt.Sprintf("%s Expected Values", good), screen)
}

func graphPersonalValues(good Good, screen *ebiten.Image) {
	jump := 1

	buckets := make(map[int]int)
	// lets add the 0 index for the graph
	buckets[0] = 0
	for _, a := range actors {
		bucketIndex := int(a.personalValues[good] / float64(jump))
		buckets[bucketIndex]++
	}

	points := make([][]float64, 0)

	for k, v := range buckets {
		points = append(points, []float64{float64(k * jump), float64(v)})
	}

	drawGraph(points, jump, 1, 30.0, 115.0, 15.0, 10.0, fmt.Sprintf("%s Personal Values", good), screen)
}

func graphOwnership(good Good, screen *ebiten.Image) {
	jump := 1

	buckets := make(map[int]int)
	// lets add the 0 index for the graph
	buckets[0] = 0
	for _, a := range actors {
		bucketIndex := a.assets[FOOD] / jump
		buckets[bucketIndex]++
	}

	points := make([][]float64, 0)

	for k, v := range buckets {
		points = append(points, []float64{float64(k * jump), float64(v)})
	}

	drawGraph(points, jump, 1, 310.0, 115.0, 15.0, 10.0, fmt.Sprintf("%s Ownership", good), screen)
}

func graphDemand(good Good, screen *ebiten.Image) {
	jump := 1

	buckets := make(map[int]int)

	for _, a := range actors {
		bucketIndex := a.willingBuyPrice(good) / jump
		buckets[bucketIndex]++
	}

	points := make([][]float64, 0)
	for k, v := range buckets {
		points = append(points, []float64{float64(k * jump), float64(v)})
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i][0] > points[j][0]
	})

	summedPoints := make([][]float64, 0)
	for i := 0; i < len(points)-1; i++ {
		points[i+1][1] += points[i][1]
		for j := points[i][0]; j > points[i+1][0]; j -= float64(jump) {
			summedPoints = append(summedPoints, []float64{j, points[i][1]})
		}
	}
	summedPoints = append(summedPoints, points[len(points)-1])
	summedPoints = append(summedPoints, []float64{points[0][0] - float64(jump), 0})

	drawGraph(summedPoints, jump, 1, 310.0, 315.0, 15.0, 10.0, fmt.Sprintf("%s Demand", good), screen)
}

func graphSupply(good Good, screen *ebiten.Image) {
	jump := 1

	buckets := make(map[int]int)

	for _, a := range actors {
		bucketIndex := a.willingSellPrice(good) / jump
		buckets[bucketIndex]++
	}

	points := make([][]float64, 0)
	for k, v := range buckets {
		points = append(points, []float64{float64(k * jump), float64(v)})
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i][0] < points[j][0]
	})

	summedPoints := make([][]float64, 0)
	for i := 0; i < len(points)-1; i++ {
		points[i+1][1] += points[i][1]
		for j := points[i][0]; j < points[i+1][0]; j += float64(jump) {
			summedPoints = append(summedPoints, []float64{j, points[i][1]})
		}
	}
	summedPoints = append(summedPoints, points[len(points)-1])
	summedPoints = append(summedPoints, []float64{points[0][0] - float64(jump), 0})

	drawGraph(summedPoints, jump, 1, 510.0, 315.0, 15.0, 10.0, fmt.Sprintf("%s Supply", good), screen)
}

// values should be in the format {{1, 4}, {5, 8}}. Note if jumpX is not 1, then the x's must be multiples of jumpX
func drawGraph(points [][]float64, jumpXAxis, jumpYAxis int, drawXOff, drawYOff, drawXZoom, drawYZoom float64, title string, screen *ebiten.Image) {
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

	// title
	ebitenutil.DebugPrintAt(screen, title, int(drawXOff), int(drawYOff+drawYZoom))

	// X axis
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff+drawXZoom*float64(maxX-minX+2), drawYOff, color.White)
	for i := int(minX); i <= int(maxX+1); i += jumpXAxis {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), int(drawXOff+drawXZoom*(float64(i-int(minX)))), int(drawYOff))
	}
	// Y axis
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff, drawYOff-drawYZoom*float64(maxY-minY+1), color.White)
	for i := int(minY); i <= int(maxY+1); i += jumpYAxis {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), int(drawXOff-drawXZoom), int(drawYOff-drawYZoom*(float64(i-int(minY))+1.0)))
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
