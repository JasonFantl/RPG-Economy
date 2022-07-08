package economy

import (
	"math/rand"
)

// implement physical attributes later
type Actor struct {
	personalValues map[Good]float64

	expectedValues        map[Good]float64
	priceSignalReactivity float64 // [0.0, 1.0] how quickly we believe new price signals to represent the average market value.
	// 0.0 -> don't update.
	// 0.5 -> half of our expected value comes from the last signal.
	// 1.0 -> last signal is the expected value.

	transacted bool
}

func NewActor() *Actor {
	actor := &Actor{
		personalValues: map[Good]float64{
			MONEY:  1.0, // ha
			FOOD:   rand.Float64()*10 + 5,
			WOOD:   rand.Float64() * 10,
			ROCKET: rand.Float64()*10 + 5,
		},
		expectedValues: map[Good]float64{
			MONEY:  1.0,
			FOOD:   1.0,
			WOOD:   1.0,
			ROCKET: 0.0,
		},
		priceSignalReactivity: 0.1,
	}

	return actor
}
