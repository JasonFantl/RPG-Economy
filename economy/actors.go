package economy

import (
	"math"
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

	assets        map[Good]int
	desiredAssets map[Good]int
}

func NewActor() *Actor {
	actor := &Actor{
		personalValues: map[Good]float64{
			ROCKET: rand.Float64()*10 + 5,
		},
		expectedValues: map[Good]float64{
			ROCKET: rand.Float64()*10 + 5,
		},
		priceSignalReactivity: 1.0,
		assets: map[Good]int{
			MONEY:  100,
			ROCKET: 10,
		},
		desiredAssets: map[Good]int{
			ROCKET: rand.Intn(10) + 20,
		},
	}

	return actor
}

func (actor *Actor) Update() {

	// desired equation
	actor.personalValues[ROCKET] = float64(actor.desiredAssets[ROCKET] - actor.assets[ROCKET])

	if rand.Float64() > 0.1 { // simulates time between activities
		return
	}
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
	return int(math.Floor(math.Min(actor.expectedValues[good], actor.personalValues[good])))
}

func (actor *Actor) willingSellPrice(good Good) int {
	return int(math.Ceil(math.Max(actor.expectedValues[good], actor.personalValues[good])))
}

func (actor *Actor) isBuyer(good Good) bool {
	return actor.expectedValues[good] < actor.personalValues[good]
}

func (actor *Actor) isSeller(good Good) bool {
	return actor.expectedValues[good] > actor.personalValues[good]
}

func (actor *Actor) canBuy(cost int) bool {
	// return true
	return actor.assets[MONEY] >= cost
}

func (actor *Actor) canSell(good Good) bool {
	// return true
	return actor.assets[good] > 0
}
