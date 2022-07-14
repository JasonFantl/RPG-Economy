package economy

import (
	"fmt"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func DrawGraphs(good Good, drawXOff, drawYOff float64, screen *ebiten.Image) {
	graphPersonalAndExpectedValues(good, 1, 6, drawXOff, drawYOff, 15.0, 15.0, screen)
	// graphSupplyDemand(good, screen)
	graphDesiredAndActualValues(good, 1, 4, drawXOff+300.0, drawYOff, 15.0, 15.0, screen)
	drawPlot(good, 1, 50, drawXOff+1000, drawYOff, 15.0, 15.0, "corr", screen)
}

func GraphWealth(screen *ebiten.Image) {
	jump := 20

	buckets := make(map[int]int)

	// lets add the 0 index for the graph
	buckets[0] = 0
	for a := range actors {
		bucketIndex := a.money / jump
		buckets[bucketIndex]++
	}

	points := make([][]float64, 0)

	for k, v := range buckets {
		points = append(points, []float64{float64(k * jump), float64(v), 0})
	}

	drawGraph(points, jump, 4, 100.0, 800.0, 15.0, 15.0, "Wealth", screen)
}

func graphPersonalAndExpectedValues(good Good, jumpXAxis, jumpYAxis int, drawXOff, drawYOff, drawXZoom, drawYZoom float64, screen *ebiten.Image) {

	buckets := []map[int]int{make(map[int]int), make(map[int]int)}

	// lets add the 0 index for the graph
	buckets[0][0] = 0
	for a := range actors {
		bucketIndex := int(a.markets[good].adjustedPersonalValue / float64(jumpXAxis))
		buckets[0][bucketIndex]++
		bucketIndex = int(a.markets[good].expectedValue / float64(jumpXAxis))
		buckets[1][bucketIndex]++
	}

	points := make([][]float64, 0)

	for i, bucket := range buckets {
		for k, v := range bucket {
			points = append(points, []float64{float64(k * jumpXAxis), float64(v), float64(i)})
		}
	}

	drawGraph(points, jumpXAxis, jumpYAxis, drawXOff, drawYOff, drawXZoom, drawYZoom, fmt.Sprintf("Personal and Expected for %s", good), screen)
}

func graphDesiredAndActualValues(good Good, jumpXAxis, jumpYAxis int, drawXOff, drawYOff, drawXZoom, drawYZoom float64, screen *ebiten.Image) {
	jump := 1

	buckets := make(map[int]int)
	bucket2s := make(map[int]int)
	// lets add the 0 index for the graph
	buckets[0] = 0
	bucket2s[0] = 0

	for a := range actors {
		bucketIndex := int(float64(a.markets[good].desiredAssets) / float64(jump))
		buckets[bucketIndex]++
		bucket2Index := int(float64(a.markets[good].ownedAssets) / float64(jump))
		bucket2s[bucket2Index]++
	}

	points := make([][]float64, 0)

	for k, v := range buckets {
		points = append(points, []float64{float64(k * jump), float64(v), 0})
	}
	for k, v := range bucket2s {
		points = append(points, []float64{float64(k * jump), float64(v), 1})
	}

	drawGraph(points, jumpXAxis, jumpYAxis, drawXOff, drawYOff, drawXZoom, drawYZoom, fmt.Sprintf("Desired and Actual %s", good), screen)
}

func graphSupplyDemand(good Good, screen *ebiten.Image) {
	jump := 1

	buckets := make(map[int]int)

	for a := range actors {
		bucketIndex := int(a.markets[good].adjustedPersonalValue) / jump
		buckets[bucketIndex]++
	}

	points := make([][]float64, 0)
	for k, v := range buckets {
		points = append(points, []float64{float64(k * jump), float64(v)})
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i][0] < points[j][0]
	})

	sum := 0.0
	summedPointsL := make([][]float64, 0)
	for i := 0; i < len(points); i++ {
		sum += points[i][1]
		summedPointsL = append(summedPointsL, []float64{points[i][0], sum, 0})
		for j := points[i][0] + float64(jump); i != len(points)-1 && j < points[i+1][0]; j += float64(jump) {
			summedPointsL = append(summedPointsL, []float64{j, sum, 0})
		}
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i][0] > points[j][0]
	})

	sum = 0
	summedPointsR := make([][]float64, 0)
	for i := 0; i < len(points); i++ {
		sum += points[i][1]
		summedPointsR = append(summedPointsR, []float64{points[i][0], sum, 1})
		for j := points[i][0] - float64(jump); i != len(points)-1 && j > points[i+1][0]; j -= float64(jump) {
			summedPointsR = append(summedPointsR, []float64{j, sum, 1})
		}
	}

	SD := append(summedPointsR, summedPointsL...)
	SD = append(SD, []float64{0, 0, 0})
	drawGraph(SD, jump, 10, 100.0, 400.0, 15.0, 15.0, fmt.Sprintf("Supply v Demand of %s", good), screen)
}

// values should be in the format {{1, 4}, {5, 8}}. Note if jumpX is not 1, then the x's must be multiples of jumpX
func drawGraph(points [][]float64, jumpXAxis, jumpYAxis int, drawXOff, drawYOff, drawXZoom, drawYZoom float64, title string, screen *ebiten.Image) {
	sort.Slice(points, func(i, j int) bool {
		return points[i][0] < points[j][0]
	})

	ODXZ, ODYZ := drawXZoom, drawYZoom
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
	ebitenutil.DebugPrintAt(screen, title, int(drawXOff), int(drawYOff+ODYZ))

	// X axis
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff+drawXZoom*float64(maxX-minX+2), drawYOff, color.White)
	for i := int(minX); i <= int(maxX+1); i += jumpXAxis {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), int(drawXOff+drawXZoom*(float64(i-int(minX)))), int(drawYOff))
	}
	// Y axis
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff, drawYOff-drawYZoom*float64(maxY-minY+1), color.White)
	for i := int(minY); i <= int(maxY+1); i += jumpYAxis {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), int(drawXOff-ODXZ), int(drawYOff-drawYZoom*float64(i-int(minY))-ODYZ))
	}

	// histogram
	for i := 0; i < len(points); i++ {
		x := drawXOff + drawXZoom*(points[i][0]-minX)
		y := drawYOff
		w := drawXZoom * float64(jumpXAxis)
		h := -drawYZoom * (points[i][1] - minY)

		c := color.RGBA{0, 0, 200, 100}
		if len(points[i]) > 2 {
			n := int(points[i][2]) + 1
			c = color.RGBA{uint8(200 * (n % 2)), uint8(200 * ((n / 2) % 2)), uint8(200 * ((n / 4) % 2)), 100}
		}

		ebitenutil.DrawRect(screen, x, y, w, h, c)
	}

}

func drawPlot(goodX Good, jumpXAxis, jumpYAxis int, drawXOff, drawYOff, drawXZoom, drawYZoom float64, title string, screen *ebiten.Image) {
	minX, maxX := 0, 0
	minY, maxY := 0, 0

	ODXZ, ODYZ := drawXZoom, drawYZoom
	drawXZoom /= float64(jumpXAxis)
	drawYZoom /= float64(jumpYAxis)

	for actor := range actors {
		x := actor.skills[goodToJob[goodX]]
		// x := actor.markets[goodX].ownedAssets
		y := actor.money
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}

	averages := make(map[int]float64)
	counts := make(map[int]int)
	for actor := range actors {
		x := actor.skills[goodToJob[goodX]]
		// x := actor.markets[goodX].ownedAssets
		y := actor.money

		averages[x] += float64(y)
		counts[x]++
		ebitenutil.DrawRect(screen, drawXOff+float64(x-minX)*drawXZoom, drawYOff+float64(-y+minY)*drawYZoom, 3, 3, color.White)
	}

	for x, y := range averages {
		y /= float64(counts[x])
		ebitenutil.DrawRect(screen, drawXOff+float64(x-minX)*drawXZoom, drawYOff+(-y+float64(minY))*drawYZoom, 5, 5, color.RGBA{250, 20, 20, 255})
	}

	// title
	ebitenutil.DebugPrintAt(screen, title, int(drawXOff), int(drawYOff+ODYZ))

	// X axis
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff+drawXZoom*float64(maxX-minX+2), drawYOff, color.White)
	for i := int(minX); i <= int(maxX+1); i += int(jumpXAxis) {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), int(drawXOff+drawXZoom*(float64(i-int(minX)))), int(drawYOff))
	}
	// Y axis
	ebitenutil.DrawLine(screen, drawXOff, drawYOff, drawXOff, drawYOff-drawYZoom*float64(maxY-minY+1), color.White)
	for i := int(minY); i <= int(maxY+1); i += int(jumpYAxis) {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), int(drawXOff-ODXZ), int(drawYOff-drawYZoom*float64(i-int(minY))-ODYZ))
	}
}
