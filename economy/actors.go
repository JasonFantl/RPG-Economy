package economy

import (
	"math/rand"
)

// implement physical attributes later
type Actor struct {
	money   int
	markets map[Good]*marketInfo
}

func NewActor() *Actor {
	d := rand.Intn(10) + 10
	actor := &Actor{
		money: 100,
		markets: map[Good]*marketInfo{
			ROCKET: &marketInfo{
				basePersonalValue:       float64(d),
				expectedValue:           0,
				ownedAssets:             10,
				desiredAssets:           d * 2,
				failedTransactionThresh: 5,
				priceSignalReactivity:   1.0,
			},
		},
	}

	return actor
}

func (actor *Actor) Update() {
	if rand.Float64() > 0.5 { // simulates time between activities
		return
	}

	actor.runMarket(ROCKET)

}
