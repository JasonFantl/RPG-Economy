package economy

import (
	"fmt"
	"math"
	"math/rand"
)

var actors map[*Actor]bool

func Initialize() {
	actors = make(map[*Actor]bool)
	for i := 0; i < 20; i++ {
		actors[NewActor()] = true
	}
}

func Update() {

	sellers := make(map[*Actor]bool)
	buyers := make(map[*Actor]bool)

	for actor := range actors {
		actor.cooldown--
		if actor.expectedValues[ROCKET] < actor.personalValues[ROCKET] {
			buyers[actor] = true
		} else {
			sellers[actor] = true
		}
	}
	fmt.Printf("BUYERS: %d, SELLERS: %d\n", len(buyers), len(sellers))

	for actor := range actors {
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
				fmt.Println("matched")
				offerPrice := seller.expectedValues[ROCKET] // we would use max, but we know sellers exp > per
				// need some wiggle room 									V
				willingBuyPrice := math.Min(buyer.expectedValues[ROCKET], buyer.personalValues[ROCKET])
				if offerPrice <= willingBuyPrice { // transaction made
					buyer.cooldown, seller.cooldown = 10, 10
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
}

func ChangePVals(d float64) {
	for actor := range actors {
		actor.personalValues[ROCKET] += d
	}
}

func isBuyer(actor *Actor, good Good) bool {
	return (actor.expectedValues[good] < actor.personalValues[good]) && actor.cooldown <= 0
}

func isSeller(actor *Actor, good Good) bool {
	return actor.expectedValues[good] > actor.personalValues[good] && actor.cooldown <= 0
}
