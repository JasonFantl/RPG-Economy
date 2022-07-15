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
	actor := &Actor{
		money: 100,
		markets: map[Good]*marketInfo{
			FOOD:  newMarket(0, 3, 6, 15.0),
			WOOD:  newMarket(3, 2, 10, 10.0),
			HOUSE: newMarket(0, 1, 2, 15.0),
		},
		skills: map[Job]int{
			FARMER:     rand.Intn(10),
			WOODWORKER: 10,
			BUILDER:    rand.Intn(10) / 9,
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
	if rand.Float64() < 0.01 && actor.markets[HOUSE].ownedAssets > 0 { // house burns down
		actor.markets[HOUSE].ownedAssets--
	}

	actor.runJob(FARMER)
	actor.runJob(BUILDER)

	actor.runMarket(FOOD)
	actor.runMarket(WOOD)
	actor.runMarket(HOUSE)

}
