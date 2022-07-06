package gui

import "math/rand"

var layout [][]SpriteID

func generateSpace(size int) {

	// initialize with basic floor

	layout = make([][]SpriteID, size)
	for x := 0; x < size; x++ {
		layout[x] = make([]SpriteID, size)
		for y := 0; y < size; y++ {

			if rand.Float32() < 0.3 {
				layout[x][y] = GRASS_FLOOR
			} else {
				layout[x][y] = DIRT_FLOOR
			}

		}
	}
}
