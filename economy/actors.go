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

	// cooldown should later be replaced with limited assets
	cooldown int // stops everyone from buying from 1 buyer or selling to 1 seller
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
	}

	return actor
}

func (actor *Actor) Update() {
	if rand.Float64() < 0.1 { // simulates time between activities
		// select another actor
		var buyer, seller *Actor = nil, nil

		if isBuyer(actor, ROCKET) { // buying
			buyer = actor
			// look for a seller
			for otherActor := range actors { // rely on random iteration
				if isSeller(otherActor, ROCKET) {
					seller = otherActor
					break
				}
			}
		} else if isSeller(actor, ROCKET) { // selling
			seller = actor
			// look for a buyer
			for otherActor := range actors { // rely on random iteration
				if isBuyer(otherActor, ROCKET) {
					buyer = otherActor
					break
				}
			}
		}

		if buyer != nil && seller != nil { // if a  pair has been found
			offerPrice := seller.expectedValues[ROCKET] // we would use max, but we know sellers exp > per
			// need some wiggle room 									V
			willingBuyPrice := math.Min(buyer.expectedValues[ROCKET], buyer.personalValues[ROCKET])
			if offerPrice <= willingBuyPrice { // transaction made
				buyer.cooldown, seller.cooldown = 50, 50 // this can alter S and D (not seen in graphs) by altering how many sellers and buyers there are
				seller.expectedValues[ROCKET] += actor.priceSignalReactivity
				buyer.expectedValues[ROCKET] -= actor.priceSignalReactivity
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
}
