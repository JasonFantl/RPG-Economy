package economy

import (
	"math"
	"math/rand"
)

// implement physical attributes later
type Actor struct {
	basePersonalValues     map[Good]float64 // if they had no goods, how much would they be willing to pay for this good?
	adjustedPersonalValues map[Good]float64 // changes based on if they have a good

	expectedValues map[Good]float64
	// 0.0 -> don't update.
	// 0.5 -> half of our expected value comes from the last signal.
	// 1.0 -> last signal is the expected value.

	assets        map[Good]int
	desiredAssets map[Good]int

	failedTransactionAttempts int     // how many times we have to fail to buy/sell
	failedTransactionThresh   int     // how many times we have to fail to buy/sell before updating our expected value
	priceSignalReactivity     float64 // how much we should change our expected values by
}

func NewActor() *Actor {
	d := rand.Intn(10) + 10
	actor := &Actor{
		basePersonalValues: map[Good]float64{
			FOOD:   rand.Float64()*10 + 10,
			ROCKET: float64(d),
		},
		adjustedPersonalValues: make(map[Good]float64),
		expectedValues: map[Good]float64{
			FOOD:   0,
			ROCKET: 0,
		},
		assets: map[Good]int{
			MONEY:  100,
			FOOD:   10,
			ROCKET: 10,
		},
		desiredAssets: map[Good]int{
			FOOD:   3,
			ROCKET: d * 2,
		},
		failedTransactionAttempts: 0,
		failedTransactionThresh:   5,
		priceSignalReactivity:     1.0,
	}

	return actor
}

func (actor *Actor) Update() {
	if rand.Float64() > 0.5 { // simulates time between activities
		return
	}

	if rand.Float64() > 0.01 { // throw some randomness in (helps for non-sellers/buyers who need to interact with the market)
		actor.expectedValues[ROCKET] += rand.Float64() - 0.5
	}

	// desired equation
	for _, good := range goods {
		actor.calcAdjustedValue(good)
	}

	// we assume buyers are the initiators, so we ignore sellers
	if actor.isBuyer(ROCKET) {
		willingBuyPrice := actor.willingBuyPrice(ROCKET)
		// look for a seller, simulates going from shop to shop
		var seller *Actor = nil
		sellingPrice := 0
		for otherActor := range actors { // rely on random iteration
			// do we want to update our expected value after each failed shop visit, or only after exhausting every shop? Doing only after every for now
			if !otherActor.isSeller(ROCKET) || !otherActor.canSell(ROCKET) { // must be a seller with goods to sell
				continue
			}
			sellingPrice = otherActor.willingSellPrice(ROCKET) // looking at the price tag
			if willingBuyPrice < sellingPrice {                // price is too high
				continue
			}
			if !actor.canBuy(sellingPrice) { // we don't have the money
				continue
			}

			// made it past all the checks, this is someone we can buy from
			seller = otherActor
			break
		}

		if seller != nil { // found a seller we can buy from
			// update expected values and asset count
			actor.buyGood(ROCKET, sellingPrice)
			seller.sellGood(ROCKET, sellingPrice)
		} else {
			actor.failedTransactionAttempts++
		}
	}

	// sellers need to update if they haven't made a sale in a while
	if actor.isSeller(ROCKET) && actor.canSell(ROCKET) {
		actor.failedTransactionAttempts++
	}

	if actor.failedTransactionAttempts > actor.failedTransactionThresh { // haven't transacted in a while, update expected values
		actor.failedTransactionAttempts = 0
		if actor.isBuyer(ROCKET) {
			actor.expectedValues[ROCKET] += actor.priceSignalReactivity
		} else if actor.isSeller(ROCKET) {
			actor.expectedValues[ROCKET] -= actor.priceSignalReactivity
		}
	}
}

// approximation, see README: APPROXIMATION 1
func (actor *Actor) calcAdjustedValue(good Good) {
	x := float64(actor.assets[good])
	d := float64(actor.desiredAssets[good])
	b := actor.basePersonalValues[good]

	if d < 1 {
		actor.adjustedPersonalValues[good] = 0
	} else {
		// linear
		// actor.adjustedPersonalValues[good] = b * (1 - x/d)
		// exponential
		actor.adjustedPersonalValues[good] = b * (1 - x*x/(d*d))
	}
}

func (actor *Actor) sellGood(good Good, cost int) {
	actor.assets[good]--
	actor.assets[MONEY] += cost
	actor.expectedValues[good] += actor.priceSignalReactivity
	actor.failedTransactionAttempts = 0
}

func (actor *Actor) buyGood(good Good, cost int) {
	actor.assets[good]++
	actor.assets[MONEY] -= cost
	actor.expectedValues[good] -= actor.priceSignalReactivity
	actor.failedTransactionAttempts = 0
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
	return actor.assets[MONEY] >= cost
}

func (actor *Actor) canSell(good Good) bool {
	return actor.assets[good] > 0
}
