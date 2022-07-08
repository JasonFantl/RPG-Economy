package economy

import (
	"math"
	"math/rand"
)

var actors []*Actor

func Initialize() {
	actors = make([]*Actor, 20)
	for i := 0; i < len(actors); i++ {
		actors[i] = NewActor()
	}
}

func Update() {

	// for primer we use a round based system

	// determine buyers and sellers
	sellers := make(map[*Actor]bool)
	buyers := make(map[*Actor]bool)

	for _, actor := range actors {
		actor.transacted = false // reset transactions for this round
		if actor.expectedValues[ROCKET] < actor.personalValues[ROCKET] {
			buyers[actor] = true
		} else {
			sellers[actor] = true
		}
	}

	// fmt.Printf("buyers: %d, sellers: %d\n", len(buyers), len(sellers))

	transactionMade := true
	for transactionMade {
		transactionMade = false

		// pair up
		remainingSellers := make([]*Actor, 0)
		remainingBuyers := make([]*Actor, 0)
		for seller := range sellers {
			if !seller.transacted {
				remainingSellers = append(remainingSellers, seller)
			}
		}
		for buyer := range buyers {
			if !buyer.transacted {
				remainingBuyers = append(remainingBuyers, buyer)
			}
		}

		for i := 0; i < len(remainingBuyers) && i < len(remainingSellers); i++ {
			// attempt to transact
			seller := remainingSellers[i]
			buyer := remainingBuyers[i]

			offerPrice := seller.expectedValues[ROCKET] // we would use max, but we know sellers exp > per
			// need some wiggle room 									V
			willingBuyPrice := math.Min(buyer.expectedValues[ROCKET], buyer.personalValues[ROCKET])
			if offerPrice <= willingBuyPrice {
				// transaction made!
				seller.transacted = true
				buyer.transacted = true

				transactionMade = true
			}
		}
	}

	// update expected values
	for _, actor := range actors {
		if sellers[actor] {
			if actor.transacted {
				// seller sold
				actor.expectedValues[ROCKET]++
			} else {
				actor.expectedValues[ROCKET]--

			}
		} else {
			if actor.transacted {
				// buyer bought
				actor.expectedValues[ROCKET]--
			} else {
				actor.expectedValues[ROCKET]++

			}
		}
	}
}

func ChangePVals() {
	for _, actor := range actors {
		actor.personalValues[ROCKET] += rand.Float64() * 0.5
	}
}

func intMin(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func intMax(i, j int) int {
	if i < j {
		return j
	}
	return i
}
