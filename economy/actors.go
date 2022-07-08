package economy

import (
	"fmt"
	"math"
	"math/rand"
)

// implement physical attributes later
type Actor struct {
	name           string
	assets         map[Good]int
	personalValues map[Good]float64

	expectedValues        map[Good]float64
	priceSignalReactivity float64 // [0.0, 1.0] how quickly we believe new price signals to represent the average market value.
	// 0.0 -> don't update.
	// 0.5 -> half of our expected value comes from the last signal.
	// 1.0 -> last signal is the expected value.

}

func NewActor() *Actor {
	actor := &Actor{
		name: fmt.Sprintf("%d", rand.Intn(1000)),
		assets: map[Good]int{
			MONEY: 100,
			FOOD:  10,
			WOOD:  rand.Intn(20),
		},
		personalValues: map[Good]float64{
			MONEY: 1.0, // ha
			FOOD:  rand.Float64()*10 + 5,
			WOOD:  rand.Float64() * 10,
		},
		expectedValues: map[Good]float64{
			MONEY: 1.0,
			FOOD:  1.0,
			WOOD:  1.0,
		},
		priceSignalReactivity: 0.1,
	}

	return actor
}

func (actor *Actor) update() {
	// randomly trade with random individuals
	if rand.Float64() < 0.1 {
		otherActor := actors[rand.Intn(len(actors))]
		if otherActor == actor { // don't trade with self
			return
		}

		// TRADING PROTOCOL

		offerPrice := actor.willingBuyPrice(FOOD)

		// stop if the buyer doesn't have the money
		if actor.assets[MONEY] < offerPrice {
			return
		}

		// ask to trade, get willingness and how much they would sell for
		otherWillTrade := otherActor.willingToSell(FOOD, offerPrice)

		// update understanding of the market
		if otherWillTrade {
			actor.updateExpectedValues(FOOD, offerPrice)
		} else {
			// on a failure to buy, we ask for what price they would sell at
			actor.updateExpectedValues(FOOD, otherActor.willingSellPrice(FOOD))
		}

		// actually make the exchange
		if otherWillTrade {
			makeExchange(actor, otherActor, FOOD, 1, offerPrice)
		}
	}
}

// willingToSell tells someone if they are willing to sell a good for the given price
func (actor *Actor) willingToSell(good Good, cost int) bool {

	if actor.assets[good] == 0 {
		return false
	}

	return actor.willingSellPrice(good) <= cost // at or above the willing sell price they will sell
}

// what price at or above this actor would be willing to buy a good
func (actor *Actor) willingSellPrice(good Good) int {
	return int(math.Ceil(math.Max(actor.expectedValues[good], actor.personalValues[good])))
}

// what price at or below this actor would be willing to buy a good
func (actor *Actor) willingBuyPrice(good Good) int {
	return int(math.Floor(math.Min(actor.expectedValues[good], actor.personalValues[good])))
}

func (actor *Actor) updateExpectedValues(good Good, cost int) {
	// next = (1 - reactivity) * last + reactivity * new
	actor.expectedValues[good] = (1.0-actor.priceSignalReactivity)*actor.expectedValues[good] + actor.priceSignalReactivity*float64(cost)
}
