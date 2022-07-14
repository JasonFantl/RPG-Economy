package economy

import (
	"math/rand"
)

// implement physical attributes later
type Actor struct {
	money   int
	markets map[Good]*marketInfo
	skills  map[Job]int
}

func NewActor() *Actor {
	d := rand.Intn(10) + 10
	actor := &Actor{
		money: 100,
		markets: map[Good]*marketInfo{
			ROCKET: {
				basePersonalValue:       float64(d),
				expectedValue:           0,
				ownedAssets:             10,
				desiredAssets:           d * 2,
				failedTransactionThresh: 5,
				priceSignalReactivity:   1.0,
			},
			FOOD: {
				basePersonalValue:       rand.Float64()*5 + 10,
				expectedValue:           0,
				ownedAssets:             0,
				desiredAssets:           rand.Intn(5) + 2,
				failedTransactionThresh: 4,
				priceSignalReactivity:   1.0,
			},
		},
		skills: map[Job]int{
			FARMER:     rand.Intn(10),
			WOODWORKER: 10,
		},
	}

	return actor
}

func (actor *Actor) Update() {
	if rand.Float64() > 0.5 { // simulates time between activities
		return
	}

	if rand.Float64() < 0.1 && actor.markets[FOOD].ownedAssets > 0 { // eat food sometimes
		actor.markets[FOOD].ownedAssets--
	}

	actor.runJob(FARMER)
	// actor.runMarket(ROCKET)
	actor.runMarket(FOOD)

}
