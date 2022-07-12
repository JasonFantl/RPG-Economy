package economy

import (
	"math"
	"math/rand"
)

// implement physical attributes later
type Actor struct {
	basePersonalValues     map[Good]float64 // if they had no goods, how much would they be willing to pay for this good?
	adjustedPersonalValues map[Good]float64 // changes based on if they have a good

	expectedValues        map[Good]float64
	priceSignalReactivity float64 // [0.0, 1.0] how quickly we believe new price signals to represent the average market value.
	// 0.0 -> don't update.
	// 0.5 -> half of our expected value comes from the last signal.
	// 1.0 -> last signal is the expected value.

	assets        map[Good]int
	desiredAssets map[Good]int
}

func NewActor() *Actor {
	actor := &Actor{
		basePersonalValues: map[Good]float64{
			FOOD:   rand.Float64()*10 + 10,
			ROCKET: rand.Float64() * 5,
		},
		adjustedPersonalValues: make(map[Good]float64),
		expectedValues: map[Good]float64{
			FOOD:   0,
			ROCKET: rand.Float64() * 10,
		},
		priceSignalReactivity: 1.0,
		assets: map[Good]int{
			MONEY:  100,
			FOOD:   10,
			ROCKET: 5,
		},
		desiredAssets: map[Good]int{
			FOOD:   3,
			ROCKET: rand.Intn(50),
		},
	}

	return actor
}

func (actor *Actor) Update() {
	if rand.Float64() > 0.1 { // simulates time between activities
		return
	}

	// desired equation
	actor.calcAdjustedValue(FOOD)
	actor.calcAdjustedValue(ROCKET)

	// select another actor
	var buyer, seller *Actor = nil, nil

	if actor.isBuyer(ROCKET) && actor.canBuy(actor.willingBuyPrice(ROCKET)) { // buying
		buyer = actor
		// look for a seller
		for otherActor := range actors { // rely on random iteration
			if otherActor.isSeller(ROCKET) && otherActor.canSell(ROCKET) {
				seller = otherActor
				break
			}
		}
	} else if actor.isSeller(ROCKET) && actor.canSell(ROCKET) { // selling
		seller = actor
		// look for a buyer
		for otherActor := range actors { // rely on random iteration
			if otherActor.isBuyer(ROCKET) {
				if otherActor.canBuy(seller.willingSellPrice(ROCKET)) {
					buyer = otherActor
					break
				}
			}
		}
	}

	if buyer != nil && seller != nil { // if a  pair has been found
		offerPrice := seller.willingSellPrice(ROCKET)
		willingBuyPrice := buyer.willingBuyPrice(ROCKET)

		if offerPrice <= willingBuyPrice { // transaction made
			seller.sellGood(ROCKET, offerPrice)
			buyer.buyGood(ROCKET, offerPrice)
		} else { // transaction failed
			seller.expectedValues[ROCKET] -= actor.priceSignalReactivity
			buyer.expectedValues[ROCKET] += actor.priceSignalReactivity
		}
	} else if buyer != nil { // buyer failed to find a seller
		buyer.expectedValues[ROCKET] += actor.priceSignalReactivity
	} else if seller != nil { // seller failed to find a buyer
		seller.expectedValues[ROCKET] -= actor.priceSignalReactivity
	}
}

func (actor *Actor) calcAdjustedValue(good Good) {
	x := float64(actor.assets[good])
	d := float64(actor.desiredAssets[good])
	b := actor.basePersonalValues[good]

	if d < 1 {
		actor.adjustedPersonalValues[good] = 0
	} else {
		actor.adjustedPersonalValues[good] = b * (1 - x/d)
	}

	// actor.adjustedPersonalValues[good] = actor.basePersonalValues[good]
}

func (actor *Actor) sellGood(good Good, cost int) {
	actor.assets[good]--
	actor.assets[MONEY] += cost
	actor.expectedValues[good] += actor.priceSignalReactivity
	// actor.personalValues[good] = float64(actor.desiredAssets[good] - actor.assets[good])
}

func (actor *Actor) buyGood(good Good, cost int) {
	actor.assets[good]++
	actor.assets[MONEY] -= cost
	actor.expectedValues[good] -= actor.priceSignalReactivity
	// actor.personalValues[good] = float64(actor.desiredAssets[good] - actor.assets[good])
}

func (actor *Actor) willingBuyPrice(good Good) int {
	return int(math.Floor(math.Min(actor.expectedValues[good], actor.adjustedPersonalValues[good])))
}

func (actor *Actor) willingSellPrice(good Good) int {
	return int(math.Ceil(math.Max(actor.expectedValues[good], actor.adjustedPersonalValues[good])))
}

func (actor *Actor) isBuyer(good Good) bool {
	return actor.expectedValues[good] < actor.adjustedPersonalValues[good]
}

func (actor *Actor) isSeller(good Good) bool {
	return actor.expectedValues[good] > actor.adjustedPersonalValues[good]
}

func (actor *Actor) canBuy(cost int) bool {
	// return true
	return actor.assets[MONEY] >= cost
}

func (actor *Actor) canSell(good Good) bool {
	// return true
	return actor.assets[good] > 0
}
